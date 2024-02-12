terraform {
  backend "s3" {
    bucket = "nameslol-deployments"
    key    = "terraform/name-updater-producer"
    region = "us-east-1"
  }
}

variable "app_name" {
  type    = string
  default = "name-updater-producer"
}

provider "aws" {
  region = "us-east-1"
}

data "aws_dynamodb_table" "nameslol" {
  name = "nameslol"
}

data "aws_sqs_queue" "name-update-queue" {
  name = "NameUpdateQueue"
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/bootstrap"
  output_path = "${path.module}/bootstrap.zip"
}

data "local_file" "lambda_zip_contents" {
  filename = data.archive_file.lambda_zip.output_path
}

resource "aws_iam_role" "lambda_exec" {
  name               = "${var.app_name}_execution_role"
  assume_role_policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Action" : "sts:AssumeRole",
        "Principal" : {
          "Service" : "lambda.amazonaws.com"
        },
        "Effect" : "Allow"
      },
    ]
  })
}

resource "aws_iam_role_policy" "lambda_exec_policy" {
  role   = aws_iam_role.lambda_exec.id
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Effect" : "Allow",
        "Action" : [
          "dynamodb:Query"
        ],
        "Resource" : [
          data.aws_dynamodb_table.nameslol.arn,
          "${data.aws_dynamodb_table.nameslol.arn}/index/*"
        ]
      },
      {
        "Effect" : "Allow",
        "Action" : [
          "sqs:SendMessage",
          "sqs:BatchSendMessage",
        ],
        "Resource" : [
          data.aws_sqs_queue.name-update-queue.arn
        ]
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_lambda_function" "default" {
  function_name    = var.app_name
  architectures    = ["arm64"]
  memory_size      = 256
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_exec.arn
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.local_file.lambda_zip_contents.content_md5
  runtime          = "provided.al2023"
  timeout          = 60

  environment {
    variables = {
      QUEUE_URL      = data.aws_sqs_queue.name-update-queue.url
      DYNAMODB_TABLE = data.aws_dynamodb_table.nameslol.name
    }
  }
}


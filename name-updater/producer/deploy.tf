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

data "aws_ssm_parameter" "riot-api-token" {
  name = "/riot-api-token"
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

resource "aws_iam_role" "scheduler_exec" {
  name = "${var.app_name}_scheduler_role"
  assume_role_policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": "scheduler.amazonaws.com"
        },
        "Action": "sts:AssumeRole"
      }
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

resource "aws_iam_role_policy" "scheduler_exec_policy" {
  role = aws_iam_role.scheduler_exec.id
  policy = jsonencode({
    "Version": "2012-10-17",
    "Statement": [
      {
        "Action": [
          "lambda:InvokeFunction"
        ],
        "Effect": "Allow",
        "Resource": [
            aws_lambda_function.default.arn
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
      RIOT_API_TOKEN = data.aws_ssm_parameter.riot-api-token.value
    }
  }
}

resource "aws_scheduler_schedule" "hourly" {
  name = "name-updater-hourly"
  schedule_expression = "cron(0 * ? * * *)" // Every hour
  state = "DISABLED"

  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn = aws_lambda_function.default.arn
    role_arn = aws_iam_role.scheduler_exec.arn

    input = jsonencode({
      "refreshType": "hourly"
    })

    retry_policy {
      maximum_retry_attempts = 0
    }
  }
}

resource "aws_scheduler_schedule" "weekly" {
  name = "name-updater-weekly"
  schedule_expression = "cron(0 0 ? * 6 *)" // Every Friday
  state = "DISABLED"

  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn = aws_lambda_function.default.arn
    role_arn = aws_iam_role.scheduler_exec.arn

    input = jsonencode({
      "refreshType": "weekly"
    })

    retry_policy {
      maximum_retry_attempts = 0
    }
  }
}

resource "aws_scheduler_schedule" "monthly" {
  name = "name-updater-monthly"
  schedule_expression = "cron(0 0 1 * ? *)" // First day of each month
  state = "DISABLED"

  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn = aws_lambda_function.default.arn
    role_arn = aws_iam_role.scheduler_exec.arn

    input = jsonencode({
      "refreshType": "monthly"
    })

    retry_policy {
      maximum_retry_attempts = 0
    }
  }
}

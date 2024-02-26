terraform {
  backend "s3" {
    bucket = "nameslol-deployments"
    key    = "terraform/name-updater-consumer"
    region = "us-east-1"
  }
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

module "lambda" {
  source                = "../../infrastructure/modules/lambda"
  app_name              = "name-updater-consumer"
  bootstrap_file_path   = "${path.module}/bootstrap"
  timeout               = 30
  memory_size           = 256
  iam_policy_statements = [
    {
      "Effect" : "Allow",
      "Action" : [
        "dynamodb:PutItem",
        "dynamodb:DeleteItem",
      ],
      "Resource" : [
        data.aws_dynamodb_table.nameslol.arn,
        "${data.aws_dynamodb_table.nameslol.arn}/index/*"
      ]
    },
    {
      "Effect" : "Allow",
      "Action" : [
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes",
      ],
      "Resource" : [
        data.aws_sqs_queue.name-update-queue.arn
      ]
    }
  ]
  environment_variables = {
    DYNAMODB_TABLE = data.aws_dynamodb_table.nameslol.name
    RIOT_API_TOKEN = data.aws_ssm_parameter.riot-api-token.value
  }
}

resource "aws_lambda_event_source_mapping" "default" {
  event_source_arn = data.aws_sqs_queue.name-update-queue.arn
  function_name    = module.lambda.lambda_function_arn
  batch_size       = 5
  scaling_config {
    maximum_concurrency = 10 // Keep below 10 to avoid exceeding Riot API rate limits
  }
}

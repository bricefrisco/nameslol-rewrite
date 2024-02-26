terraform {
  backend "s3" {
    bucket = "nameslol-deployments"
    key    = "terraform/name-updater-producer"
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
  source = "../../infrastructure/modules/lambda"
  app_name = "name-updater-producer"
  bootstrap_file_path = "${path.module}/bootstrap"
  timeout = 60
  memory_size = 256
  iam_policy_statements = [
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
  environment_variables = {
    QUEUE_URL      = data.aws_sqs_queue.name-update-queue.url
    DYNAMODB_TABLE = data.aws_dynamodb_table.nameslol.name
    RIOT_API_TOKEN = data.aws_ssm_parameter.riot-api-token.value
  }
}

resource "aws_iam_role" "scheduler_exec" {
  name = "name-updater-producer-scheduler-role"
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
          module.lambda.lambda_function_arn
        ]
      }
    ]
  })
}

resource "aws_scheduler_schedule" "hourly" {
  name = "name-updater-hourly"
  schedule_expression = "cron(0 * ? * * *)" // Every hour
  state = "DISABLED"

  flexible_time_window {
    mode = "OFF"
  }

  target {
    arn = module.lambda.lambda_function_arn
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
    arn = module.lambda.lambda_function_arn
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
    arn = module.lambda.lambda_function_arn
    role_arn = aws_iam_role.scheduler_exec.arn

    input = jsonencode({
      "refreshType": "monthly"
    })

    retry_policy {
      maximum_retry_attempts = 0
    }
  }
}

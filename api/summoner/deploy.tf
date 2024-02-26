terraform {
  backend "s3" {
    bucket = "nameslol-deployments"
    key    = "terraform/api-summoner"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

data "aws_dynamodb_table" "nameslol" {
  name = "nameslol"
}

data "aws_ssm_parameter" "riot-api-token" {
  name = "/riot-api-token"
}

module "lambda" {
  source = "../../infrastructure/modules/lambda"
  app_name = "api-summoner"
  bootstrap_file_path = "${path.module}/bootstrap"
  timeout = 15
  memory_size = 256
  iam_policy_statements = [
    {
      "Effect" : "Allow",
      "Action" : [
        "dynamodb:PutItem"
      ],
      "Resource" : [
        data.aws_dynamodb_table.nameslol.arn,
        "${data.aws_dynamodb_table.nameslol.arn}/index/*"
      ]
    },
  ]
  environment_variables = {
    DYNAMODB_TABLE = data.aws_dynamodb_table.nameslol.name
    RIOT_API_TOKEN = data.aws_ssm_parameter.riot-api-token.value
    CORS_ORIGINS   = "http://localhost:3000"
    CORS_METHODS   = "GET, OPTIONS"
  }
}

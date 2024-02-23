variable "app_name" {
  description = "The name of the application or lambda function."
  type        = string
}

variable "bootstrap_file_path" {
  description = "The path to the bootstrap file to use for the lambda function."
  type        = string
}

variable "iam_policy_statements" {
    description = "A list of IAM policy statements to attach to the lambda function."
    type        = list(object({
        Effect   = string
        Action  = list(string)
        Resource = list(string)
    }))
}

variable "environment_variables" {
  description = "A map of environment variables to set on the lambda function."
  type        = map(string)
}

variable "timeout" {
  description = "The amount of time the lambda function has to run before it times out."
  type        = number
}

variable "memory_size" {
  description = "The amount of memory to allocate to the lambda function."
  type        = number
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = var.bootstrap_file_path
  output_path = "${var.bootstrap_file_path}.zip"
}

data "local_file" "lambda_zip_contents" {
  filename = data.archive_file.lambda_zip.output_path
}

resource "aws_iam_role" "lambda_exec" {
  name               = "${var.app_name}-execution-role"
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

resource "aws_iam_role_policy_attachment" "lambda_basic_execution" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy" "lambda_exec_policy" {
  role   = aws_iam_role.lambda_exec.id
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : var.iam_policy_statements
  })
}

resource "aws_lambda_function" "default" {
  function_name    = var.app_name
  memory_size      = var.memory_size
  timeout          = var.timeout
  architectures    = ["arm64"]
  handler          = "bootstrap"
  role             = aws_iam_role.lambda_exec.arn
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.local_file.lambda_zip_contents.content_md5
  runtime          = "provided.al2023"

  environment {
    variables = var.environment_variables
  }
}
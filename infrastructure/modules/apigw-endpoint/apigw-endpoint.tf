variable "api_gateway_id" {
  description = "The ID of the API Gateway."
  type        = string
}

variable "api_gateway_root_resource_id" {
  description = "The ID of the root resource."
  type        = string
}

variable "api_gateway_execution_arn" {
  description = "The ARN of the API Gateway execution."
  type        = string
}

variable "function_name" {
  description = "The name of the lambda function."
  type        = string
}

variable "path" {
  description = "API path"
  type        = string
}

data "aws_lambda_function" "lambda-api" {
  function_name = var.function_name
}

resource "aws_api_gateway_resource" "api-resource" {
  rest_api_id = var.api_gateway_id
  parent_id   = var.api_gateway_root_resource_id
  path_part   = var.path
}

resource "aws_api_gateway_method" "api-method" {
  rest_api_id = var.api_gateway_id
  resource_id = aws_api_gateway_resource.api-resource.id
  http_method = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "api-integration" {
  rest_api_id = var.api_gateway_id
  resource_id = aws_api_gateway_resource.api-resource.id
  http_method = aws_api_gateway_method.api-method.http_method
  integration_http_method = "POST"
  type = "AWS_PROXY"
  uri = data.aws_lambda_function.lambda-api.invoke_arn
}

resource "aws_lambda_permission" "api-permission" {
  statement_id = "AllowAPIGatewayInvoke"
  action = "lambda:InvokeFunction"
  function_name = data.aws_lambda_function.lambda-api.function_name
  principal = "apigateway.amazonaws.com"
  source_arn = "${var.api_gateway_execution_arn}/*/*"
}

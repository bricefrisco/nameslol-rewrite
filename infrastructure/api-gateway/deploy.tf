terraform {
  backend "s3" {
    bucket = "nameslol-deployments"
    key    = "terraform/api-gateway"
    region = "us-east-1"
  }
}

provider "aws" {
  region = "us-east-1"
}

/***************************
 * API Gateway Resource + Deployment
 ***************************/
resource "aws_api_gateway_rest_api" "default" {
  name = "nameslol-api"
  description ="NamesLoL API"
}

resource "aws_api_gateway_deployment" "deployment" {
  depends_on  = [
    aws_api_gateway_integration.api-summoner-integration
  ]
  stage_description = "Deployment: #1"
  rest_api_id = aws_api_gateway_rest_api.default.id
  stage_name  = "prod"
}

/***************************
 * Summoner API (/summoner)
 ***************************/
data "aws_lambda_function" "api-summoner" {
  function_name = "api-summoner"
}

resource "aws_api_gateway_resource" "api-summoner-resource" {
  rest_api_id = aws_api_gateway_rest_api.default.id
  parent_id   = aws_api_gateway_rest_api.default.root_resource_id
  path_part   = "summoner"
}

resource "aws_api_gateway_method" "api-summoner-method" {
    rest_api_id = aws_api_gateway_rest_api.default.id
    resource_id = aws_api_gateway_resource.api-summoner-resource.id
    http_method = "GET"
    authorization = "NONE"
}

resource "aws_api_gateway_integration" "api-summoner-integration" {
    rest_api_id = aws_api_gateway_rest_api.default.id
    resource_id = aws_api_gateway_resource.api-summoner-resource.id
    http_method = aws_api_gateway_method.api-summoner-method.http_method
    integration_http_method = "POST"
    type = "AWS_PROXY"
    uri = data.aws_lambda_function.api-summoner.invoke_arn
}

resource "aws_lambda_permission" "api-summoner-permission" {
    statement_id = "AllowAPIGatewayInvoke"
    action = "lambda:InvokeFunction"
    function_name = data.aws_lambda_function.api-summoner.function_name
    principal = "apigateway.amazonaws.com"
    source_arn = "${aws_api_gateway_rest_api.default.execution_arn}/*/*"
}

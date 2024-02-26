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
    module.summoner-apigw-endpoint,
    module.summoners-apigw-endpoint
  ]
  stage_description = "Deployment: #3"
  rest_api_id = aws_api_gateway_rest_api.default.id
  stage_name  = "prod"
}

/***************************
 * API Gateway Endpoints
 ***************************/
module "summoner-apigw-endpoint" {
  source = "../modules/apigw-endpoint"
  api_gateway_id = aws_api_gateway_rest_api.default.id
  api_gateway_root_resource_id = aws_api_gateway_rest_api.default.root_resource_id
  api_gateway_execution_arn = aws_api_gateway_rest_api.default.execution_arn
  function_name = "api-summoner"
  path = "summoner"
}

module "summoners-apigw-endpoint" {
  source = "../modules/apigw-endpoint"
  api_gateway_id = aws_api_gateway_rest_api.default.id
  api_gateway_root_resource_id = aws_api_gateway_rest_api.default.root_resource_id
  api_gateway_execution_arn = aws_api_gateway_rest_api.default.execution_arn
  function_name = "api-summoners"
  path = "summoners"
}

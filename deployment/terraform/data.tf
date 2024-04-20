data "aws_api_gateway_rest_api" "lofyd_api" {
  name = "lofyd-api"
}

data "aws_api_gateway_resource" "forgotpassword_resource" {
  path = "/forgotpassword"
  rest_api_id = data.aws_api_gateway_rest_api.lofyd_api.id
}

data "aws_iam_role" "apigateway_role" {
  name = "lofyd-ApiGatewayRole"
}

data "aws_secretsmanager_secret" "mobile_app_cognito" {
  name = "/lofyd/${var.stage}/mobileapp/cognito" 
}

data "aws_iam_role" "lambda_role"{
  name = "${module.lambda.lambda_function_name}"
}

data "aws_iam_policy" "read_secret_policy" {
  name = "lofyd-ReadSecretsPolicy"
}
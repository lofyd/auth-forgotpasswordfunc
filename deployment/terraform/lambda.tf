module "lambda" {
  source  = "terraform-aws-modules/lambda/aws" #https://github.com/terraform-aws-modules/terraform-aws-lambda

  function_name = local.function_name
  handler       = "main"
  runtime       = "provided.al2023"
  source_path   = ["../../build"]


  environment_variables = {
    STAGE               = var.stage
  }
}

resource "aws_api_gateway_integration" "apigw_lambda_integration" {
  rest_api_id             = data.aws_api_gateway_rest_api.lofyd_api.id
  resource_id             = data.aws_api_gateway_resource.forgotpassword_resource.id
  http_method             = "POST"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.lambda.lambda_function_invoke_arn
  credentials             = data.aws_iam_role.apigateway_role.arn
}

# This allows apigw POST /signup to invoke lambda function
resource "aws_lambda_permission" "apigw_invoke_lambda_permission" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.lambda.lambda_function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  source_arn = "${data.aws_api_gateway_rest_api.lofyd_api.execution_arn}/*/POST/forgotpassword"
}

#Add read secret policy to lambda's role
resource "aws_iam_role_policy_attachment" "read_secret_policy_attachement" {
  policy_arn = data.aws_iam_policy.read_secret_policy.arn
  role       = data.aws_iam_role.lambda_role.name
}
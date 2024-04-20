# auth-forgotpasswordfunc
Lambda function that backs the API Gateway endpoint `POST /forgotpassword`

## Github Code Repository
- https://github.com/lofyd/auth-forgotpasswordfunc

## Dependent Github Code Repositories
- https://github.com/lofyd/app-infrastructure

## Dependent Services
The following AWS services are deployed by the [lofyd/app-infrastructre](https://github.com/lofyd/app-infrastructure) repository
- [Secrets Manager](https://us-east-1.console.aws.amazon.com/secretsmanager/landing?region=us-east-1)
- [Cognito](https://us-east-1.console.aws.amazon.com/cognito/v2/home?region=us-east-1)

## Run and Test Lambda Locally
[AWS Serverless Application Model](https://docs.aws.amazon.com/serverless-application-model/) (SAM) CLI can be used to run AWS Services locally for testing purposes without deploying the lambda to AWS.
 - Software
   - [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html) (Required)
   - [jq](https://jqlang.github.io/jq/) (Optional) - Formats the JSON response, so it is readable.
 - Files
   - [template.yaml](template.yaml): defines the AWS resources that will run locally. [specification](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-specification.html)
   - [signup_event.json](testevents/forgotpassword_event.json): the apigateway event json this lambda is expecting 
   - Add new *_event.json for different uses cases that need to be tested
 - How-To build and execute lambda locally
   - Build lambda archive file
        ```
        ./build.sh
        ```
   - Execute lambda locally        
        ```
        sam local invoke -e ./testevents/forgotpassword_event.json SignupFunc --profile lofyd-dev | jq
        ``` 
## Build and Deploy Lambda to AWS
- Build lambda archive file
    ```
    ./build.sh
    ```
- Deploy Lambda to AWS
    ```
    cd deployment/terraform
    terraform init
    terraform plan -out=out.txt
    terraform apply "out.txt"
    ```

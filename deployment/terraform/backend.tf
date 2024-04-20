terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
  backend "s3" {
     key            = "lofyd/auth-forgotpassword/dev/us-east-1/terraform.tfstate"
     bucket         = "lofyd-dev-terraform-state"
     dynamodb_table = "lofyd-state-locking-table"
     encrypt        = true
     profile        = "lofyd-dev"
     region         = "us-east-1"
  }
}
variable "region" {
  type        = string
  default     = "us-east1"
  description = "aws region"
}

variable "stage" {
  type        = string
  default     = "dev"
  description = "stage name used as part of resource names. prod, dev, qa. this is populated by the ci pipeline"
}

variable "aws_profile" {
  type        = string
  default     = "lofyd-dev"
  description = "aws profile to use when deploying."
}

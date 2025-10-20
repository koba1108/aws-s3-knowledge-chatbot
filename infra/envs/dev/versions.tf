terraform {
  required_version = "~> 1.13.4"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.60"
    }
  }
}

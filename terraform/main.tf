terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.20"
    }
  }
  required_version = ">= 1.5.7"
}

provider "aws" {
}

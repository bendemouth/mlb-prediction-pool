variable "aws_region" {
    description = "AWS Region"
    type = string
    default = "us-east-1"
}

variable "project_name" {
    description = "Project Name"
    type = string
    default = "mlb-prediction-pool"
}

variable "environment" {
    description = "Environment (dev, prod)"
    type = string
    default = "dev"
}

variable "mlb_api_base" {
    description = "MLB Stats API base URL"
    type = string
    default = "https://statsapi.mlb.com/api/v1"
}
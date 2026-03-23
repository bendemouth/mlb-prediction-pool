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

variable "ec2_key_public_key" {
    description = "SSH public key material for the EC2 deployer key pair (e.g. contents of ~/.ssh/id_rsa.pub)"
    type        = string
}

variable "ec2_instance_type" {
    description = "EC2 instance type for the application server"
    type        = string
    default     = "t3a.micro"
}
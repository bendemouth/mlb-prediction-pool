# Security Group for the application EC2 instance
resource "aws_security_group" "app_server" {
    name        = "${var.project_name}-${var.environment}-app-sg"
    description = "Allow HTTP and SSH inbound traffic"

    ingress {
        description = "HTTP"
        from_port   = 80
        to_port     = 80
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        description = "HTTPS (reserved for future TLS termination)"
        from_port   = 443
        to_port     = 443
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        description = "SSH"
        from_port   = 22
        to_port     = 22
        protocol    = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        from_port   = 0
        to_port     = 0
        protocol    = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# SSH key pair for EC2 access
resource "aws_key_pair" "deployer" {
    key_name   = "${var.project_name}-${var.environment}-deployer"
    public_key = var.ec2_key_public_key

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Amazon Linux 2023 latest AMI
data "aws_ami" "amazon_linux_2023" {
    most_recent = true
    owners      = ["amazon"]

    filter {
        name   = "name"
        values = ["al2023-ami-2023.*-x86_64"]
    }

    filter {
        name   = "virtualization-type"
        values = ["hvm"]
    }
}

# EC2 instance running Docker + Docker Compose
resource "aws_instance" "app_server" {
    ami                    = data.aws_ami.amazon_linux_2023.id
    instance_type          = var.ec2_instance_type
    key_name               = aws_key_pair.deployer.key_name
    vpc_security_group_ids = [aws_security_group.app_server.id]
    iam_instance_profile   = aws_iam_instance_profile.ec2_backend_profile.name

    user_data = <<-EOF
        #!/bin/bash
        set -e

        # Install Docker
        dnf update -y
        dnf install -y docker
        systemctl enable docker
        systemctl start docker
        usermod -aG docker ec2-user

        # Install Docker Compose v2 plugin
        mkdir -p /usr/local/lib/docker/cli-plugins
        curl -SL "https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64" \
            -o /usr/local/lib/docker/cli-plugins/docker-compose
        chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

        # Create app directory
        mkdir -p /opt/app
        chown ec2-user:ec2-user /opt/app
    EOF

    root_block_device {
        volume_size = 20
        volume_type = "gp3"
        encrypted   = true
    }

    tags = {
        Project     = var.project_name
        Environment = var.environment
        Name        = "${var.project_name}-${var.environment}-app-server"
    }
}

# Elastic IP for a stable public address
resource "aws_eip" "app_server" {
    instance = aws_instance.app_server.id
    domain   = "vpc"

    tags = {
        Project     = var.project_name
        Environment = var.environment
        Name        = "${var.project_name}-${var.environment}-eip"
    }
}

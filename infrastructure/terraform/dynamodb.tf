# Users Table
resource "aws_dynamodb_table" "users" {
    name = "${var.project_name}-${var.environment}-users"
    billing_mode = "PAY_PER_REQUEST"

    attribute {
        name = "userId"
        type = "S"
    }

    hash_key = "userId"

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Predictions Table
resource "aws_dynamodb_table" "predictions" {
    name = "${var.project_name}-${var.environment}-predictions"
    billing_mode = "PAY_PER_REQUEST"

    attribute {
        name = "userId"
        type = "S"
    }

    attribute {
        name = "gameId"
        type = "S"
    }

    hash_key  = "userId"
    range_key = "gameId"

    global_secondary_index {
        name            = "GameIdIndex"
        hash_key        = "gameId"
        projection_type = "ALL"
    }

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Games Table
resource "aws_dynamodb_table" "games" {
    name = "${var.project_name}-${var.environment}-games"
    billing_mode = "PAY_PER_REQUEST"

    attribute {
        name = "gameId"
        type = "S"
    }

    hash_key = "gameId"

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Models table 
resource "aws_dynamodb_table" "models" {
    name = "${var.project_name}-${var.environment}-models"
    billing_mode = "PAY_PER_REQUEST"

    attribute {
        name = "modelId"
        type = "S"
    }

    hash_key = "modelId"

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}
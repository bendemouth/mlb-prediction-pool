# Users Table
resource "aws_dynamodb_table" "users" {
    name = "${var.project_name}-${var.environment}-users"
    billing_mode = "PAY_PER_REQUEST"
    hash_key = "user_id"

    attribute {
        name = "userId"
        type = "S"
    }

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Predictions Table
resource "aws_dynamodb_table" "predictions" {
    name = "${var.project_name}-${var.environment}-predictions"
    billing_mode = "PAY_PER_REQUEST"
    hash_key = "prediction_id"

    attribute {
        name = "predictionId"
        type = "S"
    }

    attribute {
        name = "gameId"
        type = "S"
    }

    global_secondary_index {
        name               = "GameIndex"
        hash_key          = "gameId"
        projection_type   = "ALL"
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
    hash_key = "game_id"

    attribute {
        name = "gameId"
        type = "S"
    }

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Models table 
resource "aws_dynamodb_table" "models" {
    name = "${var.project_name}-${var.environment}-models"
    billing_mode = "PAY_PER_REQUEST"
    hash_key = "model_id"

    attribute {
        name = "modelId"
        type = "S"
    }

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}
# Automatically zip the Lambda function
data "archive_file" "lambda_zip" {
    type        = "zip"
    source_dir  = "${path.module}/../../lambda/data-ingestion"
    output_path = "${path.module}/../../lambda/data-ingestion/lambda.zip"
    excludes    = ["lambda.zip", "__pycache__", "*.pyc"]
}

resource "aws_lambda_function" "data_ingestion" {
    filename         = data.archive_file.lambda_zip.output_path
    function_name    = "${var.project_name}-data-ingestion-${var.environment}"
    role            = aws_iam_role.lambda_exec_role.arn
    handler         = "handler.lambda_handler"
    runtime         = "python3.11"
    timeout         = 60
    memory_size     = 256
    
    # Update lambda if changes are made
    source_code_hash = data.archive_file.lambda_zip.output_base64sha256

    environment {
        variables = {
            GAMES_TABLE        = aws_dynamodb_table.games.name
            PREDICTIONS_TABLE  = aws_dynamodb_table.predictions.name
            USERS_TABLE        = aws_dynamodb_table.users.name
            MODELS_TABLE       = aws_dynamodb_table.models.name
            DATA_BUCKET        = aws_s3_bucket.mlb_data.bucket
            USER_MODELS_BUCKET = aws_s3_bucket.user_models.bucket
        }
    }

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

resource "aws_cloudwatch_log_group" "data_ingestion_logs" {
    name              = "/aws/lambda/${aws_lambda_function.data_ingestion.function_name}"
    retention_in_days = 5
}
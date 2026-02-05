output "dynamodb_tables" {
    description = "DynamoDB table names"
    value = {
        users_table       = aws_dynamodb_table.users.name
        predictions_table = aws_dynamodb_table.predictions.name
        games_table       = aws_dynamodb_table.games.name
        models = aws_dynamodb_table.models.name
    }
}

output "s3_buckets" {
    description = "S3 bucket names"
    value = {
        mlb_data_bucket  = aws_s3_bucket.mlb_data.bucket
        user_models_bucket = aws_s3_bucket.user_models.bucket
    }
}

output "lambda_functions" {
    description = "Lambda function ARNs"
    value = {
        data_ingestion_lambda = aws_lambda_function.data_ingestion.arn
    }
}
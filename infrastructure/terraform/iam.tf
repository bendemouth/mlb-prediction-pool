resource "aws_iam_role" "lambda_exec_role" {
    name = "${var.project_name}-${var.environment}-lambda-exec-role"
    assume_role_policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Action = "sts:AssumeRole"
                Effect = "Allow"
                Principal = {
                    Service = "lambda.amazonaws.com"
                }
            }
        ]
    })

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

resource "aws_iam_role_policy" "lambda_dynamodb_policy" {
    name = "dynamodb-access-policy"
    role = aws_iam_role.lambda_exec_role.id

    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = [
                    "dynamodb:PutItem",
                    "dynamodb:GetItem",
                    "dynamodb:Query",
                    "dynamodb:Scan",
                    "dynamodb:UpdateItem",
                    "dynamodb:BatchWriteItem"
                ]
                Resource = [
                    aws_dynamodb_table.games.arn,
                    aws_dynamodb_table.predictions.arn,
                    aws_dynamodb_table.users.arn,
                    aws_dynamodb_table.models.arn
                ]
            }
        ]
    })
}

# Add S3 access policy
resource "aws_iam_role_policy" "lambda_s3_policy" {
    name = "s3-access-policy"
    role = aws_iam_role.lambda_exec_role.id

    policy = jsonencode({
        Version = "2012-10-17"
        Statement = [
            {
                Effect = "Allow"
                Action = [
                    "s3:GetObject",
                    "s3:PutObject",
                    "s3:ListBucket"
                ]
                Resource = [
                    aws_s3_bucket.mlb_data.arn,
                    "${aws_s3_bucket.mlb_data.arn}/*",
                    aws_s3_bucket.user_models.arn,
                    "${aws_s3_bucket.user_models.arn}/*"
                ]
            }
        ]
    })
}

# Add CloudWatch Logs policy
resource "aws_iam_role_policy_attachment" "lambda_logs" {
    role       = aws_iam_role.lambda_exec_role.name
    policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}
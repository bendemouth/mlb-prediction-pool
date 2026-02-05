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
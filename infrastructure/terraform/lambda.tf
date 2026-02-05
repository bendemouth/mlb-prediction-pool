resource "aws_lambda_function" "data_ingestion" {
    filename         = "${path.module}/../../lambda/data-ingestion/lambda.zip"
    function_name    = "${var.project_name}-data-ingestion-${var.environment}"
    role            = aws_iam_role.lambda_exec_role.arn
    handler         = "handler.lambda_handler"
    runtime         = "python3.11"
    timeout         = 60
    memory_size     = 256
}

resource "aws_cloudwatch_log_group" "data_ingestion_logs" {
    name = "/aws/lambda/${aws_lambda_function.data_ingestion.function_name}"
    retention_in_days = 5
}
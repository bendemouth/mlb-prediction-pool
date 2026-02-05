# EventBridge rule for daily data ingestion at 5 AM CST
resource "aws_cloudwatch_event_rule" "daily_data_ingestion" {
    name                = "${var.project_name}-daily-data-ingestion-${var.environment}"
    description         = "Trigger MLB data ingestion daily at 5 AM CST"
    schedule_expression = "cron(0 11 * * ? *)"  # 5 AM CST = 11 AM UTC

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Invoke Lambda Function
resource "aws_cloudwatch_event_target" "invoke_data_ingestion" {
    target_id = "InvokeDataIngestionLambda"
    rule     = aws_cloudwatch_event_rule.daily_data_ingestion.name
    arn      = aws_lambda_function.data_ingestion.arn
}

resource "aws_lambda_permission" "allow_eventbridge_invoke" {
    statement_id  = "AllowEventBridgeInvoke"
    action        = "lambda:InvokeFunction"
    function_name = aws_lambda_function.data_ingestion.function_name
    principal     = "events.amazonaws.com"
    source_arn    = aws_cloudwatch_event_rule.daily_data_ingestion.arn
}
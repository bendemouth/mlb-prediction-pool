resource "aws_s3_bucket" "mlb_data" {
    bucket = "${var.project_name}-${var.environment}-mlb-data"

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

resource "aws_s3_bucket_acl" "mlb_data_acl" {
    bucket = aws_s3_bucket.mlb_data.id
    acl    = "private"
}

resource "aws_s3_bucket" "user_models" {
    bucket = "${var.project_name}-${var.environment}-user-models"

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

resource "aws_s3_bucket_acl" "user_models_acl" {
    bucket = aws_s3_bucket.user_models.id
    acl    = "private"
}

resource "aws_s3_bucket_versioning" "user_models_versioning" {
    bucket = aws_s3_bucket.user_models.id

    versioning_configuration {
        status = "Enabled"
    }
}

resource "aws_s3_bucket_lifecycle_configuration" "mlb_data_lifecycle" {
    bucket = aws_s3_bucket.mlb_data.id

    rule {
        id     = "ExpireOldData"
        status = "Enabled"

        transition {
            days = 30
            storage_class = "STANDARD_IA"
        }

        transition {
            days = 90
            storage_class = "GLACIER_IR"
        }
    }
}
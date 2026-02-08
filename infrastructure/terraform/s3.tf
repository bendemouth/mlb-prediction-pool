resource "aws_s3_bucket" "mlb_data" {
    bucket = "${var.project_name}-${var.environment}-mlb-data"

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Add public access block to ensure bucket is private
resource "aws_s3_bucket_public_access_block" "mlb_data" {
    bucket = aws_s3_bucket.mlb_data.id

    block_public_acls       = true
    block_public_policy     = true
    ignore_public_acls      = true
    restrict_public_buckets = true
}

resource "aws_s3_bucket" "user_models" {
    bucket = "${var.project_name}-${var.environment}-user-models"

    tags = {
        Project     = var.project_name
        Environment = var.environment
    }
}

# Add public access block to ensure bucket is private
resource "aws_s3_bucket_public_access_block" "user_models" {
    bucket = aws_s3_bucket.user_models.id

    block_public_acls       = true
    block_public_policy     = true
    ignore_public_acls      = true
    restrict_public_buckets = true
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
            days          = 30
            storage_class = "STANDARD_IA"
        }

        transition {
            days          = 90
            storage_class = "GLACIER_IR"
        }
    }
}
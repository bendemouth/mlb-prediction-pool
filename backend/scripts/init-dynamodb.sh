#!/bin/bash

# Wait for DynamoDB local
echo "Waiting for DynamoDB Local..."
until curl -s http://dynamodb-local:8000/ > /dev/null 2>&1; do
    echo "DynamoDB not ready yet..."
    sleep 2
done

echo "DynamoDB is ready! Creating tables..."

# Set fake AWS credentials
export AWS_ACCESS_KEY_ID=fakeAccessKeyId
export AWS_SECRET_ACCESS_KEY=fakeSecretAccessKey
export AWS_DEFAULT_REGION=us-west-2

# Create Users Table
aws dynamodb create-table \
    --table-name mlb-prediction-pool-dev-users \
    --attribute-definitions AttributeName=userId,AttributeType=S \
    --key-schema AttributeName=userId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://dynamodb-local:8000 \
    --region us-west-2 || echo "Users table already exists"

# Create Games Table
aws dynamodb create-table \
    --table-name mlb-prediction-pool-dev-games \
    --attribute-definitions AttributeName=gameId,AttributeType=S \
    --key-schema AttributeName=gameId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://dynamodb-local:8000 \
    --region us-west-2 || echo "Games table already exists"

# Create Predictions Table
aws dynamodb create-table \
    --table-name mlb-prediction-pool-dev-predictions \
    --attribute-definitions \
        AttributeName=userId,AttributeType=S \
        AttributeName=gameId,AttributeType=S \
    --key-schema \
        AttributeName=userId,KeyType=HASH \
        AttributeName=gameId,KeyType=RANGE \
    --global-secondary-indexes \
        "IndexName=GameIdIndex,KeySchema=[{AttributeName=gameId,KeyType=HASH}],Projection={ProjectionType=ALL}" \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://dynamodb-local:8000 \
    --region us-west-2 || echo "Predictions table already exists"

# Create Models Table
aws dynamodb create-table \
    --table-name mlb-prediction-pool-dev-models \
    --attribute-definitions AttributeName=modelId,AttributeType=S \
    --key-schema AttributeName=modelId,KeyType=HASH \
    --billing-mode PAY_PER_REQUEST \
    --endpoint-url http://dynamodb-local:8000 \
    --region us-west-2 || echo "Models table already exists"

echo "Tables created successfully!"
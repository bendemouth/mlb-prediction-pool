module github.com/bendemouth/mlb-prediction-pool

go 1.24

require (
	github.com/aws/aws-sdk-go-v2 v1.41.3 // direct
	github.com/aws/aws-sdk-go-v2/config v1.32.11 // direct
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.31 // direct
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.54.0 // direct
	github.com/joho/godotenv v1.5.1
)

require github.com/aws/aws-sdk-go-v2/service/s3 v1.96.4

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.19.11 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.5 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.32.10 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.16 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.8 // indirect
	github.com/aws/smithy-go v1.24.2 // indirect
)

require github.com/golang-jwt/jwt/v5 v5.3.1 // direct

module get-rate-by-point

go 1.17

require (
	covid19/models v0.0.0-00010101000000-000000000000
	covid19/shared v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.37.0
)

require github.com/go-sql-driver/mysql v1.7.0 // indirect

replace covid19/shared => ../../shared

replace covid19/models => ../../models

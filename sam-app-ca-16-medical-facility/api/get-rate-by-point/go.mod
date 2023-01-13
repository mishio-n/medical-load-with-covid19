module get-rate-by-point

go 1.17

require (
	covid19/models v0.0.0-00010101000000-000000000000
	covid19/shared v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.37.0
)

require (
	github.com/DATA-DOG/go-txdb v0.1.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace covid19/shared => ../../shared

replace covid19/models => ../../models

module find-by-prefecture

go 1.17

require (
	covid19/shared v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.36.1
)

require github.com/go-sql-driver/mysql v1.7.0 // indirect

replace covid19/shared => ../../shared

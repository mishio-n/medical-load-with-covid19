require (
	covid19/models v0.0.0-00010101000000-000000000000
	covid19/shared v0.0.0-00010101000000-000000000000
	github.com/aws/aws-lambda-go v1.23.0
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

replace covid19/shared => ../shared

replace covid19/models => ../models

module daily-data-batch

go 1.16

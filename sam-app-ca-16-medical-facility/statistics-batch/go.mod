replace covid19/shared => ../shared

replace covid19/models => ../models

module statistics-batch

go 1.17

require (
	covid19/models v0.0.0-00010101000000-000000000000
	covid19/shared v0.0.0-00010101000000-000000000000
	github.com/shopspring/decimal v1.3.1
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

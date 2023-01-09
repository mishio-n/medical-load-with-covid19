module past-data-batch

go 1.17

replace covid19/shared => ../shared

replace covid19/models => ../models

require (
	covid19/models v0.0.0-00010101000000-000000000000
	covid19/shared v0.0.0-00010101000000-000000000000
)

require github.com/go-sql-driver/mysql v1.7.0 // indirect

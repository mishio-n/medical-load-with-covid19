replace covid19/shared => ../shared

module statistics-batch

go 1.17

require (
	covid19/shared v0.0.0-00010101000000-000000000000
	github.com/shopspring/decimal v1.3.1
)

require github.com/go-sql-driver/mysql v1.7.0 // indirect

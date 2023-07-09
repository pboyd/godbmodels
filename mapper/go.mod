module github.com/pboyd/godbmodels/mapper

go 1.20

require (
	github.com/jmoiron/sqlx v1.3.5
	github.com/pboyd/godbmodels/common v0.0.0
	github.com/stretchr/testify v1.8.4
)

replace github.com/pboyd/godbmodels/common => ../common

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

module github.com/pboyd/godbmodels/orm

go 1.20

require (
	github.com/pboyd/godbmodels/common v0.0.0
	github.com/stretchr/testify v1.8.4
	gorm.io/driver/sqlite v1.5.2
	gorm.io/gorm v1.25.2
)

replace github.com/pboyd/godbmodels/common => ../common

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

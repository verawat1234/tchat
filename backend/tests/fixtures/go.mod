module tchat.dev/tests/fixtures

go 1.24.3

require (
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	tchat.dev/auth v0.0.0
	tchat.dev/content v0.0.0
	tchat.dev/payment v0.0.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	golang.org/x/text v0.27.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/gorm v1.31.0 // indirect
)

replace tchat.dev/auth => ../../auth

replace tchat.dev/content => ../../content

replace tchat.dev/payment => ../../payment

replace tchat.dev/shared => ../../shared

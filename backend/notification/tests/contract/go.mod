module tchat.dev/notification/tests/contract

go 1.23

require (
	github.com/gin-gonic/gin v1.11.0
	github.com/google/uuid v1.6.0
	github.com/pact-foundation/pact-go/v2 v2.4.1
	github.com/stretchr/testify v1.11.1
	tchat.dev/notification v0.0.0
	tchat.dev/shared v0.0.0
)

replace tchat.dev/notification => ../../

replace tchat.dev/shared => ../../../shared

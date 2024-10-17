module github.com/circulohealth/sonar-backend/packages/common

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/aws/aws-lambda-go v1.28.0
	github.com/aws/aws-sdk-go v1.40.20
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.0.0
	github.com/lestrrat-go/iter v1.0.1
	github.com/lestrrat-go/jwx v1.2.6
	github.com/lib/pq v1.10.2 // indirect
	github.com/stretchr/testify v1.7.0
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.20.0
	gorm.io/driver/postgres v1.2.3
	gorm.io/gorm v1.22.5
)

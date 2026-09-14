# Setup Database Migration

## Install Golang Migrate
Run the following command to install golang-migrate
```
brew install golang-migrate
```

## Create db/migration folder
In the project directory, run the following command create a db/migration sub-directory
```
mkdir -p db/migration
```

## Create a migration file
In the project directory, run the following command to create a migration file
```
migrate create -ext sql -dir db/migration -seq init_schema
```

## Upgrade database
```
migrate -path db/migration -database "postgresql://postgres_user:postgres_password@localhost:5432/simple_bank?sslmode=disable" -verbose up 
```

## Installing SQL-C
```
brew install sqlc
```

## Immport PostgreSQL driver for Go's database/sql package
```
go get github./lib/pq
```

## Immport Testify library
```
go get github.com/stretchr/testify
```

## Installing Gin
```
go get -u github.com/gin-gonic/gin
```

## Installing Viper
```
go get github.com/spf13/viper
```

## Installing Go Mock
```
go install github.com/golang/mock/mockgen@v1.6.0
```

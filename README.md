# TWReporter's Golang Backend API

## Environment 
### Development

#### Go module
After go-api@5.0.0 is released, go-api no longer needs to be developed within $GOPATH/src directory thanks to the go module support. Make sure your go version is go@1.11 or higher to have full compatibility. You can clone to anywhere outside the $GOPATH as you wish.

```golang
$ git clone github.com/twreporter/go-api
$ cd go-api

// Run test
$ go test ./...

// Use makefile
make start
// Or
$ go run main.go

// Build server binaries
$ go build -o go-api
$ ./go-api
```

#### Deprecated
Please make sure that you install [Glide
  package manager](https://github.com/Masterminds/glide) in the environment. (Switch to [go module](https://github.com/golang/go/wiki/Modules) after v5.0.0)

```
cd $GOPATH/src/github.com/twreporter/go-api
glide install                           # Install packages and dependencies

// use Makefile
make start 
// or 
go run main.go                          # Run without live-reloading
```

### Production
```
go build
./go-api
```

## Dependencies Setup and Configurations
There are two major dependencies of go-api, one is MySQL database, 
another is MongoDB. <br/>
MySQL DB stores membership data, which is related to users.<br/>
MongoDB stores news entities, which is the content that go-api provides.<br/>

### Install docker-compose
[docker-compose installation](https://docs.docker.com/compose/install/) 

### Start/Stop MySQL and MongoDB with default settings
```
// start MySQL and MongoDB
make env-up

// stop MySQL and MongoDB
make env-down
```

### Configure MySQL Connection
Copy `configs/config.example.json` and rename as `configs/config.json`.
Change `DBSettings` fields to connect to your own database, like following example.
```
  "DBSettings": {
    "Name":     "test_membership",
    "User":     "test_membership",
    "Password": "test_membership",
    "Address":  "127.0.0.1",
    "Port":     "3306"
  },
```

### Configure MongoDB Connection
Copy `configs/config.example.json` and rename as `configs/config.json`.
Change `MongoDBSettings` fields to connect to your own database, like following example.
```
  "MongoDBSettings": {
    "URL": "localhost",
    "DBName": "plate",
    "Timeout": 5
  },
```

### AWS SES Setup
Currently the source code sends email through AWS SES,

If you want to send email through your AWS SES, just put your AWS SES config under `~/.aws/credentials`
```
[default]
aws_access_key_id = ${AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${AWS_SECRET_ACCESS_KEY}
```

Otherwise, you have to change the `utils/mail.go` to integrate with your email service.

### Database Migrations
Go-api integrates [go-migrate](https://github.com/golang-migrate/migrate) to do the schema version control. You can follow the instructions to install the cli from [here](https://github.com/golang-migrate/migrate/tree/master/cli).

Basic operations are listed below:

Create

```
# Create migration pair files up/down_000001_$FILE_NAME.$FILE_EXTENSION in $FILE_DIR
migrate create -ext $FILE_EXTENSION -dir $FILE_DIR -seq $FILE_NAME 
```

Up

```
# Upgrade to the latest version in migration directory
migrate -databsae $DATABASE_CONNECTION -path $MIGRATION_DIR up
```

Down (CAUTION)

```
# Remove all the existing versions change
migrate -database $DATABASE_CONNECTION -path $MIGRATION_DIR down
```

Goto

```
# Upgrade/Downgrade to the specific versions
migrate -database $DATABASE_CONNECTION -path $MIGRATION_DIR goto $VERSION
```

more details in [FAQ](https://github.com/golang-migrate/migrate/blob/master/FAQ.md) and [godoc](https://godoc.org/github.com/mattes/migrate)(old repo)



## Functional Testing
### Prerequisite
* Make sure the environment you run the test has a running `MySQL` server and `MongoDB` server<br/>

### How To Run Tests
```
// use Makefile
make test

// or

go test $(glide novendor)

// or print logs
go test -v $(glide novendor)
```

## API Document
See [https://twreporter.github.io/go-api/](https://twreporter.github.io/go-api/)

## License
Go-api is [MIT licensed](https://github.com/twreporter/go-api/blob/master/LICENSE)

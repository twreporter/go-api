# TWReporter's Golang Backend API

## Development
Please make sure that [Glide
  package manager](https://github.com/Masterminds/glide) was installed in your environment.

```
cd $GOPATH/src/twreporter.org/go-api
glide install                           # Install packages and dependencies
go run main.go
```

## Production
```
go build
./go-api
```


## Testing
```
$ go test $(glide novendor)             # run go test over all directories of the project except the vendor directory
```

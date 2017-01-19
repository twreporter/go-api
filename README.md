# TWReporter's Golang Backend API

## Development
Please make sure that you install [Glide
  package manager](https://github.com/Masterminds/glide) in the environment.

```
cd $GOPATH/src/twreporter.org/go-api
glide install                           # Install packages and dependencies
go run main.go                          # Run without live-reloading
```

#### Live Reloading
Note that `GOPATH/bin` should be in your `PATH`.
```
go get github.com/codegangsta/gin
gin                                     # Run with live-reloading
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

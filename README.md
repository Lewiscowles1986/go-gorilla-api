# go-gorilla-api
testing semaphore ci with basic golang

## building

```
go build main.go app.go
```

## testing

```
go test -v . ./data ./repositories ./rest ./settings -cover -coverprofile cover.out
```

## coverage

```
go tool cover -html=./cover.out
```

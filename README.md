# go-gorilla-api
testing semaphore ci with basic golang

## building

```
go build -o main .
```

## testing

### Unit

```
go test -v . ./data ./repositories ./rest ./settings -cover -coverprofile cover.out
```

### Integration

> :warning: **This uses brine-dsl**: This package is known to currently not work with Ruby v3 and (oddly) requires the ruby runtime installed.

<details>
<summary>Setting up Brine</summary>

- install [Ruby Version Manager](https://rvm.io/)
- install Ruby runtime `2.7.6` (this is just the latest 2.x as 3.x is known broken)
- (from within this folder) using cli `bundle install`

</details>

```
# Run our app
./main &
# supply base_url and run cucumber
BRINE_ROOT_URL=http://localhost:8080 cucumber -vvv
```

## coverage

```
go tool cover -html=./cover.out
```

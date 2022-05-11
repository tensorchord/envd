# Contributing

## Community Code of Conduct

Please refer to the [CNCF Community Code of Conduct v1.0](https://github.com/cncf/foundation/blob/main/code-of-conduct.md)

## Development

### Requirements

- Go 1.17
- Docker

### Build from source

```
go mod tidy
make
./bin/midi bootstrap
./bin/midi --version
```

### Running tests

To run tests you can issue

```
make test
```

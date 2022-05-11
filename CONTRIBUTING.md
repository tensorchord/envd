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

### Lint

You could run the command below

```
make lint
```

You should see output similar to the following if there is any linting issue:

```
cmd/envd/main.go:36:67: Revision not declared by package version (typecheck)
                fmt.Println(c.App.Name, version.Package, c.App.Version, version.Revision)
                                                                                ^
make: *** [Makefile:102: lint] Error 1
```

### Running tests

To run tests you can issue

```
make test
```

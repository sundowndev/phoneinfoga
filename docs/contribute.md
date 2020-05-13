# Contribute

This page describe the project structure and gives you a bit of help to start contributing.

The project is maintained by a single person: [sundowndev](https://github.com/sundowndev). Contributions are welcome !

!!! tip "Want to contribute ? Clone the project and open some pull requests !"

## Project

### Installation

See the [installation page](install.md) to install the project.

### File structure

```shell
api         # REST API code
client      # web client code
cmd         # Command-line app code
docs        # Documentation
pkg         # Code base for scanners, utils ...
scripts     # Development & deployment scripts
go.mod      # Go modules file
main.go     # Application entrypoint
```

## Testing

### Go code

```shell
go test -v ./...

# Collect coverage
go test -coverprofile=coverage.out ./...

# Open coverage file as HTML
go tool cover -html=coverage.out
```

### Typescript code

Developping on the web client.

```shell
cd client

yarn test
yarn test:unit
yarn test:e2e
```

## Formatting

### Go code

We use a shell script to format Go files.

```shell
sh ./scripts/format.sh

# You can also use GolangCI
golangci-lint run -D errcheck
```

### Typescript code

```shell
cd client

yarn lint
yarn lint:fix
```

## Documentation

We use [mkdocs](https://www.mkdocs.org/) to write our documentation.

### Install mkdocs

```
python3 -m pip install mkdocs
```

### Serve documentation on localhost

```
mkdocs serve
```

### Deploy on github pages

```
mkdocs gh-deploy
```
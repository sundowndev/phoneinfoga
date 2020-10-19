# Contribute

This page describe the project structure and gives you a bit of help to start contributing.

The project is maintained by a single person: [sundowndev](https://github.com/sundowndev). Contributions are welcome !

!!! tip "Want to contribute ? Clone the project and open some pull requests !"

## Project

### Building from source

**Requirements :**

- Node.js >= v10.x
- npm or yarn
- Go >= 1.13

**Note:** if you're using npm, just replace `yarn <command>` by `npm run <command>`.

```shell
# Build static assets
# This will create dist directory containing client's static files
$ (cd client && yarn && yarn build)

# Generate in-memory assets
# This will put content of dist directory in memory. It's usually needed to build but
# the design requires you to do it anyway.
# This step is needed at each change if you're developing on the client.
$ go generate ./...

# Build the whole project
$ go build -v .
```

If you're developing, you don't need to build at each changes, you can compile then run with the `go run` command :

```
$ go run main.go
```

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
# Run test suite
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

If you're developing on the client, you can watch changes with `yarn build:watch`.

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

```shell
python3 -m pip install mkdocs
```

### Serve documentation on localhost

This is the only command you need to start working on docs.

```shell
mkdocs serve
# or
python3 -m mkdocs serve
```

### Build website

```shell
mkdocs build
```

### Deploy on github pages

```shell
mkdocs gh-deploy
```

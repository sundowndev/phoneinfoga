---
hide:
- navigation
---

# Contribute

This page describe the project structure and gives you a bit of context to start contributing to the project.

## Project

### Building from source

**Requirements :**

- Nodejs >= v15
- npm or yarn
- Go >= 1.16

**Note:** if you're using npm, just replace `yarn <command>` by `npm run <command>`.

```shell
# Install tools needed to build, creating mocks or running tests
$ make install-tools

# Build static assets
# This will create dist directory containing client's static files
$ (cd web/client && yarn && yarn build)

# Generate in-memory assets, then build the project.
# This will put content of dist directory in a single binary file.
# It's needed to build but the design requires you to do it anyway.
# This step is needed at each change if you're developing on the client.
$ make build
```

If you're developing, you don't need to build at each changes, you can compile then run with the `go run` command :

```
$ go run main.go
```

### File structure

```shell
bin/        # Local binaries
build/      # Build package providing info about the current build
cmd/        # Command-line app code
docs/       # Documentation
examples/   # Some code examples
lib/        # Libraries 
mocks/      # Mocks
support/    # Utilities, manifests for production and local env
test/       # Utilities for testing purposes
web/        # Web server, including REST API and web client
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
cd web/client

yarn test
yarn test:unit
yarn test:e2e
```

If you're developing on the client, you can watch changes with `yarn build:watch`.

## Formatting

### Go code

We use gofmt to format Go files.

```shell
make fmt
```

### Typescript code

```shell
cd web/client

yarn lint
yarn lint:fix
```

## Documentation

We use [mkdocs](https://www.mkdocs.org/) to generate our documentation website.

### Install mkdocs

```shell
python3 -m pip install mkdocs==1.3.0 mkdocs-material==8.3.9 mkdocs-minify-plugin==0.5.0 mkdocs-redirects==1.1.0
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

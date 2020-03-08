# Contribute

This page describe the project structure and gives you a bit of help to start contributing.

The project is maintained by a single person: [sundowndev](https://github.com/sundowndev). Contributions are welcome !

!!! tip "Want to contribute ? Clone the project and open some pull requests !"

## Project

### Installation

See the [installation page](install.md) to install the project.

### File structure

```
$ tree . -I __pycache__

├── docs
├── examples
│   ├── generate.sh
│   ├── input.txt
│   ├── output_from_input.txt
│   └── output_single.txt
├── lib
│   ├── args.py
│   ├── banner.py
│   ├── colors.py
│   ├── format.py
│   ├── googlesearch.py
│   ├── __init__.py
│   ├── logger.py
│   ├── output.py
│   └── request.py
├── osint
├── scanners
│   ├── footprints.py
│   ├── __init__.py
│   ├── localscan.py
│   ├── numverify.py
│   ├── ovh.py
│   └── recon.py
├── config.example.py
├── Dockerfile
├── mkdocs.yml
├── phoneinfoga.py
└── requirements.txt
```

## Testing

### Go code

```shell
go test -v ./...

# Collect coverage
go test -coverprofile=coverage.out ./...

# Open coverage as HTML
go tool cover -html=coverage.out
```

### Typescript code

```shell
cd client

yarn test
yarn test:unit
yarn test:e2e
```

## Formatting

### Go code

We use a custom script to format Go files.

```shell
sh ./scripts/format.sh
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
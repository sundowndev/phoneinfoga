# Contribute

This page describe the project structure and gives you a bit of help to start developing.

The project is maintained by a single person: [sundowndev](https://github.com/sundowndev). Contributions are welcome.

You want to contribute ? Clone the project and open some pull requests !

## Project

### Installation

See the [installation page](install.md) to install the project.

### Structure

```
├── docs
│   ├── contribute.md
│   ├── formatting.md
│   ├── googlesearch.md
│   ├── index.md
│   ├── install.md
│   ├── resources.md
│   └── usage.md
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
│   ├── disposable_num_providers.json
│   ├── individuals.json
│   ├── reputation.json
│   └── social_medias.json
├── scanners
│   ├── footprints.py
│   ├── __init__.py
│   ├── localscan.py
│   ├── numverify.py
│   ├── ovh.py
│   └── recon.py
├── config.example.py
├── Dockerfile
├── LICENSE
├── mkdocs.yml
├── phoneinfoga.py
├── README.md
└── requirements.txt
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
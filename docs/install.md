To install PhoneInfoga, you'll need to download the binary or build the software from source code.

## Download binary

Follow the instructions :

```shell
# Using curl
# Get the latest version
PHONEINFOGA_VERSION=$(curl -s https://api.github.com/repos/sundowndev/phoneinfoga/releases/latest | grep tag_name | cut -d '"' -f 4)

curl -sSL "https://github.com/sundowndev/phoneinfoga/releases/download/$PHONEINFOGA_VERSION/phoneinfoga_$(uname -s)_$(uname -m).tar.gz" -o ./phoneinfoga.tar.gz
tar xfv phoneinfoga.tar.gz

# Run the software
./phoneinfoga version

# Install it globally
mv ./phoneinfoga /usr/bin/phoneinfoga
```

## Install globally with Go

```shell
go get -u github.com/sundowndev/PhoneInfoga
```

## Build from source

Follow the instructions :

```shell
# Clone the repository
git clone https://github.com/sundowndev/PhoneInfoga
cd PhoneInfoga/

# Install requirements
go get -v -t -d ./...

# You need Packr to inject assets inside the binary
packr build -o phoneinfoga

./phoneinfoga
```

## Using Docker

### From docker hub

You can pull the repository directly from Docker hub

```shell
docker pull sundowndev/phoneinfoga:latest
```

Then run the tool

```shell
docker run --rm -it sundowndev/phoneinfoga version
```

### Docker-compose

You can use a single docker-compose file to run the tool without downloading the source code.

```
version: '3.7'

services:
    phoneinfoga:
      container_name: phoneinfoga
      restart: on-failure
      build:
        context: .
        dockerfile: Dockerfile
      command:
        - "serve"
      ports:
        - "80:5000"
```

### From the source code

You can download the source code, then build the docker images

#### Build

This will automatically pull, build then setup services (such as web client and REST API)

```shell
docker-compose up -d
```

#### CLI usage

```shell
docker-compose run --rm phoneinfoga --help
```

#### Troubleshooting

All output is sent to stdout so it can be inspected by running:

```shell
docker logs -f <container-id|container-name>
```

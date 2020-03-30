To install PhoneInfoga, you'll need to download the binary or build the software from its source code.

!!! info
    For now, only Linux and MacOS are supported. If you don't see your OS/arch on the [release page on GitHub](https://github.com/sundowndev/PhoneInfoga/releases), it means it's not explicitly supported. You can always build from source by yourself. Want your OS to be supported ? Please [open an issue on GitHub](https://github.com/sundowndev/PhoneInfoga/issues).

## Download the binary

Follow the instructions :

- Go to [release page on GitHub](https://github.com/sundowndev/PhoneInfoga/releases)
- Choose your OS and architecture
- Download the archive, extract the binary then run it in a terminal

You can also do it from the terminal:

```shell
LATEST_VERSION=$(curl -s https://api.github.com/repos/sundowndev/phoneinfoga/releases/latest | grep tag_name | cut -d '"' -f 4)

# Download the archive
curl -sSL "https://github.com/sundowndev/phoneinfoga/releases/download/$LATEST_VERSION/phoneinfoga_$(uname -s)_$(uname -m).tar.gz" -o ./phoneinfoga.tar.gz

# Extract the binary
tar xfv phoneinfoga.tar.gz

# Run the software
./PhoneInfoga --help

# Install it globally
mv ./PhoneInfoga /usr/bin/phoneinfoga
```

## Building from source

You shouldn't need to, but in case the binary isn't working you can still build the software from source.

Follow the instructions :

```shell
# Clone the repository
git clone https://github.com/sundowndev/PhoneInfoga
cd PhoneInfoga/

# Install requirements
go get -v -t -d ./...

# Install packr2
go get -u github.com/gobuffalo/packr/v2/packr2

# Build web client assets
(cd client && yarn && yarn build)

# You need Packr v2 to inject assets inside the binary
packr2 build -o phoneinfoga

packr2 clean

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

Build the image 

```shell
docker-compose build
```

#### CLI usage

```shell
docker-compose run --rm phoneinfoga --help
```

#### Run web services

```shell
docker-compose up -d
```

##### Disable web client

Edit `docker-compose.yml` and add the `--no-client` option

```yaml
# docker-compose.yml
command:
    - "serve --no-client"
```

#### Troubleshooting

All output is sent to stdout so it can be inspected by running:

```shell
docker logs -f <container-id|container-name>
```

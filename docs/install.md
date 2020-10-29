To install PhoneInfoga, you'll need to download the binary or build the software from its source code.

!!! info
    For now, only Linux, MacOS and Windows are supported. If you don't see your OS/arch on the [release page on GitHub](https://github.com/sundowndev/PhoneInfoga/releases), it means it's not explicitly supported. You can build from source by yourself anyway. Want your OS to be supported ? Please [open an issue on GitHub](https://github.com/sundowndev/PhoneInfoga/issues).

## Binary installation (recommanded)

Follow the instructions :

- Go to [release page on GitHub](https://github.com/sundowndev/PhoneInfoga/releases)
- Choose your OS and architecture
- Download the archive, extract the binary then run it in a terminal

You can also do it from the terminal (UNIX systems only) :

```shell
# Download latest release in the current directory
curl -sSL https://raw.githubusercontent.com/sundowndev/PhoneInfoga/master/support/scripts/install | bash

# Check the binary
./phoneinfoga version

# You can also install it globally
sudo mv ./phoneinfoga /usr/bin/phoneinfoga
```

To ensure your system is supported, please check the output of `echo "$(uname -s)_$(uname -m)"` in your terminal and see if it's available on the [GitHub release page](https://github.com/sundowndev/PhoneInfoga/releases).

## Using Docker

!!! info
    Be careful when using `latest` tag, it's updated directly from the master branch. We recommend using [`v2` or `stable` tags](https://hub.docker.com/r/sundowndev/phoneinfoga/tags) to only get release updates.

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
      image: phoneinfoga:latest
      command: serve
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
command: "serve --no-client"
```

#### Troubleshooting

All output is sent to stdout so it can be inspected by running:

```shell
docker logs -f <container-id|container-name>
```

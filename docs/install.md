To install PhoneInfoga, you'll need to download source code then install dependencies.

Requirements : 

- python3 and python3-pip OR Docker
- git OR wget and curl

## Manual installation

### Clone the repository

```shell
git clone https://github.com/sundowndev/PhoneInfoga
cd PhoneInfoga/
```

You can also download the source code archive : 

```shell
wget $(curl -s https://api.github.com/repos/sundowndev/phoneinfoga/releases/latest | grep tarball_url | cut -d '"' -f 4) -O PhoneInfoga.tar.gz
tar -xvzf PhoneInfoga.tar.gz
cd sundowndev*
```

### Install requirements

```shell
python3 -m pip install -r requirements.txt --user
```

### Create the config file

```shell
cp config.example.py config.py 
```

To ensure everything works, use the `-v` option to show the version : 

```shell
python3 phoneinfoga.py -v
```

### Install the Geckodriver

The Geckodriver is the Firefox webdriver for Selenium, which is used by PhoneInfoga to perform queries to Google search and handle captcha. Firefox is actually the only webdriver supported by PhoneInfoga. Want to hack it to use chrome or another driver instead ? See [this file](https://github.com/sundowndev/PhoneInfoga/blob/8179fe4857ca7df2d843119e2123c260e8401818/lib/googlesearch.py#L35).

#### Linux

##### Download

Go to the [geckodriver releases page](https://github.com/mozilla/geckodriver/releases). Find the latest version of the driver for your platform and download it. For example: 

```shell
wget https://github.com/mozilla/geckodriver/releases/download/v0.24.0/geckodriver-v0.24.0-linux64.tar.gz
```

##### Extract the file

```shell
sudo sh -c 'tar -x geckodriver -zf geckodriver-*.tar.gz -O > /usr/bin/geckodriver'
```

##### Make it executable

```shell
sudo chmod +x /usr/bin/geckodriver
```

##### Remove the archive

```shell
rm geckodriver-*.tar.gz
```

**NOTE:** You also have to install Firefox browser v65+. To verify everything is fine, be sure the following commands work:

- `which firefox` should return something like `/usr/bin/firefox`
- `which geckodriver` should return something like `/usr/bin/geckodriver`

#### Windows or MacOS

- Go to the [geckodriver releases page](https://github.com/mozilla/geckodriver/releases). Find the latest version of the driver for your platform and download it.
- Extract the archive
- Run the executable and follow the instructions

## Using Docker

### From docker hub

You can pull the repository directly from Docker hub

```shell
docker pull sundowndev/phoneinfoga:latest
```

Then run the tool

```shell
docker run --rm -it sundowndev/phoneinfoga --help
```
**WARNING**: This image only contain the python tool and not the Selenium hub which is useful to query Google. In order to use Selenium driver, you must use the docker-compose configuration, as described below.

### Docker-compose

You can use a single docker-compose file to run the tool without downloading the source code.

```
version: "3"

services:
  phoneinfoga:
    image: sundowndev/phoneinfoga
    container_name: phoneinfoga
    restart: on-failure
    environment:
      webdriverRemote: 'http://selenium-hub:4444/wd/hub'

  selenium-hub:
    image: selenium/hub:3.141.59-palladium
    container_name: selenium-hub
    ports:
      - "4444:4444"

  firefox:
    image: selenium/node-firefox:3.141.59-palladium
    volumes:
      - /dev/shm:/dev/shm
    depends_on:
      - selenium-hub
    environment:
      - HUB_HOST=selenium-hub
      - HUB_PORT=4444
```

### From the source code

You can download the source code, then build the docker images

#### Build

This will automatically pull, build then setup services

```shell
docker-compose up -d
```

#### Usage

```shell
docker-compose run --rm phoneinfoga --help
```

#### Troubleshooting

All output is sent to stdout so it can be inspected by running:

```shell
docker logs -f <container-id|container-name>
```

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
python3 -m pip install -r requirements.txt
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

#### Linux

Go to the [geckodriver releases page](https://github.com/mozilla/geckodriver/releases). Find the latest version of the driver for your platform and download it. For example: 

```
wget https://github.com/mozilla/geckodriver/releases/download/v0.24.0/geckodriver-v0.24.0-linux64.tar.gz
```

Extract the file with:

```
tar xvfz geckodriver-v0.24.0-linux64.tar.gz
```

Make it executable:

```
chmod +x geckodriver
```

Add the driver to your PATH so other tools can find it:

```
export PATH=$PATH:/path-to-extracted-file/.
```

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

### From the source code

Or, you can download the source code, then build the docker image

#### Build

```shell
docker build --rm=true -t phoneinfoga/latest .
```

#### Usage

```shell
docker run --rm -it phoneinfoga/latest --help
```
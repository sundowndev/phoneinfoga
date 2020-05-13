#!/usr/bin sh

# This script is used to push the image to Docker hub
docker login
docker build --rm=true -t sundowndev/phoneinfoga:latest .
docker tag $(docker images -q sundowndev/phoneinfoga) sundowndev/phoneinfoga:latest
docker push sundowndev/phoneinfoga

echo 'Script executed.'
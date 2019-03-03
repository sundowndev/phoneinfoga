FROM jfloff/alpine-python:3.6-onbuild

ADD . /opt/phoneinfoga

WORKDIR /opt/phoneinfoga

ENTRYPOINT ["python", "phoneinfoga.py"]

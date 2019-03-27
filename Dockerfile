FROM jfloff/alpine-python:3.6-onbuild

ADD . /opt/phoneinfoga

COPY ./config.example.py /opt/phoneinfoga/config.py

WORKDIR /opt/phoneinfoga

ENTRYPOINT ["python", "phoneinfoga.py"]

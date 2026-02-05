FROM ubuntu:latest
LABEL authors="yu"

ENTRYPOINT ["top", "-b"]
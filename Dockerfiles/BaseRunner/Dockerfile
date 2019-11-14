# BUILD-USING:    docker build -t code-kombat/base-runner .
# RUN-USING:      docker run --rm --name=runner-base code-kombat/base-runner --help

# Pull base image.
FROM ubuntu:18.04

ENV DEBIAN_FRONTEND noninteractive

# create non root user
RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
RUN mkdir /home/appuser
RUN chown -R appuser:appuser /home/appuser/

## install php libraries

USER appuser
WORKDIR /home/appuser
RUN composer require symfony/process
USER root

# Pull base image.
FROM buildpack-deps:bionic-scm

RUN apt-get update -qq

RUN apt-get install -y software-properties-common

# create non root user
RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
RUN mkdir /home/appuser
RUN chown -R appuser:appuser /home/appuser/

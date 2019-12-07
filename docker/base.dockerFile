# Pull base image.
FROM buildpack-deps:bionic-scm

RUN RUN apt-get update -qq

# create non root user
RUN groupadd -g 999 appuser && \
    useradd -r -u 999 -g appuser appuser
RUN mkdir /home/appuser
RUN chown -R appuser:appuser /home/appuser/

# pull base-image
FROM code-kombat/base-runner

ENV DEBIAN_FRONTEND noninteractive

## install pip
RUN apt-get update && apt-get install -yq --no-install-recommends python-pip\
    python3-pip\
    python3-setuptools\
    python-dev\
    build-essential\
    libssl-dev\
    libffi-dev

## install python packages
RUN pip3 install timeout-decorator

COPY test_frameworks /usr/lib/python3.6/code_kombat_test_frameworks

USER appuser


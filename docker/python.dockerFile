# pull base-image
FROM sandman/base-runner

## install pip
RUN apt-get update && apt-get install -yq --no-install-recommends python-pip\
    python3-pip\
    python3-setuptools\
    python-dev\
    build-essential\
    libssl-dev\
    libffi-dev

USER appuser
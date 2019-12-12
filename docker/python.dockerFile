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

COPY entrypoint.sh entrypoint.sh
COPY language_conf.json language_conf.json
COPY sandman sandman
RUN chmod +x entrypoint.sh
RUN chmod +x sandman

USER appuser
ENV USER=appuser HOME=/home/appuser

ENTRYPOINT ["./entrypoint.sh"]

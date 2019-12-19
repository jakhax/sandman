# pull base-image
FROM sandman/base-runner

RUN apt-get update -qq

RUN apt-get install -y build-essential\
    checkinstall\
    libssl-dev\
    libffi-dev\
    python3-pip\
    python3-setuptools\
    python3-dev


## install python packages
RUN pip3 install timeout-decorator

# set  workdir
WORKDIR /home/appuser

COPY entrypoint.sh entrypoint.sh
COPY language_conf.json language_conf.json
COPY sandman sandman
RUN chmod +x entrypoint.sh
RUN chmod +x sandman

# language test framework
COPY frameworks/python3/ /usr/lib/python3.6/

USER appuser
ENV USER=appuser HOME=/home/appuser

ENTRYPOINT ["./entrypoint.sh"]

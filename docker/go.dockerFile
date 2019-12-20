# pull base-image
FROM sandman/base-runner

RUN apt-get update -qq

ENV GOLANG_VERSION 1.13.5

ENV GOLANG_DOWNLOAD_URL https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 512103d7ad296467814a6e3f635631bd35574cab3369a97a323c9a585ccaa569

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /home/appuser/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN go get github.com/onsi/ginkgo/ginkgo  # installs the ginkgo CLI
RUN go get github.com/onsi/gomega         # fetches the matcher library

COPY frameworks/golang/codewars $GOPATH/src/codewars
RUN go install codewars/reporter

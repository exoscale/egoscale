FROM golang:1.10-stretch

ARG DEBIAN_FRONTEND=noninteractive

RUN go get -u github.com/golang/dep/cmd/dep \
 && go get -u -d github.com/goreleaser/goreleaser/... \
 && go get -u -d github.com/goreleaser/nfpm/... \
 && apt-get update -q \
 && apt-get upgrade -qy \
 && apt-get install -qy \
        rpm \
 && cd $GOPATH/src/github.com/goreleaser/nfpm \
 && dep ensure -v -vendor-only \
 && go install \
 && cd ../goreleaser \
 && dep ensure -v -vendor-only \
 && go install \
 && cd /

ADD ops.asc ops.asc
RUN gpg --allow-secret-key-import --import ops.asc

VOLUME /go/src/github.com/exoscale/egoscale
WORKDIR /go/src/github.com/exoscale/egoscale

CMD ['goreleaser', '--snapshot']

FROM golang:latest
ENV GOPATH /go
WORKDIR /go/src/github.com/ltacker
RUN git clone https://github.com/ltacker/request-chain.git
WORKDIR /go/src/github.com/ltacker/request-chain
RUN make install

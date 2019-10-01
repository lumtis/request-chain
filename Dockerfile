FROM golang:latest
ENV GOPATH /go
WORKDIR /go/src/github.com/ltacker
RUN git clone https://github.com/ltacker/request-chain.git
WORKDIR /go/src/github.com/ltacker/request-chain
RUN make install
RUN rcd init tacker --chain-id wacken
ADD utils/genesis.json /root/.rcd/config/
RUN rcd unsafe-reset-all
ENTRYPOINT ["rcd", "start"]

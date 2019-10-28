FROM golang:latest
ENV GOPATH /go
WORKDIR /go/src/github.com/ltacker
RUN git clone https://github.com/ltacker/request-chain.git
WORKDIR /go/src/github.com/ltacker/request-chain
RUN make install
RUN rcd init tacker --chain-id wacken
ADD utils/genesis.json /root/.rcd/config/
ADD utils/priv_validator_key.json /root/.rcd/config/
RUN rcd unsafe-reset-all
RUN rccli config output json
RUN rccli config indent true
RUN rccli config trust-node true
RUN rccli config chain-id wacken
RUN rccli config node rcd:26657

FROM golang:latest
ENV GOPATH /go
WORKDIR /go/src/github.com/ltacker
RUN git clone https://github.com/ltacker/request-chain.git
WORKDIR /go/src/github.com/ltacker/request-chain

RUN make install
RUN apt-get update -y && apt-get install -y expect
RUN rcd init tacker --chain-id wacken

# Import genesis state into daemon
ADD utils/genesis.json /root/.rcd/config/
ADD utils/priv_validator_key.json /root/.rcd/config/
RUN rcd unsafe-reset-all

# Set config for the cli
RUN rccli config output json
RUN rccli config indent true
RUN rccli config trust-node true
RUN rccli config chain-id wacken
RUN rccli config node rcd:26657

# Import the private key into cli to sign transactions
ADD utils/keysfilePierre /go/src/github.com/ltacker/request-chain
ADD utils/importPierre /go/src/github.com/ltacker/request-chain
ADD utils/keysfilePerrine /go/src/github.com/ltacker/request-chain
ADD utils/importPerrine /go/src/github.com/ltacker/request-chain
RUN chmod +x importPierre
RUN chmod +x importPerrine
RUN /usr/bin/expect -f importPierre
RUN /usr/bin/expect -f importPerrine

# ENTRYPOINT ["rcd", "start"]

FROM golang:1.14.4-buster as builder

RUN apt-get update && apt-get install ruby vim-common -y

RUN apt-get install flex bison -y
RUN wget http://www.tcpdump.org/release/libpcap-1.9.1.tar.gz && tar xzf libpcap-1.9.1.tar.gz && cd libpcap-1.9.1 && ./configure && make install

#RUN go get github.com/google/gopacket
#RUN go get -u golang.org/x/lint/golint

#WORKDIR /go/src/github.com/buger/goreplay/
#ADD . /go/src/github.com/buger/goreplay/

#RUN go get


WORKDIR /app

COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.io,direct
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -ldflags="-s -w -extldflags \"-static\""

FROM alpine:3.12

COPY --from=builder /app/goreplay .

ENTRYPOINT ["./goreplay"]

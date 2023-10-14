FROM golang:1.19.1 as builder

ENV GO111MODULE=on \
      CGO_ENABLED=0 \
      GOOS=linux \
      GOARCH=amd64
ENV GOPROXY=https://goproxy.cn

WORKDIR /gogin
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o test .


FROM scratch

COPY ./conf /conf
COPY ./docs /docs
COPY --from=builder /gogin/test /

EXPOSE 9090
CMD ["/test","conf/conf.ini"]
FROM golang

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY . .
# RUN go mod download
RUN go build -o bin .

CMD ["/build/bin"]
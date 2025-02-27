FROM golang:latest

WORKDIR /pb/src

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
--mount=type=cache,target=$GOPATH/pkg/mod \
    go build -v -o /pb/app

EXPOSE 8080
WORKDIR /pb
# start PocketBase

CMD ["/pb/app", "serve", "--http=0.0.0.0:8080","--dev"]

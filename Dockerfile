FROM golang:latest

WORKDIR /pb/src
RUN --mount=type=cache,target=/root/.cache/go-build 
RUN --mount=type=cache,target=$GOPATH/pkg/mod 
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./

COPY . .
RUN go build -v -o /pb/app

EXPOSE 8080
WORKDIR /pb
# start PocketBase

CMD ["/pb/app", "serve", "--http=0.0.0.0:8080","--dev"]

FROM golang:1.21

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/peer-z && mkdir -p $HOME/.config

EXPOSE 33300:33300

CMD ["peer-z", "-d", "-noupnp"]

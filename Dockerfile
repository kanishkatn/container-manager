FROM docker:dind

# Install Go
ENV GO_VERSION=1.21.1
RUN apk add --no-cache curl git && \
    curl -LO https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Set up Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# create a directory, copy the code and build the app
WORKDIR /container-manager
COPY . .

RUN go mod download
RUN go build -o container-manager .

CMD ["./container-manager"]
FROM 

WORKDIR /app

ARG SERVICE_NAME=go-clean-architecture
RUN mkdir -p bin && mkdir -p ${SERVICE_NAME}

# Download lib
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Build file binary
WORKDIR /app/${SERVICE_NAME}
COPY . ./
WORKDIR /app/${SERVICE_NAME}/cmd
RUN go build -o /app/bin/main .
WORKDIR /app
RUN rm -rf ${SERVICE_NAME}/

EXPOSE 8088
EXPOSE 8089

ENV GOOS linux
ENV CGO_ENABLED 0

CMD ["/app/bin/main"] 
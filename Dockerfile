FROM golang:latest as builder
WORKDIR /build/
COPY go.mod go.sum ./
RUN go mod download
ENV CGO_ENABLE=0
ENV GOOS=linux 
ENV GOARCH=amd64
COPY . ./
RUN go build 

FROM ubuntu:latest
WORKDIR /app/
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean
COPY --from=builder /build/rabbithole2 /bin/rabbithole
RUN chmod +x /bin/rabbithole
CMD /bin/rabbithole
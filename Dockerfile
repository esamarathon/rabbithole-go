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
COPY --from=builder /build/rabbithole2 /bin/rabbithole
RUN chmod +x /bin/rabbithole
WORKDIR /rabbithole/
CMD /bin/rabbithole
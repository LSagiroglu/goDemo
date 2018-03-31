FROM golang:1.10-alpine as buildImage
ARG VERSION=v1.0
RUN apk add --no-cache --update git openssl ca-certificates  && \	
    update-ca-certificates 

COPY . /go/src/github.com/lsagiroglu/godemo
WORKDIR /go/src/github.com/lsagiroglu/godemo
RUN go get ./...
WORKDIR /go/src/github.com/lsagiroglu/godemo/cmd/godemo
RUN env CGO_ENABLED=0 go build -a -installsuffix cgo -o godemo

FROM scratch
LABEL maintainer "Levent SAGIROGLU <LSagiroglu@gmail.com>"
ENV DOMAIN "" 
ENV EMAIL "" 
COPY --from=buildImage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=buildImage /go/src/github.com/lsagiroglu/godemo/cmd/godemo/godemo /bin/godemo
EXPOSE 80 443
ENTRYPOINT ["/bin/godemo"]

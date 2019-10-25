FROM golang:1.8

WORKDIR /go/src/app
COPY ./api/api .

CMD ["/go/src/app/api"]
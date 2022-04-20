FROM golang:1.18 as build-env

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...

RUN go install -v -buildvcs=false ./...

RUN CGO_ENABLED=0 go build -buildvcs=false -o /go/bin/app

FROM gcr.io/distroless/static

COPY --from=build-env /go/bin/app /
CMD ["/app"]
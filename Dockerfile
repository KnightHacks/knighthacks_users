FROM golang:1.20-alpine as build-env

WORKDIR /go/src/app
COPY . .

RUN go mod download

RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -buildvcs=false -o /go/bin/app

FROM gcr.io/distroless/static

COPY --from=build-env /go/bin/app /
CMD ["/app"]
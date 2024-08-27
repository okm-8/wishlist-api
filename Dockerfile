FROM golang:1.22-bookworm AS dev-env

ARG CGO_ENABLED=0

WORKDIR /usr/src/app

COPY ./go.mod ./go.sum /usr/src/app/

RUN go get ./...

VOLUME ["/usr/src/app"]

CMD ["tail", "-f", "/dev/null"]


FROM golang:1.22 AS test

ARG CGO_ENABLED=0

COPY . /usr/src/app

WORKDIR /usr/src/app

CMD [ "make", "clean-test", "test" ]


FROM golang:1.22 AS build

ARG CGO_ENABLED=0

COPY . /usr/src/app

WORKDIR /usr/src/app

RUN make build


FROM scratch AS runtime

COPY --from=build /usr/src/app/bin/wishlist-api /wishlist-api
COPY --from=build /usr/src/app/migrations /migrations

ENTRYPOINT ["/wishlist-api"]
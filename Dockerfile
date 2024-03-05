FROM golang:1.22.0-alpine3.19 as build

WORKDIR /home

COPY . .

RUN export CGO_ENABLED=0; \
    go mod download && \
    go build -a -tags netgo -trimpath -ldflags="-s -w -extldflags \"-static\"" -o /home/main ./*.go

FROM scratch

COPY --from=build /home/main /home/main

EXPOSE 8080

CMD ["/home/main"]
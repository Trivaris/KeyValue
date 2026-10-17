FROM golang:1.26 AS build

COPY go.mod /
COPY cmd /cmd

# ENV GOPATH=$TMPDIR
# ENV GOCACHE=$TMPDIR/go-cache
ENV CGO_ENABLED=0
RUN go build -o /KeyValue /cmd/server

FROM scratch AS install

COPY --from=build --chmod=755 /KeyValue /KeyValue

EXPOSE 8080
ENTRYPOINT [ "/KeyValue" ]

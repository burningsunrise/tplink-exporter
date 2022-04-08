FROM --platform=${BUILDPLATFORM} golang:alpine AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /out/tplink-exporter .

FROM alpine:latest
COPY --from=build /out/tplink-exporter /
ENTRYPOINT [ "/tplink-exporter" ]
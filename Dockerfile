FROM golang:1.16-alpine AS build
ENV CGO_ENABLED=false
WORKDIR /app
COPY . ./
RUN go build
ENTRYPOINT [ "/app/user-registry" ]

FROM alpine
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=build /app/user-registry .
CMD ["user-registry"]
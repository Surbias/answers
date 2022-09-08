From golang:1.18.0-alpine AS builder

RUN apk add --no-cache git

RUN mkdir /app
ADD . /app
WORKDIR /app

# RUN go clean --modcache
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app .

FROM scratch AS final

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app /app

ENTRYPOINT ["/app/answers"]
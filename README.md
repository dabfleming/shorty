# shorty
URL Shortener for a coding exercise

## Development

`go get github.com/dabfleming/shorty/...` then use `docker-compose up` to provide a dev database with a few seeded records, then run shorty and load `http://localhost:8080/`

## Suggested Improvements

- Wrap errors
- Better verify Input URLs
- Graceful shutdown of server on SIGTERM etc
- Structured logging
- Remove port from IP addresses

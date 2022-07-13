# Gateway Proxy

Needed to test out building a reverse proxy in go so here we are.

## Usage

    source .env
    go run servicea.go &
    go run serviceb.go &
    go run gateway.go

## Accessing

    curl http://localhost:1130/serva
    curl -H "AuthToken: T0K3N" http://localhost:1330/secure/servb

## Notes

No idea if this is done right or if it's useful.

Built it based on this article here: https://hackernoon.com/writing-a-reverse-proxy-in-just-one-line-with-go-c1edfa78c84b

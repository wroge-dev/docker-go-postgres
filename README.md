## Golang Dev Template with Docker

- Serving HTML files + postgres API
- Live reload etc by [Go Develop](https://github.com/zephinzer/golang-dev)!
- You don't need Go installed for this! (docker, docker-compose)

## Usage

- Git Clone outside of GOPATH

- Developing

```
docker-compose up
...open localhost
...do your stuff...

docker-compose down
```

- Build Production Image

```
docker build -t your/image:tag .
```
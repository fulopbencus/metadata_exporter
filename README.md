# metadata_exporter
![Go Report Card](https://goreportcard.com/badge/github.com/fulopbencus/metadata_exporter)

Prometheus metadata exporter written in Go.
It counts the 200's responses in the body trough a HTTP GET method.

# Running containerized

## Docker ![ ](assets/docker.png "Docker icon")

For the pre-built image, pull from ![ ](assets/github_icon.png "Github icon") :
```bash
docker pull fulopbence/metadata_exporter:latest
```

Then run with the following command:

```bash
docker run -d -p 9091:9091 --name=metadata_exporter metadata-exporter -url <your_url>
```

## To build your own image from the Dockerfile

```bash
docker build -t metadata-exporter .
```

# Running source ![ ](assets/go.png "Go icon")

cd inside the library, then:

```bash
go run main.go -port <your_port> -url <your_url>
```

## Testing 

```bash
curl http://localhost:9091/metrics
```

should give you something like this:

```
# HELP metadata_count Number of metadata entries in the host's response
# TYPE metadata_count gauge
metadata_count 6
```
### Flags

| Flag | Default | About |
|---------|-------------|----------------|
| -port | 9091 | The port where the HTTP server is listening|
| -url | empty | The URL's body that gets parsed for 200's endings|

## Requirements

go 1.20

## License

GPL-2.0-or-later

## Auth, support

Made by:

- SZTAKI HBIT

Authors:

- Fülöp Bence <bence.fulop@sztaki.hu>

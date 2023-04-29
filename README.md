# prom-custom-metrics

This commandline utility runs the provided shell commands (check config.yaml file) concurrently in a goroutine and generates the prometheus metrics which can then be read by the node_exporter.

Timeout can be configured using the context

The config file , log file  and the prom file should be writable by the user running the command.

## configuration
- provide the name and the command to run in config.yaml
- path for config file can be configured in type.go (configPath)
- path to generated log file can be configured in type.go (logfilePath)
- path to generated prometheus file can be configured in type.go (promFilePath)

## build
- build locally `go build -o prommetric` 
- build for specific os/architecture `env GOOS=linux GOARCH=amd64 go build  -o prommetric`

## running locally
- run using `go run main.go type.go`
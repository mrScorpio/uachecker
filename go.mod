module github.com/mrscorpio/uachecker

go 1.26

replace uadata => ./uadata

require (
	github.com/influxdata/influxdb v1.13.1
	google.golang.org/grpc v1.84.0
	uadata v0.0.0-00010101000000-000000000000
)

require github.com/go-telegram/bot v1.17.0 // indirect

require (
	github.com/mrscorpio/uahelper v0.0.0-20261022064322-ebf7b4e510f6
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260706201446-f0a921348800 // indirect
	google.golang.org/protobuf v1.36.12 // indirect
)

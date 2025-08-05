module messaging_service

go 1.24.5

require (
	GOGOGO v0.0.0
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
)

replace GOGOGO v0.0.0 => ../../../

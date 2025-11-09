module github.com/AxMdv/go-rocket-factory/inventory

replace github.com/AxMdv/go-rocket-factory/shared => ../shared

go 1.24.8

require (
	github.com/AxMdv/go-rocket-factory/shared v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit/v7 v7.8.1
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

# go-project-template

## Project structure

```
github.com/user/some_project/
├── pkg/ (common own-created packages for all services )
|   ├── errors/
|   ├── log/
|   ├── metrics/
|   ├── sd/
|   |   ├── consul/
|   |   └── kubernetes/
|   └── tracing/
├── internal/ (internal packages)
|   ├── somelib/
├── services/
|   ├── account/
|   |   ├── pb/
|   |   |   ├── account.proto
|   |   |   └── account.pb.go
|   |   ├── handler.go
|   |   ├── main.go
|   |   ├── main_test.go
|   |   ├── Dockerfile
|   |   └── README.md
|   ├── auth/
|   ├── frontend/
|   └── user/
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

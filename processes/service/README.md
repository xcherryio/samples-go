### Command to build this `my_service_mock.go`

Run this at the root of the project:
```shell
mockgen -source=processes/service/my_service.go -package=service -destination=processes/service/my_service_mock.go
```
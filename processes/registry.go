package processes

import (
	"github.com/xdblab/xdb-golang-samples/processes/service"
	"github.com/xdblab/xdb-golang-samples/processes/signup"
	"github.com/xdblab/xdb-golang-sdk/xdb"
)

var registry = xdb.NewRegistry()

func init() {

	svc := service.NewMyService()

	err := registry.AddProcesses(
		signup.NewSignupProcess(svc),
	)
	if err != nil {
		panic(err)
	}
}

func GetRegistry() xdb.Registry {
	return registry
}

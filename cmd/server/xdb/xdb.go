package xdb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-samples/processes"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"log"
	"net/http"
	"sync"
)

// BuildCLI is the main entry point for the iwf server
func BuildCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "xdb golang samples"
	app.Usage = "xdb golang samples"
	app.Version = "beta"

	app.Commands = []cli.Command{
		{
			Name:    "start",
			Aliases: []string{""},
			Usage:   "start xdb golang samples",
			Action:  start,
		},
	}
	return app
}

func start(c *cli.Context) {
	fmt.Println("start running samples")
	closeFn := startWorkflowWorker()
	// TODO improve the waiting with process signal
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	closeFn()
}

var client = xdb.NewClient(processes.GetRegistry(), nil)
var workerService = xdb.NewWorkerService(processes.GetRegistry(), nil)

func startWorkflowWorker() (closeFunc func()) {
	router := gin.Default()
	router.POST(xdb.ApiPathAsyncStateWaitUntil, apiStateWaitUntil)
	router.POST(xdb.ApiPathAsyncStateExecute, apiStateExecute)

	router.GET("/signup/start", signupProcessStart)
	router.GET("/signup/verify", signupProcessVerifyEmail)

	wfServer := &http.Server{
		Addr:    ":" + xdb.DefaultWorkerPort,
		Handler: router,
	}
	go func() {
		if err := wfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return func() { wfServer.Close() }
}

func apiStateWaitUntil(c *gin.Context) {
	var req xdbapi.AsyncStateWaitUntilRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := workerService.HandleAsyncStateWaitUntil(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}
func apiStateExecute(c *gin.Context) {
	var req xdbapi.AsyncStateExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := workerService.HandleAsyncStateExecute(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

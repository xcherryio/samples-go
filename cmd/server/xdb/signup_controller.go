package xdb

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-samples/processes/signup"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"net/http"
)

func signupProcessStart(c *gin.Context) {
	userId := c.Query("userId")
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	email := c.Query("email")
	if userId == "" || firstName == "" || lastName == "" || email == "" {
		c.JSON(http.StatusBadRequest, "must provide userId, firstName, lastName, email via URL parameter")
		return
	}

	form := signup.SubmitForm{
		UserId:    userId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	prc := signup.Process{}
	processExecutionId, err := client.StartProcessWithOptions(
		c.Request.Context(), prc, userId, form, &xdb.ProcessStartOptions{
			GlobalAttributeOptions: xdb.NewGlobalAttributeOptions(
				xdb.DBTableConfig{
					TableName: signup.UserTableName,
					PKValue:   userId,
					InitialAttributes: map[string]interface{}{
						signup.AttrEmail:     email,
						signup.AttrFirstName: firstName,
						signup.AttrLastName:  lastName,
					},
					InitialWriteConflictMode: xdbapi.RETURN_ERROR_ON_CONFLICT.Ptr(),
				},
			),
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("started processId: %v runId: %v", userId, processExecutionId))
	return
}

func signupProcessVerifyEmail(c *gin.Context) {
	userId := c.Query("userId")
	if userId != "" {
		err := client.PublishToLocalQueue(
			context.Background(), userId, signup.VerifyQueue, nil, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		} else {
			c.JSON(http.StatusOK, struct{}{})
		}
		return
	}
	c.JSON(http.StatusBadRequest, "must provide workflowId via URL parameter")
}

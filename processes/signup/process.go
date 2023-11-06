package signup

import (
	"fmt"
	"github.com/xdblab/xdb-apis/goapi/xdbapi"
	"github.com/xdblab/xdb-golang-samples/processes/service"
	"github.com/xdblab/xdb-golang-sdk/xdb"
	"time"
)

func NewSignupProcess(svc service.MyService) xdb.Process {

	return &Process{
		svc: svc,
	}
}

type Process struct {
	xdb.ProcessDefaults

	svc service.MyService
}

func (p Process) GetAsyncStateSchema() xdb.StateSchema {
	return xdb.NewStateSchema(
		&SubmitState{svc: p.svc},
		&VerifyState{svc: p.svc},
	)
}

const (
	VerifyQueue   = "Verify"
	AttrFirstName = "firstName"
	AttrLastName  = "lastName"
	AttrEmail     = "email"
	UserTableName = "sample_user_table"
)

func (p Process) GetPersistenceSchema() xdb.PersistenceSchema {
	return xdb.NewPersistenceSchema(
		xdb.NewEmptyLocalAttributesSchema(),
		xdb.NewGlobalAttributesSchema(
			xdb.NewDBTableSchema(
				UserTableName, "user_id",
				xdbapi.NO_LOCKING,
				xdb.NewDBColumnDef(AttrFirstName, "first_name", true),
				xdb.NewDBColumnDef(AttrLastName, "last_name", true),
				xdb.NewDBColumnDef(AttrEmail, "email", true)),
		),
	)
}

type SubmitState struct {
	svc service.MyService

	xdb.AsyncStateDefaultsSkipWaitUntil
}

func (s SubmitState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var firstName, email string
	persistence.GetGlobalAttribute(AttrFirstName, &firstName)
	persistence.GetGlobalAttribute(AttrEmail, &email)
	err := s.svc.SendEmail(
		email, fmt.Sprintf("%v, please verify your email", firstName),
		".....more content",
	)
	if err != nil {
		return nil, err
	}

	return xdb.SingleNextState(&VerifyState{}, nil), nil
}

type VerifyState struct {
	svc service.MyService

	xdb.AsyncStateDefaults
}

func (v VerifyState) WaitUntil(
	ctx xdb.XdbContext, input xdb.Object, communication xdb.Communication,
) (*xdb.CommandRequest, error) {
	return xdb.AnyOf(
		xdb.NewTimerCommand(time.Second*30),
		xdb.NewLocalQueueCommand(VerifyQueue, 1),
	), nil
}

func (v VerifyState) Execute(
	ctx xdb.XdbContext, input xdb.Object, commandResults xdb.CommandResults, persistence xdb.Persistence,
	communication xdb.Communication,
) (*xdb.StateDecision, error) {
	var firstName, email string
	persistence.GetGlobalAttribute(AttrFirstName, &firstName)
	persistence.GetGlobalAttribute(AttrEmail, &email)

	if commandResults.GetFirstLocalQueueCommand().GetStatus() == xdbapi.COMPLETED_COMMAND {
		err := v.svc.SendEmail(
			email, fmt.Sprintf("%v, welcome!!!!", firstName),
			".....more content",
		)
		if err != nil {
			return nil, err
		}
		return xdb.GracefulCompletingProcess, nil
	}
	err := v.svc.SendEmail(
		email, fmt.Sprintf("%v, REMINDER: verify your email", firstName),
		".....more content",
	)
	if err != nil {
		return nil, err
	}
	return xdb.SingleNextState(&VerifyState{}, nil), nil
}

package signup

type Status string

const (
	StatusPendingVerification Status = "PendingVerification"
	StatusVerified            Status = "Verified"
)

type SubmitForm struct {
	UserId    string
	Email     string
	FirstName string
	LastName  string
}

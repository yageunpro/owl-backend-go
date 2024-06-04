package appointment

type Status string

const (
	DRAFT   Status = "DRAFT"
	CONFIRM Status = "CONFIRM"
	DONE    Status = "DONE"
	CANCEL  Status = "CANCEL"
	DELETE  Status = "DELETE"
)

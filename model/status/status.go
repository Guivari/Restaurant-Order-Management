package status

type Status string

const (
	Pending   Status = "PENDING"
	Confirm   Status = "CONFIRMED"
	Complete  Status = "COMPLETED"
	Cancelled Status = "CANCELLED"
)

var AllStat = []Status{
	Pending,
	Confirm,
	Complete,
	Cancelled,
}

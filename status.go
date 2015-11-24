package zendesk

type Status string

const (
	New     Status = "new"
	Open    Status = "open"
	Pending Status = "pending"
	Hold    Status = "hold"
	Solved  Status = "solved"
	Closed  Status = "closed"
)

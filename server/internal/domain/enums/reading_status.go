package enums

type ReadingStatus string

const (
	Pending   ReadingStatus = "pending"
	Scheduled ReadingStatus = "scheduled"
	Cancelled ReadingStatus = "cancelled"
	Refunded  ReadingStatus = "refunded"
)

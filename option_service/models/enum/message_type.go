package enum

type MessageType int8

const (
	LeaveWord MessageType = 1
	Complaint MessageType = 2
	inquiry   MessageType = 3
	AfterSale MessageType = 4
	AskBuy    MessageType = 5
)

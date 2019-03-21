package alertsv2

type Report struct {
	AckTime        int    `json:"ackTime,omitempty"`
	CloseTime      int    `json:"closeTime,omitempty"`
	AcknowledgedBy string `json:"acknowledgedBy,omitempty"`
	ClosedBy       string `json:"closedBy,omitempty"`
}

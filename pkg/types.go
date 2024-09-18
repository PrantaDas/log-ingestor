package pkg

type Message struct {
	ID             string `json:"requestId"`
	Error          string `json:"error"`
	Source         string `json:"source"`
	AdditionalInfo string `json:"additionInfo"`
}

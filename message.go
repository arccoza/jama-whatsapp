type Chat struct {
	ID string
	Name string
	Type string // Group or Direct
	Protocol string // whatsapp, wechat, google chat, FB messenger
}

type Message struct {
	ID string
	CID string // Chat ID
	Timestamp int64
	From string
	To string
	Text string
	Attachments []Attachment
	Status string // sending, sent, received, read
}

type Attachment stuct {
	Type int
	Mime string
	URL string
	Location [2]int64
}

type Handler func(msg *Message)

package restapi

type Channel struct {
	Name string `json:"chan,omitempty"`
}

type Subscription struct {
	Channel string `json:"chan,omitempty"`
	Name    string `json:"sub,omitempty"`
}

type Message struct {
	Channel      string `json:"chan,omitempty"`
	Subscription string `json:"sub,omitempty"`
	Content      []byte `json:"msg,omitempty"`
}

type Operation struct {
	Status string `json:"status,omitempty"`
	Err    string `json:"err,omitempty"`
}

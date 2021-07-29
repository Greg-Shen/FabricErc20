package event

// 建立事件
type Event struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int    `json:"value"`
}

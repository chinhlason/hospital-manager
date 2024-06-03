package kafka

type Payload interface {
	any
}

type Schema[T Payload] struct {
	Id        int    `json:"id"`
	Key       string `json:"key"`
	Timestamp string `json:"timestamp"`
	Payload   T      `json:"payload"`
}

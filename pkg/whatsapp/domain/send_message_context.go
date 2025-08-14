package domain

type ContextMessage struct {
	Context *Context `json:"context"`
}
type Context struct {
	MessageId string `json:"message_id"`
}

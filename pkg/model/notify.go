package model

type Notify struct {
	Action string
	Body   Request
}

func NewNotify(action string, body Request) *Notify {
	return &Notify{
		Action: action,
		Body:   body,
	}
}

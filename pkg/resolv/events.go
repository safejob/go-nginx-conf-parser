package resolv

type Events struct {
	BasicContext `json:"events"`
}

func NewEvents() *Events {
	return &Events{BasicContext{
		Name:     "events",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}

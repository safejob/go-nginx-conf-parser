package resolv

type Server struct {
	BasicContext
}

func NewServer() *Server {
	return &Server{BasicContext{
		Name:     "server",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}

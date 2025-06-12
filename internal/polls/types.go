package polls

type Poll struct {
	ID      string
	Title   string
	Options []string
	IsOpen  bool
}

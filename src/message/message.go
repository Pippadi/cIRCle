package message

type Message struct {
	From    string
	To      string
	Content string
}

type Command struct {
	Person string
	Action string
	Args   []string
}

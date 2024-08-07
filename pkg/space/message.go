package space

type MessageChannel struct {
	ClassName string `json:"className"`
	ID        string `json:"id"`
}

type MessageContentSectionElement struct {
	ClassName string `json:"className"`
	Content   string `json:"content"`
}

type MessageContentSection struct {
	ClassName string                          `json:"className"`
	Header    string                          `json:"header"`
	Elements  []*MessageContentSectionElement `json:"elements"`
	Footer    string                          `json:"footer"`
}

type MessageContent struct {
	ClassName string                   `json:"className"`
	Style     string                   `json:"style"`
	Sections  []*MessageContentSection `json:"sections"`
}

type ChannelIdentifier struct {
	ClassName string          `json:"className"`
	Channel   *MessageChannel `json:"channel"`
}

type Message struct {
	Channel *ChannelIdentifier `json:"channel"`
	Content *MessageContent    `json:"content"`
}

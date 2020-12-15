package notification

import (
	"fmt"
	"strings"

	"github.com/ns3777k/alertmanager-webhook-space/pkg/alertmanager"
	"github.com/ns3777k/alertmanager-webhook-space/pkg/space"
)

func MapAlertToMessage(alert *alertmanager.Alert, channelID string) *space.Message {
	style := "SUCCESS"
	header := fmt.Sprintf("[**[%s] %s**](%s)",
		strings.ToUpper(alert.Status),
		alert.Annotations.Summary,
		alert.URL,
	)

	if alert.Status == "firing" {
		style = "ERROR"
	}

	return &space.Message{
		Recipient: &space.MessageRecipient{
			ClassName: "MessageRecipient.Channel",
			Channel: &space.MessageChannel{
				ClassName: "ChatChannel.FromId",
				ID:        channelID,
			},
		},
		Content: &space.MessageContent{
			ClassName: "ChatMessage.Block",
			Style:     style,
			Sections: []*space.MessageContentSection{
				{
					ClassName: "MessageSection",
					Header:    header,
					Elements: []*space.MessageContentSectionElement{
						{
							ClassName: "MessageText",
							Content:   alert.Annotations.Description,
						},
					},
				},
			},
		},
	}
}

package main

/*

Discord has a new component v2 api that discordgo doesn't support

Need this flag to enable components v2 for the specific message

`IS_COMPONENTS_V2 = 1 << 15`

*/

const IsComponentsV2 = 1 << 15
const MessageIsEphemeral = 1 << 6

type MessageSend struct {
	Flags      int           `json:"flags"`
	Components []interface{} `json:"components"`
}

type InteractionCallbackType int

const (
	ChannelMessageWithSource         InteractionCallbackType = 4
	DeferredChannelMessageWithSource InteractionCallbackType = 5
	DeferredUpdateMessage            InteractionCallbackType = 6
	Modal                            InteractionCallbackType = 9
)

type InteractionResponse struct {
	Type InteractionCallbackType `json:"type"`
	Data MessageSend             `json:"data"`
}

func NewInteractionResponse(responseType InteractionCallbackType, data MessageSend) *InteractionResponse {
	switch responseType {
	case ChannelMessageWithSource:
	case DeferredChannelMessageWithSource:
	case DeferredUpdateMessage:
	case Modal:
		// Nothing
	default:
		panic("invalid response type")
	}

	return &InteractionResponse{
		Type: responseType,
		Data: data,
	}
}

type TextDisplay struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

func NewTextDisplay(content string) *TextDisplay {
	return &TextDisplay{
		Type:    10,
		Content: content,
	}
}

type ActionRow struct {
	Type       int           `json:"type"`
	Components []interface{} `json:"components"`
}

func NewActionRow(components []interface{}) *ActionRow {
	return &ActionRow{
		Type:       1,
		Components: components,
	}
}

type Button struct {
	Type     int    `json:"type"`
	Style    int    `json:"style"`
	Label    string `json:"label"`
	CustomID string `json:"custom_id"`
}

func NewButton(style int, label string, customID string) *Button {
	return &Button{
		Type:     2,
		Style:    style,
		Label:    label,
		CustomID: customID,
	}
}

type StringSelect struct {
	Type        int           `json:"type"`
	Options     []interface{} `json:"options"`
	Placeholder string        `json:"placeholder"`
	MinValues   int           `json:"min_values"`
	MaxValues   int           `json:"max_values"`
	CustomID    string        `json:"custom_id"`
}

func NewStringSelect(placeholder string, minValues, maxValues int, customID string, options []interface{}) *StringSelect {
	return &StringSelect{
		Type:        3,
		Placeholder: placeholder,
		MinValues:   minValues,
		MaxValues:   maxValues,
		CustomID:    customID,
		Options:     options,
	}
}

type StringOption struct {
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

type Container struct {
	Type        int           `json:"type"`
	AccentColor int           `json:"accent_color"`
	Components  []interface{} `json:"components"`
}

func NewContainer(accentColor int, components []interface{}) *Container {
	return &Container{
		Type:        17,
		AccentColor: accentColor,
		Components:  components,
	}
}

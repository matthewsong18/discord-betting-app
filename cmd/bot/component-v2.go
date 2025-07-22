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

type InteractionResponse struct {
	Type int         `json:"type"`
	Data MessageSend `json:"data"`
}

type TextDisplay struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ActionRow struct {
	Type       int           `json:"type"`
	Components []interface{} `json:"components"`
}

type Button struct {
	Type     int    `json:"type"`
	Style    int    `json:"style"`
	Label    string `json:"label"`
	CustomID string `json:"custom_id"`
}
type StringSelect struct {
	Type        int           `json:"type"`
	Options     []interface{} `json:"options"`
	Placeholder string        `json:"placeholder"`
	MinValues   int           `json:"min_values"`
	MaxValues   int           `json:"max_values"`
	CustomID    string        `json:"custom_id"`
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

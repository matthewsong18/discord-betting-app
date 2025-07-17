package main

/*

Discord has a new component v2 api that discordgo doesn't support

Need this flag to enable components v2 for the specific message

`IS_COMPONENTS_V2 = 1 << 15`

*/

var IsComponentsV2 = 1 << 15

type MessageSend struct {
	Flags      int           `json:"flags"`
	Components []interface{} `json:"components"`
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

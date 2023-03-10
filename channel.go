package aisera

type ChannelType int

const (
	ChannelType_InvalidChannel ChannelType = iota
	ChannelType_Facebook
	ChannelType_Slack
	ChannelType_Webchat
	ChannelType_Skype
	ChannelType_FormIntercept
	ChannelType_Email
	ChannelType_Teams
	ChannelType_SnowTIQ
	ChannelType_API
	ChannelType_CognitiveSearch
	ChannelType_SMS
	ChannelType_Widget
	ChannelType_ZoomChat
	ChannelType_WebexTeams
	ChannelType_Glip
	ChannelType_Twitter
	ChannelType_IVR
	ChannelType_RingCentralAssist
	ChannelType_GoogleChat
	ChannelType_Listening
	ChannelType_TicketAIAgentAssist
	ChannelType_SunshineConversation
)

type ChannelZoomSize struct {
	Width         any `json:"width,omitempty"`
	HeightPercent any `json:"heightPercent,omitempty"`
}

func (ChannelAccessParams) Default() ChannelAccessParams {
	return ChannelAccessParams{
		"authType":         "Email",
		"showHistory":      true,
		"showAiseraLogo":   true,
		"whitelistDomains": []string{},
	}
}

type ChannelAccessParams GenericKeyValue

type Channel struct {
	ID           int                 `json:"id,omitempty"`
	ChannelType  ChannelType         `json:"channel_type_id,omitempty"`
	AccessParams ChannelAccessParams `json:"access_params,omitempty"`
	Name         string              `json:"name,omitempty"`
	Bots         []Bot               `json:"bots,omitempty"`
}

type Channels []Channel

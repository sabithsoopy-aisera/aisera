package aisera

type BotConfig map[string]interface{}

type Bot struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Domain  string    `json:"domain"`
	BotType string    `json:"bot_type"`
	Config  BotConfig `json:"config"`
}

type Bots []Bot

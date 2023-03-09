package aisera

import (
	"fmt"
	"io"
	"strings"
)

type BotConfig map[string]interface{}

type Bot struct {
	ID                int       `json:"id,omitempty"`
	Name              string    `json:"name,omitempty"`
	Domain            string    `json:"domain,omitempty"`
	BotType           string    `json:"bot_type,omitempty"`
	Config            BotConfig `json:"config,omitempty"`
	Language          string    `json:"language,omitempty"`
	DomainDisplayName string    `json:"domain_display_name,omitempty"`
	Description       string    `json:"description,omitempty"`
}

func (b Bot) Validate() error {
	missingFields := []string{}
	if len(missingFields) == 0 {
		return nil
	}
	return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ","))
}

type Bots []Bot

func (b Bots) FilterBy(bot Bot) Bots {
	result := make(Bots, 0, len(b))
	for i := range b {
		if bot.ID != 0 && b[i].ID == bot.ID {
			return Bots{bot}
		}
		if bot.Name != "" && strings.Contains(strings.ToLower(b[i].Name), strings.ToLower(bot.Name)) {
			result = append(result, bot)
		}
	}
	return result
}

type DeleteEntityRequest struct {
	EntityID int `json:"entityId"`
}

func (d DeleteEntityRequest) JSONReader() io.Reader {
	return ToJSONReader(d)
}

type CreateEntityRequest struct {
	Entity interface{} `json:"entity"`
}

func (c CreateEntityRequest) JSONReader() io.Reader {
	return ToJSONReader(c)
}

package aisera

import (
	"fmt"
	"io"
	"strings"
)

type GenericKeyValue map[string]interface{}

type BotConfig GenericKeyValue

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
			return Bots{b[i]}
		}
		if bot.Name != "" && strings.Contains(strings.ToLower(b[i].Name), strings.ToLower(bot.Name)) {
			result = append(result, bot)
		}
	}
	return result
}

type JobParams struct {
	BotID int `json:"bot_id"`
}

type RetrainRequest struct {
	JobName   string    `json:"jobName"`
	JobParams JobParams `json:"jobParams"`
}

func (r RetrainRequest) JSONReader() io.Reader {
	return ToJSONReader(r)
}

type MappingRequest struct {
	AddedEntityIds []int  `json:"addedEntityIds,omitempty"`
	BotID          int    `json:"botId,omitempty"`
	EntityType     string `json:"entityType,omitempty"`
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

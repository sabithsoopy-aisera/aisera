package integration_test

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	aisera "go.aisera.cloud"
)

var _ = Describe("Bot lifecycle", func() {
	var (
		ctx          = context.Background()
		botName      = fmt.Sprintf("aisera integration_test: %d", time.Now().UnixNano())
		botCount     int
		newBotID     int
		newChannelID int
		channelName  = fmt.Sprintf("%s channe", botName)
	)
	It("can create , get and delete bots", func() {
		By("Getting all the bots that exists in the system", func() {
			bots, err := aiseraOffering.Bots(ctx, aisera.Filter{})
			Expect(err).NotTo(HaveOccurred())
			Expect(bots).NotTo(HaveLen(0))
			botCount = len(bots)
		})
		By("Creating the bot", func() {
			var err error
			log.Printf("creating the bot: %s", botName)
			newBotID, err = aiseraOffering.CreateBot(ctx, aisera.Bot{
				Name:              botName,
				DomainDisplayName: "KB Test",
				BotType:           "Conversation",
				Language:          "en",
				Config: aisera.BotConfig{
					"domainClassifierThreshold":          0,
					"endOfSessionSurvey":                 true,
					"enableKBTranslation":                false,
					"notificationChannels":               "LastActive",
					"conversationHistoryTimeStampFormat": "MM/dd/yyyy hh:mm:ss a z",
					"boostScore":                         0,
					"domain":                             "IT7.4",
					"hasUnhandledRequests":               false,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(newBotID).NotTo(Equal(0))
		})

		defer By("Deleting the bot that was created", func() {
			err := aiseraOffering.DeleteBot(ctx, newBotID)
			Expect(err).NotTo(HaveOccurred())
		})

		By("Getting the bots, we should see the bot created", func() {
			bots, err := aiseraOffering.Bots(ctx, aisera.Filter{
				Fields: []string{"name", "id", "baseDataSources.external_system_type_id", "channels.channel_type_id"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(bots).NotTo(HaveLen(botCount))
			createdBot := bots.FilterBy(aisera.Bot{
				ID: newBotID,
			})
			Expect(createdBot).To(HaveLen(1))
			Expect(createdBot[0].Name).To(Equal(botName))
		})
		By("Getting the bots, we should NOT see the bot created", func() {
			bots, err := aiseraOffering.Bots(ctx, aisera.Filter{})
			Expect(err).NotTo(HaveOccurred())
			createdBot := bots.FilterBy(aisera.Bot{
				ID: newBotID,
			})
			Expect(createdBot).To(HaveLen(0))
		})
		By("Creating channel", func() {
			var err error
			log.Printf("creating the channel: %s", channelName)
			newChannelID, err = aiseraOffering.CreateChannel(ctx, aisera.Channel{
				ChannelType:  aisera.ChannelType_Webchat,
				Name:         channelName,
				AccessParams: aisera.ChannelAccessParams{}.Default(),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(newChannelID).NotTo(Equal(0))
		})
		defer By("Deleting the channel", func() {
			err := aiseraOffering.DeleteChannel(ctx, newChannelID)
			Expect(err).NotTo(HaveOccurred())
		})
		By("Mapping the bot to channel", func() {
			err := aiseraOffering.MapBotToChannel(ctx, newBotID, []int{newChannelID})
			Expect(err).NotTo(HaveOccurred())
		})
		By("getting the channel, there should be bots info", func() {
			channels, err := aiseraOffering.Channels(ctx, aisera.Filter{
				Criteria: []aisera.Criteria{
					{
						Field:   []string{"bots.id"},
						Operand: fmt.Sprintf("%d", newBotID),
					},
				},
				Fields: []string{"bots.id", "name", "channel_type_id"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(HaveLen(1))
			Expect(channels[0].Bots).To(HaveLen(1))
		})
	})
})

package integration_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	aisera "go.aisera.cloud"
)

var _ = Describe("Bot lifecycle", func() {
	var (
		ctx      = context.Background()
		botName  = fmt.Sprintf("aisera integration_test: %d", time.Now().UnixNano())
		botCount int
		newBotID int
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
		By("Deleting the bot that was created", func() {
			err := aiseraOffering.DeleteBot(ctx, newBotID)
			Expect(err).NotTo(HaveOccurred())
		})
		By("Getting the bots, we should NOT see the bot created", func() {
			bots, err := aiseraOffering.Bots(ctx, aisera.Filter{})
			Expect(err).NotTo(HaveOccurred())
			createdBot := bots.FilterBy(aisera.Bot{
				ID: newBotID,
			})
			Expect(createdBot).To(HaveLen(0))
		})
	})
})

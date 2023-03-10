package integration_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	aisera "go.aisera.cloud"
)

var _ = Describe("Channel lifecycle", func() {
	var (
		ctx          = context.Background()
		channelName  = fmt.Sprintf("channel integration_test: %d", time.Now().UnixNano())
		newChannelID int
	)
	It("can do CRUD operations on channel", func() {
		By("Getting the channel", func() {
			channels, err := aiseraOffering.Channels(ctx, aisera.Filter{})
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).NotTo(HaveLen(0))
		})
		By("Creating channel", func() {
			var err error
			newChannelID, err = aiseraOffering.CreateChannel(ctx, aisera.Channel{
				ChannelType:  aisera.ChannelType_Webchat,
				Name:         channelName,
				AccessParams: aisera.ChannelAccessParams{}.Default(),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(newChannelID).NotTo(Equal(0))
		})
		By("Getting the channel that was just created", func() {
			channels, err := aiseraOffering.Channels(ctx, aisera.Filter{
				Criteria: []aisera.Criteria{
					{
						CanImply: true,
						Field:    []string{"name"},
						Operator: "LIKE",
						Operand:  channelName,
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(channels).To(HaveLen(1))
		})
		By("deleting the channel that was created", func() {
			err := aiseraOffering.DeleteChannel(ctx, newChannelID)
			Expect(err).NotTo(HaveOccurred())
		})
		By("deleting the channel again", func() {
			err := aiseraOffering.DeleteChannel(ctx, newChannelID)
			Expect(err).To(MatchError(`invalid status code: 500`))
		})
	})
})

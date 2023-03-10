package integration_test

import (
	"context"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	aisera "go.aisera.cloud"
)

var aiseraOffering aisera.Offerings

func TestAiseraOffering(t *testing.T) {
	if os.Getenv("AISERA_USERNAME") == "" || os.Getenv("AISERA_PASSWORD") == "" {
		t.Skip("skipping the test as required env variables are not set")
	}
	RegisterFailHandler(Fail)
	RunSpecs(t, "Aisera Offering")
}

var _ = SynchronizedBeforeSuite(func(ctx context.Context) []byte {
	var err error
	aiseraOffering, err = aisera.Login(ctx, aisera.LoginRequest{
		Username: os.Getenv("AISERA_USERNAME"),
		Password: os.Getenv("AISERA_PASSWORD"),
	})
	Expect(err).NotTo(HaveOccurred())
	return nil
}, func(data []byte) {
})

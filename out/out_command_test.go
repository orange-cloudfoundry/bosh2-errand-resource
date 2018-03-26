package out_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-community/bosh2-errand-resource/bosh/boshfakes"
	"github.com/cloudfoundry-community/bosh2-errand-resource/concourse"
	"github.com/cloudfoundry-community/bosh2-errand-resource/out"
)

var _ = Describe("OutCommand", func() {
	var (
		outCommand out.OutCommand
		director   *boshfakes.FakeDirector
	)

	BeforeEach(func() {
		director = new(boshfakes.FakeDirector)
		outCommand = out.NewOutCommand(director, nil, "")
	})

	Describe("Run", func() {
		var outRequest concourse.OutRequest

		BeforeEach(func() {
			outRequest = concourse.OutRequest{
				Source: concourse.Source{
					Target: "director.example.com",
				},
				Params: concourse.OutParams{
					ErrandName: "test_errand",
				},
			}
		})

		It("runs errand with defaults", func() {
			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			Expect(director.RunErrandCallCount()).To(Equal(1))
			params := director.RunErrandArgsForCall(0)
			Expect(params.ErrandName).To(Equal("test_errand"))
			Expect(params.KeepAlive).To(Equal(false))
			Expect(params.WhenChanged).To(Equal(false))
		})

		It("runs errand with KeepAlive", func() {
			outRequest.Params.KeepAlive = true
			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			Expect(director.RunErrandCallCount()).To(Equal(1))
			params := director.RunErrandArgsForCall(0)
			Expect(params.ErrandName).To(Equal("test_errand"))
			Expect(params.KeepAlive).To(Equal(true))
			Expect(params.WhenChanged).To(Equal(false))
		})

		It("runs errand with WhenChanged", func() {
			outRequest.Params.WhenChanged = true
			_, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			Expect(director.RunErrandCallCount()).To(Equal(1))
			params := director.RunErrandArgsForCall(0)
			Expect(params.ErrandName).To(Equal("test_errand"))
			Expect(params.KeepAlive).To(Equal(false))
			Expect(params.WhenChanged).To(Equal(true))
		})

		It("returns the new version", func() {
			sillyBytes := []byte{0xFE, 0xED, 0xDE, 0xAD, 0xBE, 0xEF}
			director.DownloadManifestReturns(sillyBytes, nil)

			outResponse, err := outCommand.Run(outRequest)
			Expect(err).ToNot(HaveOccurred())

			Expect(outResponse).To(Equal(out.OutResponse{
				Version: concourse.Version{
					ManifestSha1: "33bf00cb7a45258748f833a47230124fcc8fa3a4",
					Target:       "director.example.com",
				},
				Metadata: []concourse.Metadata{},
			}))
		})
	})
})

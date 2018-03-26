package concourse_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/cloudfoundry-community/bosh2-errand-resource/concourse"
)

var _ = Describe("NewOutRequest", func() {
	It("converts the config into an OutRequest", func() {
		config := []byte(`{
			"params": {
				"name": "smoke_tests",
        "keep_alive": true,
        "when_changed": true
			},
			"source": {
				"deployment": "mydeployment",
				"target": "director.example.com",
				"client": "foo",
				"client_secret": "foobar",
				"vars_store": {
					"provider": "gcs",
					"config": {
						"some": "dynamic",
						"keys": "per-provider"
					}
				}
			}
		}`)

		source, err := concourse.NewOutRequest(config, "")
		Expect(err).NotTo(HaveOccurred())

		Expect(source).To(Equal(concourse.OutRequest{
			Source: concourse.Source{
				Deployment:   "mydeployment",
				Target:       "director.example.com",
				Client:       "foo",
				ClientSecret: "foobar",
				VarsStore: concourse.VarsStore{
					Provider: "gcs",
					Config: map[string]interface{}{
						"some": "dynamic",
						"keys": "per-provider",
					},
				},
			},
			Params: concourse.OutParams{
				ErrandName:  "smoke_tests",
				KeepAlive:   true,
				WhenChanged: true,
			},
		}))
	})

	Context("when source_file param is passed", func() {
		It("overrides source with the values in the source_file", func() {
			sourceFile, _ := ioutil.TempFile("", "")
			sourceFile.WriteString(`{
				"deployment": "fileDeployment",
				"target": "fileDirector.com",
				"client_secret": "fileSecret",
				"vars_store": {
					"provider": "fileProvider",
					"config": {
						"file": "vars"
					}
				}
			}`)
			sourceFile.Close()

			configTemplate := `{
				"params": {
					"name": "smoke_tests",
					"source_file": "%s"
				},
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "original_client",
					"client_secret": "foobar",
					"vars_store": {
						"provider": "gcs",
						"config": {
							"some": "dynamic",
							"keys": "per-provider"
						}
					}
				}
			}`
			config := []byte(fmt.Sprintf(
				configTemplate,
				filepath.Base(sourceFile.Name()),
			))

			source, err := concourse.NewOutRequest(config, filepath.Dir(sourceFile.Name()))
			Expect(err).NotTo(HaveOccurred())

			Expect(source).To(Equal(concourse.OutRequest{
				Source: concourse.Source{
					Deployment:   "fileDeployment",
					Target:       "fileDirector.com",
					Client:       "original_client",
					ClientSecret: "fileSecret",
					VarsStore: concourse.VarsStore{
						Provider: "fileProvider",
						Config: map[string]interface{}{
							"file": "vars",
							"some": "dynamic",
							"keys": "per-provider",
						},
					},
				},
				Params: concourse.OutParams{
					ErrandName: "smoke_tests",
				},
			}))
		})
	})

	Context("when decoding fails", func() {
		It("errors", func() {
			config := []byte("not-json")

			_, err := concourse.NewOutRequest(config, "")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when a required parameter is missing", func() {
		It("returns an error with each missing parameter", func() {
			config := []byte(`{
				"source": {
					"deployment": "mydeployment",
					"target": "director.example.com",
					"client": "foo",
					"client_secret": "foobar"
				}
			}`)

			_, err := concourse.NewOutRequest(config, "")
			Expect(err).To(HaveOccurred())

			Expect(err.Error()).To(ContainSubstring("name"))
		})
	})
})

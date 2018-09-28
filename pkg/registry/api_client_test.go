package registry_test

import (
	. "github.com/onsi/gomega"
	"io/ioutil"
	"testing"
	"time"

	"net/http"

	"code.cloudfoundry.org/indicators/pkg/go_test"
	"code.cloudfoundry.org/indicators/pkg/registry"
)

func TestAPIClient_IndicatorDocuments(t *testing.T) {
	c := registry.NewAPIClient("http://localhost:8080", http.DefaultClient)

	t.Run("it fetches the payload on the /v1/indicator-documents endpoint", func(t *testing.T) {
		g := NewGomegaWithT(t)

		json, e := ioutil.ReadFile("../../pkg/registry/test_fixtures/example_response.json")
		g.Expect(e).ToNot(HaveOccurred())
		http.HandleFunc("/v1/indicator-documents", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write(json)
		})

		server := http.Server{
			Addr: "localhost:8080",
		}

		go server.ListenAndServe()
		defer server.Close()

		go_test.WaitForHTTPServer("localhost:8080", time.Second)

		documents, e := c.IndicatorDocuments()
		g.Expect(e).ToNot(HaveOccurred())

		g.Expect(documents[0].Labels["product"]).To(Equal("my-component"))
	})

	t.Run("it returns an error if the the client get fails", func(t *testing.T) {
		g := NewGomegaWithT(t)

		_, err := c.IndicatorDocuments()
		g.Expect(err).To(HaveOccurred())
	})
}
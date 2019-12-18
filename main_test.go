package main_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/EngineerBetter/cf-mysql-example-app"
)

var _ = Describe("The App", func() {
	It("Does a thing", func() {
		repo := NewInMemoryRepository()
		handler := NewMysqlHandler(repo)
		server := httptest.NewServer(handler)
		defer server.Close()
		client := http.DefaultClient

		request, err := http.NewRequest(http.MethodPut, server.URL+"/somevalue", strings.NewReader("testvalue"))
		Expect(err).ShouldNot(HaveOccurred())
		response, err := client.Do(request)
		Expect(err).ShouldNot(HaveOccurred())
		buff := new(bytes.Buffer)
		buff.ReadFrom(response.Body)
		responseBody := buff.String()
		Expect(responseBody).Should(Equal("created\n"))

		response, err = http.Get(server.URL + "/somevalue")
		Expect(err).ShouldNot(HaveOccurred())
		buff = new(bytes.Buffer)
		buff.ReadFrom(response.Body)
		responseBody = buff.String()
		Expect(responseBody).Should(Equal("testvalue\n"))

		request, err = http.NewRequest(http.MethodDelete, server.URL+"/somevalue", nil)
		Expect(err).ShouldNot(HaveOccurred())
		response, err = client.Do(request)
		Expect(err).ShouldNot(HaveOccurred())
		buff = new(bytes.Buffer)
		buff.ReadFrom(response.Body)
		responseBody = buff.String()
		Expect(responseBody).Should(Equal("deleted\n"))

		response, err = http.Get(server.URL + "/somevalue")
		Expect(err).ShouldNot(HaveOccurred())
		buff = new(bytes.Buffer)
		buff.ReadFrom(response.Body)
		responseBody = buff.String()
		Expect(responseBody).Should(Equal("key not found\n"))
	})
})

type InMemoryRepository struct {
	KeyValues map[string]string
}

func NewInMemoryRepository() InMemoryRepository {
	var r InMemoryRepository
	r.KeyValues = make(map[string]string)
	return r
}

func (r InMemoryRepository) Write(key, value string) error {
	r.KeyValues[key] = value
	return nil
}

func (r InMemoryRepository) Read(key string) (string, error) {
	return r.KeyValues[key], nil
}

func (r InMemoryRepository) Delete(key string) error {
	delete(r.KeyValues, key)
	return nil
}

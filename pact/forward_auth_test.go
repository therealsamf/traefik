package pact

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/middlewares/auth"
	"github.com/containous/traefik/v2/pkg/testhelpers"
	"github.com/go-check/check"
	"github.com/pact-foundation/pact-go/dsl"
	checker "github.com/vdemeester/shakers"
)

type ForwardAuthPactSuite struct{}

func (suite *ForwardAuthPactSuite) TestForwardAuth(c *check.C) {
	pact := &dsl.Pact{
		Consumer: "traefik.ForwardAuth",
		Provider: "xenia-auth-proxy",
		Host:     "localhost",
	}
	defer pact.Teardown()

	var forwardAuth http.Handler
	serverHandler := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		forwardAuth.ServeHTTP(responseWriter, request)
	})

	testServer := httptest.NewServer(serverHandler)
	defer testServer.Close()

	var test = func() (err error) {
		next := http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			responseWriter.WriteHeader(http.StatusOK)
		})

		authConfig := dynamic.ForwardAuth{
			Address: fmt.Sprintf("http://localhost:%d/", pact.Server.Port),
		}

		var _err error
		forwardAuth, _err = auth.NewForward(context.Background(), next, authConfig, "forwardAuthPactTest")

		if _err != nil {
			return err
		}

		request := testhelpers.MustNewRequest(http.MethodGet, testServer.URL, nil)
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}

		body, _err := ioutil.ReadAll(response.Body)
		if _err != nil {
			return _err
		}
		_err = response.Body.Close()
		if _err != nil {
			return _err
		}
		c.Log(string(body))

		c.Assert(response.StatusCode, check.Equals, http.StatusOK)

		return nil
	}

	_url, err := url.Parse(testServer.URL)
	c.Assert(err, checker.IsNil)

	pact.
		AddInteraction().
		Given("ForwardAuth middleware is configured correctly").
		UponReceiving("A request to the ForwardAuth middleware").
		WithRequest(dsl.Request{
			Method: "GET",
			Path:   dsl.String("/"),
			Headers: dsl.MapMatcher{
				"X-Forwarded-Host":  dsl.String(_url.Host),
				"X-Forwarded-Proto": dsl.String(_url.Scheme),
				"X-Forwarded-Uri":   dsl.String("/"),
			},
		}).
		WillRespondWith(dsl.Response{
			Status: 200,
		})

	if err := pact.Verify(test); err != nil {
		c.Logf("Error on Verify: %v", err)
	}
}

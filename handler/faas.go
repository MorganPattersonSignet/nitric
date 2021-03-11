package handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nitric-dev/membrane/sources"
)

// FaaSHandler - source handler for the membrane when operating in FaaS mode
type FaasHandler struct {
	host string
}

func errorToInternalServerError(err error) *http.Response {
	return &http.Response{
		Status:     "Internal Server Error",
		StatusCode: 500,
		// TODO: Eat error in non development modes
		// TODO: Log the error to an external log sink
		Body: ioutil.NopCloser(bytes.NewReader([]byte(err.Error()))),
	}
}

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *FaasHandler) HandleEvent(source *sources.Event) error {
	address := fmt.Sprint("http://%s", h.host)
	httpRequest, _ := http.NewRequest("POST", address, ioutil.NopCloser(bytes.NewReader(source.Payload)))
	httpRequest.Header.Add("x-nitric-request-id", source.ID)
	httpRequest.Header.Add("x-nitric-source-type", sources.SourceType_Subscription.String())
	httpRequest.Header.Add("x-nitric-source", source.Topic)

	// TODO: Handle response or error and response appropriately
	resp, err := http.DefaultClient.Do(httpRequest)

	if resp != nil && resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	} else if resp != nil {
		respMessage, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("Error processing event (%d): %s", resp.StatusCode, string(respMessage))
	}

	return fmt.Errorf("Error processing event: %s", err.Error())
}

// HandleHttpRequest - Handles an HTTP request by forwarding it as an HTTP request.
func (h *FaasHandler) HandleHttpRequest(source *sources.HttpRequest) *http.Response {
	address := fmt.Sprintf("http://%s", h.host)
	httpRequest, err := http.NewRequest("POST", address, source.Body)

	if err != nil {
		return errorToInternalServerError(err)
	}

	httpRequest.Header = source.Header
	httpRequest.Header.Add("x-nitric-source-type", sources.SourceType_Request.String())
	httpRequest.Header.Add("x-nitric-source", fmt.Sprintf("%s:%s", source.Method, source.Path))

	resp, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return errorToInternalServerError(err)
	}

	return resp
}

func NewFaasHandler(host string) *FaasHandler {
	return &FaasHandler{
		host: host,
	}
}

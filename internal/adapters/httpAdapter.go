package adapters

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

type Adapter struct {
	client *http.Client
}

func NewAdapter(client *http.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) Post(ctx context.Context, url string, headers map[string]string, payload []byte) ([]byte, error) {
	httpRequest, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		bytes.NewReader(payload),
	)

	if err != nil {
		return []byte(""), errCantBuildRequest
	}

	for key, value := range headers {
		httpRequest.Header.Set(key, value)
	}

	resp, err := a.client.Do(httpRequest)

	if err != nil {
		return []byte(""), errRequestCallFailed
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), errCantReadBody
	}

	defer resp.Body.Close()

	return body, nil
}

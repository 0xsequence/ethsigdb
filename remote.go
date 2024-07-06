package ethsigdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type RemoteLookup struct {
	// TODO: persist the http client here
	// TODO: use goware/cachestore to at least an in-memory cache
}

func NewRemoteLookup() (*RemoteLookup, error) {
	return &RemoteLookup{}, nil
}

func (f *RemoteLookup) FindEventSig(ctx context.Context, eventSig string) (*RemoteEventSigResponse, error) {
	request := fmt.Sprintf("https://api.openchain.xyz/signature-database/v1/lookup?event=%v&filter=false", eventSig)
	_, body, err := f.doRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	var out RemoteEventSigResponse
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (f *RemoteLookup) doRequest(ctx context.Context, request string) (int, []byte, error) {
	const maxRetries = 10
	retryCount := 0

	client := http.DefaultClient

retry:
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, request, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("unable to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if err = ctx.Err(); err != nil {
		return 0, nil, fmt.Errorf("aborted because context was done: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode > 499 && retryCount < maxRetries {
			select {
			case <-ctx.Done():
				// done
			default:
				retryCount++
				delay := time.Duration(retryCount) * time.Second * 2
				// fmt.Printf("fourbyte request to endpoint %s gave %s, trying again\n", request, resp.Status)
				time.Sleep(delay)
				goto retry
			}
		}

		return resp.StatusCode, nil, fmt.Errorf("fourbyte fail, status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("unable to read response body: %v: %w", request, err)
	}

	return resp.StatusCode, body, nil
}

type RemoteEventSigResponse struct {
	OK     bool `json:"ok"`
	Result struct {
		// Event map key is the event sig, ie.
		// 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
		Event map[string][]struct {
			// Name is the event name, ie. `Transfer(address,address,uint256)`
			Name string `json:"name"`

			// Filtered means its potentially junk event
			Filtered bool `json:"filtered"`
		} `json:"event"`
	} `json:"result"`
}

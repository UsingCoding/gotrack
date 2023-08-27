package youtrack

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"gotrack/pkg/app/tracker"
)

const (
	issueTimeTrackURL = "%s/api/issues/%s/timeTracking"
	addWorkItemUrl    = "%s/api/issues/%s/timeTracking/workItems"
)

func NewClient(host string, token string) tracker.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &client{
		host:  host,
		token: token,
		client: http.Client{
			Transport: tr,
		},
	}
}

type client struct {
	host  string
	token string

	client http.Client
}

func (c *client) IssueExists(ctx context.Context, issueID string) (bool, error) {
	url := fmt.Sprintf(issueTimeTrackURL, c.host, issueID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return false, errors.Wrap(err, "failed to create request")
	}
	defer req.Body.Close()

	c.authRequest(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return false, errors.Wrap(err, "failed to make addWorkItem request")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		data, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return false, errors.Wrap(err2, "failed to read response body")
		}

		return false, errors.Errorf("unknown status code %d returned, body %s", resp.StatusCode, string(data))
	}
}

type addWorkItemRequest struct {
	Type     string `json:"$type"`
	Text     string `json:"text"`
	Date     int64  `json:"date"`
	Duration struct {
		Presentation string `json:"presentation"`
		Type         string `json:"$type"`
	} `json:"duration"`
}

func (c *client) Track(ctx context.Context, params tracker.TrackParams) error {
	url := fmt.Sprintf(addWorkItemUrl, c.host, params.IssueID)

	d, err := marshalDuration(params.Duration)
	if err != nil {
		return errors.Wrap(err, "failed to marshal duration")
	}

	body := addWorkItemRequest{
		Type: "IssueWorkItem",
		Text: params.Text,
		Date: params.Date.UnixMilli(),
		Duration: struct {
			Presentation string `json:"presentation"`
			Type         string `json:"$type"`
		}{
			Presentation: d,
			Type:         "DurationValue",
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return errors.Wrap(err, "failed to marshal addWorkItemRequest")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	c.authRequest(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make addWorkItem request")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed read response data")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
		return errors.WithStack(tracker.ErrIssueNotFound)
	default:
		return errors.Errorf("unknown status code %d returned, body %s", resp.StatusCode, string(respData))
	}
}

func (c *client) authRequest(r *http.Request) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	r.Header.Add("Content-Type", "application/json")
}

func marshalDuration(duration time.Duration) (string, error) {
	if duration < time.Minute {
		return "", errors.New("cannot track less than minute")
	}

	minutes := int(duration / time.Minute)

	return fmt.Sprintf("%dm", minutes), nil
}

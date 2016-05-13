package gosparkpost

import (
	"encoding/json"
	"fmt"
	URL "net/url"
)

// https://developers.sparkpost.com/api/#/reference/suppression-list
var suppressionListsPathFormat = "/api/v%d/suppression-list"

type SuppressionEntry struct {
	// Email is used when list is stored
	Email string `json:"email,omitempty"`

	// Recipient is used when a list is returned
	Recipient string `json:"recipient,omitempty"`

	Transactional    bool   `json:"transactional,omitempty"`
	NonTransactional bool   `json:"non_transactional,omitempty"`
	Source           string `json:"source,omitempty"`
	Description      string `json:"description,omitempty"`
	Updated          string `json:"updated,omitempty"`
	Created          string `json:"created,omitempty"`
}

type SuppressionListWrapper struct {
	Results    []*SuppressionEntry `json:"results,omitempty"`
	Recipients []SuppressionEntry  `json:"recipients,omitempty"`
}

func (c *Client) SuppressionList() (*SuppressionListWrapper, error) {
	path := fmt.Sprintf(suppressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)

	return doSuppressionRequest(c, finalUrl)
}

func (c *Client) SuppressionRetrieve(recipientEmail string) (*SuppressionListWrapper, error) {
	path := fmt.Sprintf(suppressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, recipientEmail)

	return doSuppressionRequest(c, finalUrl)
}

func (c *Client) SuppressionSearch(parameters map[string]string) (*SuppressionListWrapper, error) {
	var finalUrl string
	path := fmt.Sprintf(suppressionListsPathFormat, c.Config.ApiVersion)

	if parameters == nil || len(parameters) == 0 {
		finalUrl = fmt.Sprintf("%s%s", c.Config.BaseUrl, path)
	} else {
		params := URL.Values{}
		for k, v := range parameters {
			params.Add(k, v)
		}

		finalUrl = fmt.Sprintf("%s%s?%s", c.Config.BaseUrl, path, params.Encode())
	}

	return doSuppressionRequest(c, finalUrl)
}

func (c *Client) SuppressionDelete(recipientEmail string) (res *Response, err error) {
	path := fmt.Sprintf(suppressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s/%s", c.Config.BaseUrl, path, recipientEmail)

	res, err = c.HttpDelete(finalUrl)
	if err != nil {
		return nil, err
	}

	if res.HTTP.StatusCode >= 200 && res.HTTP.StatusCode <= 299 {
		return

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("SuppressionEntry", "delete")
		if err != nil {
			return nil, err
		}

		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return
}

func (c *Client) SuppressionInsertOrUpdate(entries []SuppressionEntry) (err error) {
	if entries == nil {
		err = fmt.Errorf("send `entries` cannot be nil here")
		return
	}

	path := fmt.Sprintf(suppressionListsPathFormat, c.Config.ApiVersion)
	finalUrl := fmt.Sprintf("%s%s", c.Config.BaseUrl, path)

	list := SuppressionListWrapper{nil, entries}

	return c.send(finalUrl, list)

}

func (c *Client) send(finalUrl string, recipients SuppressionListWrapper) (err error) {
	jsonBytes, err := json.Marshal(recipients)
	if err != nil {
		return
	}

	res, err := c.HttpPut(finalUrl, jsonBytes)
	if err != nil {
		return
	}

	if err = res.AssertJson(); err != nil {
		return
	}

	err = res.ParseResponse()
	if err != nil {
		return
	}

	if res.HTTP.StatusCode == 200 {

	} else if len(res.Errors) > 0 {
		// handle common errors
		err = res.PrettyError("Transmission", "create")
		if err != nil {
			return
		}

		err = fmt.Errorf("%d: %s", res.HTTP.StatusCode, string(res.Body))
	}

	return
}

func doSuppressionRequest(c *Client, finalUrl string) (*SuppressionListWrapper, error) {
	// Send off our request
	res, err := c.HttpGet(finalUrl)
	if err != nil {
		return nil, err
	}

	// Assert that we got a JSON Content-Type back
	if err = res.AssertJson(); err != nil {
		return nil, err
	}

	// Get the Content
	bodyBytes, err := res.ReadBody()
	if err != nil {
		return nil, err
	}

	// Parse expected response structure
	var resMap SuppressionListWrapper
	err = json.Unmarshal(bodyBytes, &resMap)

	if err != nil {
		return nil, err
	}

	return &resMap, err
}

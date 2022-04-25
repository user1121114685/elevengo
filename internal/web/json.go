package web

import (
	"bytes"
	"encoding/json"
	"github.com/deadblue/gostream/quietly"
	"io"
)

// CallJsonApi calls remote HTTP API, and parses its result as JSON.
func (c *Client) CallJsonApi(url string, qs Params, form Params, resp interface{}) (err error) {
	// Prepare request
	var body io.ReadCloser
	if form != nil {
		body, err = c.PostForm(url, qs, form)
	} else {
		body, err = c.Get(url, qs)
	}
	if err != nil {
		return
	}
	defer quietly.Close(body)
	// Parse response
	if resp != nil {
		decoder := json.NewDecoder(body)
		err = decoder.Decode(resp)
	}
	return
}

func (c *Client) CallJsonpApi(url string, qs Params, resp interface{}) (err error) {
	body, err := c.Get(url, qs)
	if err != nil {
		return
	}
	defer quietly.Close(body)
	data, err := io.ReadAll(body)
	if err != nil {
		return
	}
	left, right := bytes.IndexByte(data, '('), bytes.LastIndexByte(data, ')')
	if left < 0 || right < 0 {
		err = &json.SyntaxError{Offset: 0}
	} else {
		err = json.Unmarshal(data[left+1:right], resp)
	}
	return
}
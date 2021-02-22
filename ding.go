package notifiers

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Ding impl notifier, notify msg by Dingding webhook.
type Ding struct {
	webhook    string
	secret     string
	httpclient *http.Client
	NoopCloser
}

// TextMsg is dingding text message type.
type TextMsg struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

// Resp is dingding api response.
type Resp struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// NewDing create a instance.
func NewDing(webhook string, secret string) *Ding {
	return &Ding{
		webhook:    webhook,
		secret:     secret,
		httpclient: http.DefaultClient,
	}
}

// GetName impl Notifier.GetName.
func (d *Ding) GetName() string {
	return "dingding"
}

// Notify impl Notifier.Notify.
func (d *Ding) Notify(ctx context.Context, to string, msg Message) error {
	textMsg := &TextMsg{
		Msgtype: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: msg.Content},
	}

	data, err := json.Marshal(textMsg)

	if err != nil {
		return errors.Wrap(err, "json marshal")
	}

	finalUrl := d.webhook + d.getSignQuery(time.Now())
	req, err := http.NewRequest(http.MethodPost, finalUrl, bytes.NewReader(data))
	if err != nil {
		return errors.Wrap(err, "create request")
	}
	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(ctx)

	r, err := d.httpclient.Do(req)
	if err != nil {
		return errors.Wrap(err, "request")
	}
	defer r.Body.Close()

	var resp Resp
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		return errors.Wrap(err, "decode resp")
	}

	if resp.Errcode != 0 {
		return errors.Wrap(ErrNotify, resp.Errmsg)
	}

	return nil
}

func (d *Ding) getSignQuery(now time.Time) string {
	if d.secret == "" {
		return ""
	}
	ts := now.Unix() * 1000
	sign := d.getSign(ts)
	return fmt.Sprintf("&timestamp=%d&sign=%s", ts, sign)
}

func (d *Ding) getSign(ts int64) string {
	signStr := fmt.Sprintf("%d\n%s", ts, d.secret)
	h := hmac.New(sha256.New, []byte(d.secret))
	_, _ = h.Write([]byte(signStr))
	res := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(res)
}

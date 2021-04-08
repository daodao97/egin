package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

const maxContentLength = 1024 * 1024

type Mac struct {
	AccessKey string `json:"access_key"`
	SecretKey []byte `json:"secret_key"`
	Debug     bool   `json:"debug"`
}

func (m *Mac) SignRequest(req *http.Request) (token string, err error) {
	h := hmac.New(sha1.New, m.SecretKey)
	hb := bytes.NewBuffer(nil)

	u := req.URL
	_, _ = io.WriteString(h, req.Method+" "+u.Path)
	if m.Debug {
		_, _ = io.WriteString(hb, req.Method+" "+u.Path)
	}
	if u.RawQuery != "" {
		_, _ = io.WriteString(h, "?"+u.RawQuery)
		if m.Debug {
			_, _ = io.WriteString(hb, "?"+u.RawQuery)
		}
	}
	_, _ = io.WriteString(h, "\nHost: "+req.Host)
	if m.Debug {
		_, _ = io.WriteString(hb, "\nHost: "+req.Host)
	}

	ctType := req.Header.Get("Content-Type")
	if ctType != "" {
		_, _ = io.WriteString(h, "\nContent-Type: "+ctType)
		if m.Debug {
			_, _ = io.WriteString(hb, "\nContent-Type: "+ctType)
		}
	}

	_, _ = io.WriteString(h, "\n\n")
	if m.Debug {
		_, _ = io.WriteString(hb, "\n\n")
	}

	if incBody(req, ctType) {
		b, er := ioutil.ReadAll(req.Body)
		if er != nil {
			return "", er
		}
		_, _ = h.Write(b)
		if m.Debug {
			_, _ = hb.Write(b)
		}
		req.Body = &readCloser{bytes.NewReader(b), req.Body}
	}
	if m.Debug {
		glog.Info("SignRequest:", hb.String())
	}

	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return m.AccessKey + ":" + sign, nil

}

func incBody(req *http.Request, ctType string) bool {
	typeOk := ctType != "" && ctType != "application/octet-stream"
	lengthOk := req.ContentLength > 0 && req.ContentLength < maxContentLength
	return typeOk && lengthOk && req.Body != nil
}

type readCloser struct {
	io.Reader
	io.Closer
}

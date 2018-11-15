package httpbody

import (
	"io/ioutil"
	"net/url"
	"strings"
	"testing"
)

func TestNewFormBody(t *testing.T) {
	contentType, bodyReader := NewFormBody2(
		url.Values{
			"v1": []string{"a"},
			"v2": []string{"b"},
		},
		Files{
			"f1": []File{
				NewFile("/Users/henrylee2cn/f11.txt", strings.NewReader("f11 text.")),
			},
		},
	)
	b, _ := ioutil.ReadAll(bodyReader)
	t.Logf("\nContent-Type:\n%s\nBody:\n%s", contentType, b)
}

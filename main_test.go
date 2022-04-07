package main

import (
	"bytes"
	"html/template"
	"testing"
)

func TestTemplate_GetURLs(t *testing.T) {
	var (
		buf  bytes.Buffer
		data = TemplateArgs{
			HttpVersion:       "http",
			ImageFingerprints: generateImageFingerprints(10),
			Script:            template.JS(loadTimeScript),
		}
	)
	err := template.Must(template.New("").Parse(indexHTML)).Execute(&buf, data)
	if err != nil {
		t.Fail()
	}
	urls := getURLs(&buf, 1000)
	if len(urls) != 10 {
		t.Fail()
	}
}

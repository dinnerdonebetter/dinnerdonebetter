package main

import (
	"bytes"
	"embed"
	"net/http"

	"github.com/yuin/goldmark"
	g "maragu.dev/gomponents"
)

//go:embed content/*.md
var contentFS embed.FS

var md = goldmark.New()

func (s *ConsumerFrontendServer) TermsPage(_ http.ResponseWriter, _ *http.Request) (g.Node, error) {
	html, err := renderMarkdown("content/terms.md")
	if err != nil {
		return nil, err
	}
	return legalPage("Terms of Service", g.Raw(html)), nil
}

func (s *ConsumerFrontendServer) PrivacyPage(_ http.ResponseWriter, _ *http.Request) (g.Node, error) {
	html, err := renderMarkdown("content/privacy.md")
	if err != nil {
		return nil, err
	}
	return legalPage("Privacy Policy", g.Raw(html)), nil
}

func renderMarkdown(path string) (string, error) {
	raw, err := contentFS.ReadFile(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err = md.Convert(raw, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

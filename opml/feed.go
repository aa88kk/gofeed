package opml

import (
	"encoding/json"
	//"time"
)

type Feed struct {
	Title    string
	Version  string
	Outlines []*Outline
}

type Outline struct {
	Title       string
	XmlUrl      string
	HtmlUrl     string
	Type        string
	Text        string
	Description string
	Outlines    []*Outline
}

func (f Feed) String() string {
	json, _ := json.MarshalIndent(f, "", "    ")
	return string(json)
}

package opml

import (
	//"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/mmcdole/gofeed/internal/shared"
	"github.com/mmcdole/goxpp"
)

// Parser is an Atom Parser
type Parser struct{}

// Parse parses an xml feed into an atom.Feed
func (ap *Parser) Parse(feed io.Reader) (*Feed, error) {
	p := xpp.NewXMLPullParser(feed, false, shared.NewReaderLabel)

	_, err := shared.FindRoot(p)
	if err != nil {
		return nil, err
	}

	return ap.parseRoot(p)
}

func (ap *Parser) parseOutline(p *xpp.XMLPullParser) ([]*Outline, error) {
	fmt.Println("parseOutline()")

	outlines := []*Outline{}
	depth := p.Depth
	for {
		tok, err := p.NextTag()
		if err != nil {
			fmt.Println("error2:", err, p.EventName(p.Event), p.Name, depth)
			return nil, err
		}

		if tok == xpp.EndTag {
			//fmt.Println("tag end", p.EventName(tok), p.Name, p.Depth)
			return outlines, nil
			if p.Depth <= depth {
			}
		}

		if err := p.Expect(xpp.StartTag, "outline"); err != nil {
			fmt.Println("expect error", err, p.Name)
			return nil, err
		}

		ol := &Outline{}
		ol.Title = p.Attribute("title")
		ol.Text = p.Attribute("text")
		ol.Description = p.Attribute("description")
		ol.Type = p.Attribute("type")
		ol.XmlUrl = p.Attribute("xmlUrl")
		ol.HtmlUrl = p.Attribute("htmlUrl")
		fmt.Println(strings.Repeat("    ", p.Depth), p.Name, ol.Text, p.Depth)

		outlines = append(outlines, ol)

		tok, _ = p.Next()
		//fmt.Println("next :", p.EventName(tok), p.Name, p.Depth)

		if tok == xpp.Text {
			//embed
			//fmt.Println("embed outline next:", p.Name, p.Depth, p.EventName(tok), p.Text)
			ol.Outlines, _ = ap.parseOutline(p)
			//fmt.Println("enbed end.")
		}
	}
	return outlines, nil
}

func (ap *Parser) parseRoot(p *xpp.XMLPullParser) (*Feed, error) {
	if err := p.Expect(xpp.StartTag, "opml"); err != nil {
		return nil, err
	}

	opml := &Feed{}
	opml.Outlines = []*Outline{}
	opml.Version = p.Attribute("version")

	for {
		tok, err := shared.NextTag(p)
		if p.Event == xpp.EndDocument {
			fmt.Println("enddoc")
			break
		}

		if err != nil {
			fmt.Println("err:", p.EventName(tok))
			return nil, err
		}

		if tok == xpp.EndTag {
			//fmt.Println("end:", p.Name)
			break
		}

		if tok == xpp.StartTag {
			name := strings.ToLower(p.Name)
			switch name {
			case "title":
				err := p.DecodeElement(&opml.Title)
				if err != nil {
					return nil, err
				}
				fmt.Println("title:", opml.Title)
			case "body":
				opml.Outlines, _ = ap.parseOutline(p)
			default:
				err := p.Skip()
				if err != nil {
					return nil, err
				}

			}
		}
	}
	return opml, nil
}

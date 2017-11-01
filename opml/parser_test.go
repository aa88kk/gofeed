package opml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

func TestParserParse(t *testing.T) {
	files, _ := filepath.Glob("../testdata/parser/opml/*.opml")

	for _, f := range files {
		base := filepath.Base(f)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("Testing %s... \n", name)

		// Get actual source feed
		ff := fmt.Sprintf("../testdata/parser/opml/%s.opml", name)
		f, err := ioutil.ReadFile(ff)
		if err != nil {
			t.Fatal(err)
		}

		fp := &Parser{}
		opml, err := fp.Parse(bytes.NewReader(f))
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(opml.Outlines)

		// Get json encoded expected feed result
		ef := fmt.Sprintf("../testdata/parser/opml/%s.json", name)
		e, err := ioutil.ReadFile(ef)
		if err != nil {
			t.Fatal(err)
		}

		// Unmarshal expected feed
		expected := &Feed{}
		json.Unmarshal(e, &expected)

		if assert.Equal(t, expected, opml, "Feed file %s.xml did not match expected output %s.json", name, name) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

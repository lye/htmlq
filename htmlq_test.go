package htmlq

import (
	"testing"
)

func Test1(t *testing.T) {
	var hq HtmlQ
	hq.ParseString(`<!DOCTYPE html><html><body><span class="post-id">102</span><div class="msg">Hello</div><span class="post-id">103</span><div class="msg">There</div></body></html>`)

	found := []string{}

	hq.Find("span.post-id").ForEach(func(hq HtmlQ) {
		found = append(found, hq.Text())
	})

	if len(found) != 2 {
		t.Fatalf("Not enough things found: %#v\n", found)
	}

	if found[0] != "102" {
		t.Fatalf("First found item didn't match expected: %#v\n", found)
	}

	if found[1] != "103" {
		t.Fatalf("Second found item didn't match expected: %#v\n", found)
	}
}

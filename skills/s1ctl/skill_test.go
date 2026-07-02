package skill

import (
	"strings"
	"testing"
)

func TestMarkdownEmbedded(t *testing.T) {
	md := Markdown()
	if md == "" {
		t.Fatal("embedded SKILL.md is empty")
	}
	if !strings.Contains(md, "s1ctl") {
		t.Fatal("embedded SKILL.md does not mention s1ctl")
	}
}

func TestParse(t *testing.T) {
	doc := Parse()
	if doc.Name != "s1ctl" {
		t.Fatalf("expected name s1ctl, got %q", doc.Name)
	}
	if doc.Description == "" {
		t.Fatal("description is empty")
	}
	if doc.Body == "" {
		t.Fatal("body is empty")
	}
	if !strings.Contains(doc.Body, "## Session bootstrap") {
		t.Fatal("body missing expected section")
	}
}

func TestSplitFrontmatter(t *testing.T) {
	fm, body := splitFrontmatter("---\nname: test\n---\n# Body")
	if fm != "name: test" {
		t.Fatalf("unexpected frontmatter: %q", fm)
	}
	if body != "# Body" {
		t.Fatalf("unexpected body: %q", body)
	}

	fm2, body2 := splitFrontmatter("no frontmatter")
	if fm2 != "" || body2 != "no frontmatter" {
		t.Fatalf("unexpected no-fm result: %q / %q", fm2, body2)
	}
}

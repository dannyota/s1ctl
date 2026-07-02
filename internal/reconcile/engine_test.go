package reconcile

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// fakeDecode parses a trivial "name: X" line so LoadDir/WriteDir can round-trip
// without pulling in a YAML dependency.
func fakeDecode(data []byte) (Object, error) {
	for line := range strings.SplitSeq(string(data), "\n") {
		if v, ok := strings.CutPrefix(line, "name: "); ok {
			return Object{Name: v, Body: data}, nil
		}
	}
	return Object{}, errors.New("no name field")
}

func TestApplyCreatesThenUpdatesInOrder(t *testing.T) {
	var order []string
	s := Surface{
		Name:    "device rule",
		Command: "devicecontrol",
		Create: func(_ context.Context, local Object) error {
			order = append(order, "create:"+local.Name)
			return nil
		},
		Update: func(_ context.Context, id string, local Object) error {
			order = append(order, "update:"+id+":"+local.Name)
			return nil
		},
	}

	c1 := obj("c1", "", "b")
	c2 := obj("c2", "", "b")
	u1 := obj("u1", "", "b")
	u2 := obj("u2", "", "b")
	p := Plan{Items: []Item{
		{Kind: "update", Name: "u2", ID: "idu2", Local: &u2},
		{Kind: "create", Name: "c2", Local: &c2},
		{Kind: "create", Name: "c1", Local: &c1},
		{Kind: "update", Name: "u1", ID: "idu1", Local: &u1},
	}}

	sum, err := Apply(context.Background(), s, p, io_Discard())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != (Summary{Created: 2, Updated: 2}) {
		t.Errorf("summary = %+v, want {2 2 0}", sum)
	}
	want := []string{"create:c1", "create:c2", "update:idu1:u1", "update:idu2:u2"}
	if !reflect.DeepEqual(order, want) {
		t.Errorf("call order = %v, want %v", order, want)
	}
}

func TestApplyPartialFailure(t *testing.T) {
	s := Surface{
		Name: "device rule",
		Create: func(_ context.Context, local Object) error {
			if local.Name == "bad" {
				return errors.New("boom")
			}
			return nil
		},
		Update: func(_ context.Context, _ string, _ Object) error { return nil },
	}
	good := obj("good", "", "b")
	bad := obj("bad", "", "b")
	up := obj("up", "", "b")
	p := Plan{Items: []Item{
		{Kind: "create", Name: "good", Local: &good},
		{Kind: "create", Name: "bad", Local: &bad},
		{Kind: "update", Name: "up", ID: "idup", Local: &up},
	}}

	var warn bytes.Buffer
	sum, err := Apply(context.Background(), s, p, &warn)
	if err == nil {
		t.Fatal("expected error on partial failure")
	}
	if !strings.Contains(err.Error(), "failed") {
		t.Errorf("error %q should mention 'failed'", err.Error())
	}
	if sum != (Summary{Created: 1, Updated: 1, Failed: 1}) {
		t.Errorf("summary = %+v, want {1 1 1}", sum)
	}
	if !strings.Contains(warn.String(), "bad") {
		t.Errorf("warn output %q should name the failing item", warn.String())
	}
}

func TestApplyNilCreateIsFailure(t *testing.T) {
	s := Surface{
		Name: "cloud policy",
		Caps: Capabilities{NoCreate: true},
		// Create is nil.
		Update: func(_ context.Context, _ string, _ Object) error { return nil },
	}
	c := obj("newpolicy", "", "b")
	p := Plan{Items: []Item{{Kind: "create", Name: "newpolicy", Local: &c}}}

	var warn bytes.Buffer
	sum, err := Apply(context.Background(), s, p, &warn)
	if err == nil {
		t.Fatal("expected error when creating on a NoCreate surface")
	}
	if sum.Failed != 1 || sum.Created != 0 {
		t.Errorf("summary = %+v, want Failed=1 Created=0", sum)
	}
	if !strings.Contains(warn.String(), "cannot create") {
		t.Errorf("warn output %q should contain 'cannot create'", warn.String())
	}
}

func TestApplyCleanPlanNoError(t *testing.T) {
	s := Surface{
		Name:   "device rule",
		Create: func(_ context.Context, _ Object) error { return nil },
		Update: func(_ context.Context, _ string, _ Object) error { return nil },
	}
	p := Plan{Items: []Item{{Kind: "unchanged"}, {Kind: "live-only"}}}
	sum, err := Apply(context.Background(), s, p, io_Discard())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sum != (Summary{}) {
		t.Errorf("summary = %+v, want zero", sum)
	}
}

func TestLoadDirMissingReturnsNil(t *testing.T) {
	objs, err := LoadDir(filepath.Join(t.TempDir(), "does-not-exist"), fakeDecode)
	if err != nil {
		t.Fatalf("missing dir should not error: %v", err)
	}
	if objs != nil {
		t.Errorf("missing dir should return nil objects, got %v", objs)
	}
}

func TestLoadDirDecodeErrorNamesFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.yaml"), []byte("nope: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadDir(dir, fakeDecode)
	if err == nil {
		t.Fatal("expected decode error")
	}
	if !strings.Contains(err.Error(), "bad.yaml") {
		t.Errorf("error %q should name the offending file", err.Error())
	}
}

func TestWriteLoadRoundTrip(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "out")
	objs := []Object{
		{Name: "foo", Body: []byte("name: foo\nx: 1\n")},
		{Name: "bar", Body: []byte("name: bar\ny: 2\n")},
	}
	stale, err := WriteDir(dir, objs)
	if err != nil {
		t.Fatal(err)
	}
	if len(stale) != 0 {
		t.Errorf("fresh dir should have no stale files, got %v", stale)
	}

	loaded, err := LoadDir(dir, fakeDecode)
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]Object{}
	for _, o := range loaded {
		byName[o.Name] = o
	}
	if len(byName) != 2 {
		t.Fatalf("round trip: want 2 objects, got %d", len(byName))
	}
	for _, want := range objs {
		got, ok := byName[want.Name]
		if !ok {
			t.Errorf("missing object %q after round trip", want.Name)
			continue
		}
		if !bytes.Equal(got.Body, want.Body) {
			t.Errorf("body mismatch for %q: %q != %q", want.Name, got.Body, want.Body)
		}
	}
}

func TestWriteDirDuplicateStemSuffix(t *testing.T) {
	dir := t.TempDir()
	objs := []Object{
		{Name: "Same Name", Body: []byte("name: Same Name\n")},
		{Name: "Same Name", Body: []byte("name: Same Name\n")},
		{Name: "Same Name", Body: []byte("name: Same Name\n")},
	}
	if _, err := WriteDir(dir, objs); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"same-name.yaml", "same-name-1.yaml", "same-name-2.yaml"} {
		if _, err := os.Stat(filepath.Join(dir, want)); err != nil {
			t.Errorf("expected file %q: %v", want, err)
		}
	}
}

func TestWriteDirReportsStale(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "extra.yaml"), []byte("name: extra\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "keeper.yml"), []byte("name: keeper\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stale, err := WriteDir(dir, []Object{{Name: "foo", Body: []byte("name: foo\n")}})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(stale)
	want := []string{"extra.yaml", "keeper.yml"}
	if !reflect.DeepEqual(stale, want) {
		t.Errorf("stale = %v, want %v", stale, want)
	}
}

func TestSanitizeName(t *testing.T) {
	cases := map[string]string{
		"Simple Name":  "simple-name",
		"  Trimmed  ":  "trimmed",
		"Weird/Chars!": "weirdchars",
		"":             "object",
		"   ":          "object",
	}
	for in, want := range cases {
		if got := sanitizeName(in); got != want {
			t.Errorf("sanitizeName(%q) = %q, want %q", in, got, want)
		}
	}
}

// io_Discard returns a throwaway writer for tests that ignore warnings.
func io_Discard() *bytes.Buffer { return &bytes.Buffer{} }

package reconcile

import (
	"reflect"
	"strings"
	"testing"
)

func obj(name, id, body string) Object {
	return Object{Name: name, ID: id, Body: []byte(body)}
}

func TestBuildPlanClassification(t *testing.T) {
	cases := []struct {
		name      string
		local     []Object
		live      []Object
		create    int
		update    int
		unchanged int
		liveOnly  int
		warnings  int
		wantErr   bool
	}{
		{
			name:   "empty inputs",
			local:  nil,
			live:   nil,
			create: 0,
		},
		{
			name:   "local only is create",
			local:  []Object{obj("a", "", "body-a")},
			live:   nil,
			create: 1,
		},
		{
			name:     "live only is reported",
			local:    nil,
			live:     []Object{obj("a", "id-a", "body-a")},
			liveOnly: 1,
		},
		{
			name:      "byte equal is unchanged",
			local:     []Object{obj("a", "", "body-a")},
			live:      []Object{obj("a", "id-a", "body-a")},
			unchanged: 1,
		},
		{
			name:   "differing body is update",
			local:  []Object{obj("a", "", "body-new")},
			live:   []Object{obj("a", "id-a", "body-old")},
			update: 1,
		},
		{
			name: "mixed",
			local: []Object{
				obj("keep", "", "same"),
				obj("change", "", "new"),
				obj("add", "", "fresh"),
			},
			live: []Object{
				obj("keep", "id-keep", "same"),
				obj("change", "id-change", "old"),
				obj("gone", "id-gone", "orphan"),
			},
			create:    1,
			update:    1,
			unchanged: 1,
			liveOnly:  1,
		},
		{
			name:      "duplicate live first wins",
			local:     []Object{obj("a", "", "body-a")},
			live:      []Object{obj("a", "id-first", "body-a"), obj("a", "id-second", "body-x")},
			unchanged: 1,
			warnings:  1,
		},
		{
			name:    "duplicate local is error",
			local:   []Object{obj("a", "", "b1"), obj("a", "", "b2")},
			live:    nil,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p, warnings, err := BuildPlan(tc.local, tc.live)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if len(p.Items) != 0 {
					t.Fatalf("expected empty plan on error, got %d items", len(p.Items))
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := len(p.Creates()); got != tc.create {
				t.Errorf("creates = %d, want %d", got, tc.create)
			}
			if got := len(p.Updates()); got != tc.update {
				t.Errorf("updates = %d, want %d", got, tc.update)
			}
			if got := len(p.Unchanged()); got != tc.unchanged {
				t.Errorf("unchanged = %d, want %d", got, tc.unchanged)
			}
			if got := len(p.LiveOnly()); got != tc.liveOnly {
				t.Errorf("live-only = %d, want %d", got, tc.liveOnly)
			}
			if got := len(warnings); got != tc.warnings {
				t.Errorf("warnings = %d, want %d", got, tc.warnings)
			}
		})
	}
}

func TestBuildPlanUpdateCarriesLiveID(t *testing.T) {
	p, _, err := BuildPlan(
		[]Object{obj("a", "", "new")},
		[]Object{obj("a", "id-a", "old")},
	)
	if err != nil {
		t.Fatal(err)
	}
	ups := p.Updates()
	if len(ups) != 1 {
		t.Fatalf("want 1 update, got %d", len(ups))
	}
	if ups[0].ID != "id-a" {
		t.Errorf("update ID = %q, want id-a", ups[0].ID)
	}
	if ups[0].Local == nil || ups[0].Live == nil {
		t.Errorf("update item must carry both Local and Live")
	}
}

func TestBuildPlanCreateHasNoLiveID(t *testing.T) {
	p, _, err := BuildPlan([]Object{obj("a", "", "body")}, nil)
	if err != nil {
		t.Fatal(err)
	}
	cs := p.Creates()
	if len(cs) != 1 {
		t.Fatalf("want 1 create, got %d", len(cs))
	}
	if cs[0].ID != "" {
		t.Errorf("create ID = %q, want empty", cs[0].ID)
	}
	if cs[0].Live != nil {
		t.Errorf("create item must have nil Live")
	}
}

func TestBuildPlanDuplicateLiveWarning(t *testing.T) {
	_, warnings, err := BuildPlan(
		[]Object{obj("a", "", "body-a")},
		[]Object{obj("a", "id-first", "body-a"), obj("a", "id-second", "body-x")},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 1 {
		t.Fatalf("want 1 warning, got %d", len(warnings))
	}
	if !strings.Contains(warnings[0], "duplicate") || !strings.Contains(warnings[0], "a") {
		t.Errorf("warning %q missing 'duplicate' or the name", warnings[0])
	}
}

func TestBuildPlanDuplicateLocalError(t *testing.T) {
	_, _, err := BuildPlan([]Object{obj("dup", "", "b1"), obj("dup", "", "b2")}, nil)
	if err == nil {
		t.Fatal("want error for duplicate local names")
	}
	if !strings.Contains(err.Error(), "dup") {
		t.Errorf("error %q should name the duplicate", err.Error())
	}
}

func TestBuildPlanInputOrderIndependent(t *testing.T) {
	local := []Object{obj("c", "", "3"), obj("a", "", "1"), obj("b", "", "2")}
	live := []Object{obj("b", "idb", "2"), obj("z", "idz", "9"), obj("a", "ida", "0")}

	p1, w1, err1 := BuildPlan(local, live)
	if err1 != nil {
		t.Fatal(err1)
	}

	rl := []Object{local[2], local[0], local[1]}
	rv := []Object{live[2], live[0], live[1]}
	p2, w2, err2 := BuildPlan(rl, rv)
	if err2 != nil {
		t.Fatal(err2)
	}

	if !reflect.DeepEqual(p1, p2) {
		t.Errorf("plan not order independent:\n%+v\n%+v", p1, p2)
	}
	if !reflect.DeepEqual(w1, w2) {
		t.Errorf("warnings not order independent:\n%+v\n%+v", w1, w2)
	}
}

func TestPlanEmpty(t *testing.T) {
	unchangedOnly := Plan{Items: []Item{{Kind: "unchanged"}, {Kind: "live-only"}}}
	if !unchangedOnly.Empty() {
		t.Errorf("plan with only unchanged/live-only should be Empty")
	}
	withCreate := Plan{Items: []Item{{Kind: "create"}}}
	if withCreate.Empty() {
		t.Errorf("plan with a create is not Empty")
	}
	withUpdate := Plan{Items: []Item{{Kind: "update"}}}
	if withUpdate.Empty() {
		t.Errorf("plan with an update is not Empty")
	}
}

func TestDescribe(t *testing.T) {
	cases := []struct {
		name string
		plan Plan
		want string
	}{
		{
			name: "create and update with unchanged",
			plan: Plan{Items: []Item{
				{Kind: "create"},
				{Kind: "create"},
				{Kind: "update"},
				{Kind: "unchanged"},
				{Kind: "unchanged"},
				{Kind: "unchanged"},
			}},
			want: "create 2 device rules, update 1 device rule (3 unchanged) from dir",
		},
		{
			name: "create only no unchanged",
			plan: Plan{Items: []Item{{Kind: "create"}}},
			want: "create 1 device rule from dir",
		},
		{
			name: "update only",
			plan: Plan{Items: []Item{{Kind: "update"}, {Kind: "update"}}},
			want: "update 2 device rules from dir",
		},
		{
			name: "nothing to change",
			plan: Plan{Items: []Item{{Kind: "unchanged"}, {Kind: "live-only"}}},
			want: "no changes for device rules in dir",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Describe(tc.plan, "device rule", "dir"); got != tc.want {
				t.Errorf("Describe = %q, want %q", got, tc.want)
			}
		})
	}
}

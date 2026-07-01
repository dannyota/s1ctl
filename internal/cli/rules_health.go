package cli

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newRulesHealthCmd() *cobra.Command {
	var siteIDs []string

	cmd := &cobra.Command{
		Use:   "health",
		Short: "Classify rules by operational state",
		Long: `Fetch all custom detection rules and classify them as firing (active
with alerts), silent (active with zero alerts), disabled, or erroring
(reached alert limit). Helps identify rules that need attention.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}

			params := &mgmt.RuleListParams{SiteIDs: siteIDs, Limit: 1000}
			rules, _, err := fetchAllREST("rule", func(cur string) ([]mgmt.Rule, *mgmt.Pagination, error) {
				params.Cursor = cur
				return c.RulesList(cmd.Context(), params)
			})
			if err != nil {
				return err
			}

			type classified struct {
				Rule  mgmt.Rule
				State string
			}

			var items []classified
			var firing, silent, disabled, erroring, other int

			for _, r := range rules {
				var state string
				switch {
				case r.ReachedLimit:
					state = "erroring"
					erroring++
				case r.Status == mgmt.RuleStatusDisabled:
					state = "disabled"
					disabled++
				case r.Status == mgmt.RuleStatusActive && r.GeneratedAlerts > 0:
					state = "firing"
					firing++
				case r.Status == mgmt.RuleStatusActive && r.GeneratedAlerts == 0:
					state = "silent"
					silent++
				default:
					state = string(r.Status)
					other++
				}
				items = append(items, classified{Rule: r, State: state})
			}

			sort.Slice(items, func(i, j int) bool {
				order := map[string]int{"erroring": 0, "silent": 1, "firing": 2, "disabled": 3}
				oi, oki := order[items[i].State]
				oj, okj := order[items[j].State]
				if !oki {
					oi = 4
				}
				if !okj {
					oj = 4
				}
				if oi != oj {
					return oi < oj
				}
				return items[i].Rule.GeneratedAlerts > items[j].Rule.GeneratedAlerts
			})

			if outputFormat == "json" {
				type jsonItem struct {
					ID     string `json:"id"`
					Name   string `json:"name"`
					State  string `json:"state"`
					Alerts int    `json:"alerts"`
					Status string `json:"status"`
					OS     string `json:"os"`
					Scope  string `json:"scope"`
				}
				out := make([]jsonItem, len(items))
				for i, it := range items {
					out[i] = jsonItem{
						ID:     it.Rule.ID,
						Name:   it.Rule.Name,
						State:  it.State,
						Alerts: it.Rule.GeneratedAlerts,
						Status: string(it.Rule.Status),
						OS:     parseOSTarget(it.Rule.S1QL),
						Scope:  string(it.Rule.Scope),
					}
				}
				return printJSON(cmd.OutOrStdout(), out)
			}

			headers := []string{"Name", "State", "Alerts", "Severity", "OS", "Scope", "Response"}
			rows := make([][]string, len(items))
			for i, it := range items {
				response := "-"
				if it.Rule.TreatAsThreat != "" && it.Rule.TreatAsThreat != mgmt.RuleTreatUndefined {
					response = string(it.Rule.TreatAsThreat)
				}
				rows[i] = []string{
					truncate(it.Rule.Name, 40),
					it.State,
					fmt.Sprintf("%d", it.Rule.GeneratedAlerts),
					string(it.Rule.Severity),
					parseOSTarget(it.Rule.S1QL),
					string(it.Rule.Scope),
					response,
				}
			}
			printTable(headers, rows)

			fmt.Fprintf(cmd.OutOrStdout(), "\n%s: %d firing, %d silent, %d disabled, %d erroring",
				pluralize(len(rules), "rule"), firing, silent, disabled, erroring)
			if other > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), ", %d other", other)
			}
			fmt.Fprintln(cmd.OutOrStdout())
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&siteIDs, "site-id", nil, "filter by site ID")
	return cmd
}

var (
	osEqualRe = regexp.MustCompile(`(?i)endpoint\.os\s*=\s*'([^']+)'`)
	osInRe    = regexp.MustCompile(`(?i)endpoint\.os\s+in\s*\(([^)]+)\)`)
	osValRe   = regexp.MustCompile(`'([^']+)'`)
)

func parseOSTarget(s1ql string) string {
	seen := map[string]bool{}
	var oses []string
	addOS := func(os string) {
		os = strings.ToLower(os)
		if !seen[os] {
			seen[os] = true
			oses = append(oses, os)
		}
	}
	for _, m := range osEqualRe.FindAllStringSubmatch(s1ql, -1) {
		addOS(m[1])
	}
	for _, m := range osInRe.FindAllStringSubmatch(s1ql, -1) {
		for _, v := range osValRe.FindAllStringSubmatch(m[1], -1) {
			addOS(v[1])
		}
	}
	if len(oses) == 0 {
		return "any"
	}
	return strings.Join(oses, ",")
}

package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newAuditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit",
		Short: "View mutation audit log",
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newAuditListCmd())
	return cmd
}

func newAuditListCmd() *cobra.Command {
	var last int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recent audit log entries",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			path := auditLogPath()
			if path == "" {
				return fmt.Errorf("cannot determine home directory")
			}
			f, err := os.Open(path)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Fprintln(cmd.OutOrStdout(), "No audit entries yet.")
					return nil
				}
				return err
			}
			defer f.Close() //nolint:errcheck

			var records []auditRecord
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				var rec auditRecord
				if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
					continue
				}
				records = append(records, rec)
			}
			if err := scanner.Err(); err != nil {
				return err
			}

			if len(records) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No audit entries yet.")
				return nil
			}

			start := 0
			if last > 0 && last < len(records) {
				start = len(records) - last
			}
			visible := records[start:]

			headers := []string{"Timestamp", "Command", "Action", "Target", "Result"}
			rows := make([][]string, len(visible))
			for i, r := range visible {
				rows[i] = []string{r.Timestamp, r.Command, r.Action, r.Target, r.Result}
			}
			return printOutput(cmd.OutOrStdout(), headers, rows, visible, len(visible), len(visible), "audit entry", true)
		},
	}
	cmd.Flags().IntVar(&last, "last", 25, "number of recent entries to show")
	return cmd
}

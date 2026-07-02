package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

func newMaintenanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage task maintenance-window configuration",
		Long: `Manage maintenance-window and concurrency configuration for background tasks
(the tasks-configuration API).

Each configuration is scoped and keyed by a task type (--task-type), e.g.
agents_upgrade. 'get'/'set' use the classic per-day maintenance-window format;
'get-flexible'/'set-flexible' use the flexible policy_payload format, which
requires the flexible maintenance-window SKU. 'export' writes all window
occurrences for a scope as CSV (flexible format only).`,
	}
	requireSubcommand(cmd)
	cmd.AddCommand(newMaintenanceGetCmd(false))
	cmd.AddCommand(newMaintenanceGetCmd(true))
	cmd.AddCommand(newMaintenanceSetCmd())
	cmd.AddCommand(newMaintenanceSetFlexibleCmd())
	cmd.AddCommand(newMaintenanceExportCmd())
	return cmd
}

// maintenanceScopeFlags binds the shared scope + task-type flags and builds the
// params used by every maintenance command.
type maintenanceScopeFlags struct {
	taskType   string
	siteIDs    []string
	accountIDs []string
	groupIDs   []string
	tenant     bool
}

func (s *maintenanceScopeFlags) register(cmd *cobra.Command) {
	cmd.Flags().StringVar(&s.taskType, "task-type", "", "task type, e.g. agents_upgrade (required)")
	cmd.Flags().StringSliceVar(&s.siteIDs, "site-id", nil, "scope to site IDs")
	cmd.Flags().StringSliceVar(&s.accountIDs, "account-id", nil, "scope to account IDs")
	cmd.Flags().StringSliceVar(&s.groupIDs, "group-id", nil, "scope to group IDs")
	cmd.Flags().BoolVar(&s.tenant, "tenant", false, "scope to the global (tenant) level")
}

func (s *maintenanceScopeFlags) params() (*mgmt.TasksConfigParams, error) {
	if s.taskType == "" {
		return nil, fmt.Errorf("--task-type is required")
	}
	return &mgmt.TasksConfigParams{
		TaskType:   mgmt.TaskType(s.taskType),
		SiteIDs:    s.siteIDs,
		AccountIDs: s.accountIDs,
		GroupIDs:   s.groupIDs,
		Tenant:     s.tenant,
	}, nil
}

func newMaintenanceGetCmd(flexible bool) *cobra.Command {
	use, short := "get", "Get the maintenance-window configuration for a scope"
	if flexible {
		use, short = "get-flexible", "Get the flexible maintenance-window configuration for a scope"
	}
	var scope maintenanceScopeFlags

	cmd := &cobra.Command{
		Use:   use,
		Short: short,
		RunE: func(cmd *cobra.Command, _ []string) error {
			params, err := scope.params()
			if err != nil {
				return err
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			var cfg *mgmt.TasksConfig
			if flexible {
				cfg, err = c.TasksConfigFlexibleGet(cmd.Context(), params)
			} else {
				cfg, err = c.TasksConfigGet(cmd.Context(), params)
			}
			if err != nil {
				return err
			}
			return printJSON(cmd.OutOrStdout(), cfg)
		},
	}
	scope.register(cmd)
	return cmd
}

func newMaintenanceSetCmd() *cobra.Command {
	var scope maintenanceScopeFlags
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "set --task-type <type> --from-file <data.json>",
		Short: "Set the maintenance-window configuration for a scope",
		Long: `Set the classic per-day maintenance-window configuration for a scope.

--from-file supplies the configuration data payload (maxConcurrent, timezoneGmt,
maintenanceWindowsByDay, inherit flags); the scope and task type come from the
--task-type and scope flags.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			params, err := scope.params()
			if err != nil {
				return err
			}
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			raw, err := os.ReadFile(fromFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", fromFile, err)
			}
			var data mgmt.TasksConfigData
			if err := json.Unmarshal(raw, &data); err != nil {
				return fmt.Errorf("parse %s: %w", fromFile, err)
			}
			body := mgmt.TasksConfigWrite{
				Data: data,
				Filter: mgmt.TasksConfigFilter{
					TaskType:   params.TaskType,
					SiteIDs:    params.SiteIDs,
					AccountIDs: params.AccountIDs,
					GroupIDs:   params.GroupIDs,
					Tenant:     params.Tenant,
				},
			}

			action := fmt.Sprintf("set %s maintenance config from %s", params.TaskType, fromFile)
			return guard(cmd.OutOrStdout(), "maintenance set", action, fromFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if _, err := c.TasksConfigUpdate(cmd.Context(), body); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "taskType": string(params.TaskType)})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Updated %s maintenance configuration\n", params.TaskType)
				return nil
			})
		},
	}
	scope.register(cmd)
	cmd.Flags().StringVar(&fromFile, "from-file", "", "configuration data JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newMaintenanceSetFlexibleCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "set-flexible --from-file <body.json>",
		Short: "Set the flexible maintenance-window configuration",
		Long: `Set the flexible (policy_payload) maintenance-window configuration.

The flexible format is SKU-gated and open-ended, so --from-file must contain the
full request body: a "data" object with the policy payload and a "filter" object
with the task type and scope. The body is sent verbatim.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			raw, err := os.ReadFile(fromFile)
			if err != nil {
				return fmt.Errorf("read %s: %w", fromFile, err)
			}
			if !json.Valid(raw) {
				return fmt.Errorf("parse %s: not valid JSON", fromFile)
			}

			action := fmt.Sprintf("set flexible maintenance config from %s", fromFile)
			return guard(cmd.OutOrStdout(), "maintenance set-flexible", action, fromFile, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if _, err := c.TasksConfigFlexibleUpdate(cmd.Context(), json.RawMessage(raw)); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "updated", "format": "flexible"})
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Updated flexible maintenance configuration")
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "full request body JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newMaintenanceExportCmd() *cobra.Command {
	var scope maintenanceScopeFlags
	var outFile string

	cmd := &cobra.Command{
		Use:   "export --task-type <type> --out <file>",
		Short: "Export maintenance-window occurrences as CSV",
		Long: `Export all maintenance-window occurrences for a scope as CSV. Only the flexible
(policy_payload) maintenance-window format is supported.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			params, err := scope.params()
			if err != nil {
				return err
			}
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			data, err := c.MaintenanceWindowsExport(cmd.Context(), params)
			if err != nil {
				return err
			}
			if outFile == "-" {
				_, err = cmd.OutOrStdout().Write(data)
				return err
			}
			if err := os.WriteFile(outFile, data, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Exported maintenance windows to %s\n", outFile)
			return nil
		},
	}
	scope.register(cmd)
	cmd.Flags().StringVar(&outFile, "out", "maintenance-windows.csv", "output file (use - for stdout)")
	return cmd
}

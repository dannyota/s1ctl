package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCloudOnboardingOnboardCmd() *cobra.Command {
	var fromFile string
	var yes bool

	cmd := &cobra.Command{
		Use:   "onboard --from-file <request.json>",
		Short: "Onboard a cloud entity from a JSON file",
		Long: `Onboard a new cloud entity (AWS account, GCP project, Azure
subscription, etc.) using a CnappCloudOnboardingRequest JSON file. The file
must contain the full onboarding payload including cloudProvider, products,
and credentials.

Dry-run by default; pass --yes to apply.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			return guard(cmd.OutOrStdout(), "cloud-onboarding onboard", "onboard cloud entity from "+fromFile, fromFile, yes, func() error {
				input, err := readRuleJSONFile(fromFile)
				if err != nil {
					return err
				}
				c, err := gqlClient()
				if err != nil {
					return err
				}
				resp, err := c.CnappOnboard(cmd.Context(), input, nil)
				if err != nil {
					return err
				}
				msg := ""
				success := false
				if resp != nil {
					msg = resp.Message
					success = resp.IsSuccess
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{
						"isSuccess": success,
						"message":   msg,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Onboard: success=%v message=%s\n", success, orDash(msg))
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "path to onboarding request JSON file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply (default: dry-run)")
	return markJSON(cmd)
}

func newCloudOnboardingDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete <account-id> [account-id...]",
		Short: "Delete (offboard) cloud entities",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return guard(cmd.OutOrStdout(), "cloud-onboarding delete", "delete "+pluralize(len(args), "cloud entity"), joinIDs(args), yes, func() error {
				c, err := gqlClient()
				if err != nil {
					return err
				}
				resp, err := c.CnappDelete(cmd.Context(), args, nil)
				if err != nil {
					return err
				}
				msg := ""
				success := false
				if resp != nil {
					msg = resp.Message
					success = resp.IsSuccess
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]any{
						"isSuccess":  success,
						"message":    msg,
						"accountIds": args,
					})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Delete: success=%v message=%s\n", success, orDash(msg))
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply (default: dry-run)")
	return markJSON(cmd)
}

func joinIDs(ids []string) string {
	if len(ids) <= 3 {
		return fmt.Sprintf("%v", ids)
	}
	return fmt.Sprintf("%v and %d more", ids[:3], len(ids)-3)
}

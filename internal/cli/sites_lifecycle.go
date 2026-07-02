package cli

import (
	"fmt"
	"io"
	"strconv"

	"github.com/spf13/cobra"

	"danny.vn/s1/mgmt"
)

// noteSensitiveOutput prints a one-line stderr reminder that a command emitted
// secret material to stdout. It never contains the secret itself.
func noteSensitiveOutput(w io.Writer) {
	fmt.Fprintln(w, "Note: output contains sensitive material — handle accordingly.")
}

// validateReactivateChoice enforces that exactly one of --unlimited or
// --expiration is set on a reactivate command. Requiring an explicit choice
// prevents silently reactivating as perpetual/never-expire.
func validateReactivateChoice(unlimited bool, expiration string) error {
	if unlimited == (expiration != "") {
		return fmt.Errorf("specify exactly one of --unlimited or --expiration")
	}
	return nil
}

func newSitesReactivateCmd() *cobra.Command {
	var (
		unlimited  bool
		expiration string
		yes        bool
	)
	cmd := &cobra.Command{
		Use:   "reactivate <site-id>",
		Short: "Reactivate an expired site",
		Long: `Reactivate an expired site. Specify exactly one of --unlimited (no
expiration) or --expiration (an RFC3339 timestamp) to set the new license
window.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if err := validateReactivateChoice(unlimited, expiration); err != nil {
				return err
			}
			return guard(cmd.OutOrStdout(), "sites reactivate", "reactivate site "+id, id, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.SitesReactivate(cmd.Context(), id, unlimited, expiration); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "reactivated", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Reactivated site %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&unlimited, "unlimited", false, "reactivate with no expiration")
	cmd.Flags().StringVar(&expiration, "expiration", "", "new expiration as an RFC3339 timestamp")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newSitesExpireCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "expire <site-id>",
		Short: "Expire a site immediately",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "sites expire", "expire site "+id, id, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				if err := c.SitesExpireNow(cmd.Context(), id); err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), map[string]string{"status": "expired", "id": id})
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Expired site %s\n", id)
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newSitesDuplicateCmd() *cobra.Command {
	var (
		name, sourceSiteID, policySource string
		totalLicenses                    int
		copyUsers, unlimitedLicenses     bool
		yes                              bool
	)
	cmd := &cobra.Command{
		Use:   "duplicate",
		Short: "Duplicate an existing site",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if name == "" {
				return fmt.Errorf("--name is required")
			}
			if sourceSiteID == "" {
				return fmt.Errorf("--source-site-id is required")
			}
			src, err := strconv.ParseInt(sourceSiteID, 10, 64)
			if err != nil {
				return fmt.Errorf("--source-site-id must be numeric: %w", err)
			}
			ps := mgmt.SitePolicySource(policySource)
			switch ps {
			case mgmt.PolicySourceInheritGlobal, mgmt.PolicySourceCopySourceSite, mgmt.PolicySourceCustom:
			default:
				return fmt.Errorf("--policy-source must be one of inherit_global, copy_source_site, custom")
			}
			data := mgmt.SiteDuplicate{
				Name:              name,
				SourceSiteID:      src,
				PolicySource:      ps,
				CopyUsers:         copyUsers,
				UnlimitedLicenses: unlimitedLicenses,
			}
			if cmd.Flags().Changed("total-licenses") {
				data.TotalLicenses = &totalLicenses
			}
			action := fmt.Sprintf("duplicate site from %s as %q", sourceSiteID, name)
			return guard(cmd.OutOrStdout(), "sites duplicate", action, name, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				s, err := c.SitesDuplicate(cmd.Context(), data)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), s)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Created site %s (%s)\n", s.Name, s.ID)
				return nil
			})
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "new site name (required)")
	cmd.Flags().StringVar(&sourceSiteID, "source-site-id", "", "source site ID to copy from (required)")
	cmd.Flags().StringVar(&policySource, "policy-source", string(mgmt.PolicySourceInheritGlobal), "policy origin: inherit_global, copy_source_site, custom")
	cmd.Flags().BoolVar(&copyUsers, "copy-users", false, "copy users from the source site")
	cmd.Flags().IntVar(&totalLicenses, "total-licenses", 0, "total licenses for the new site")
	cmd.Flags().BoolVar(&unlimitedLicenses, "unlimited-licenses", false, "unlimited licenses")
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newSitesRegenerateKeyCmd() *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "regenerate-key <site-id>",
		Short: "Regenerate a site registration key",
		Long: `Regenerate a site's registration key. On apply, the new registration
token is printed to stdout — treat it as a secret.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			return guard(cmd.OutOrStdout(), "sites regenerate-key", "regenerate registration key for site "+id, id, yes, func() error {
				c, err := mgmtClient()
				if err != nil {
					return err
				}
				tok, err := c.SitesRegenerateKey(cmd.Context(), id)
				if err != nil {
					return err
				}
				if outputFormat == "json" {
					return printJSON(cmd.OutOrStdout(), tok)
				}
				fmt.Fprintln(cmd.OutOrStdout(), tok.Value())
				noteSensitiveOutput(cmd.ErrOrStderr())
				return nil
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "apply the action (default: dry-run)")
	return cmd
}

func newSitesTokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token <site-id>",
		Short: "Print a site's registration token",
		Long: `Print a site's current registration token to stdout. The token is
sensitive registration material.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := mgmtClient()
			if err != nil {
				return err
			}
			tok, err := c.SitesToken(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			if outputFormat == "json" {
				if err := printJSON(cmd.OutOrStdout(), tok); err != nil {
					return err
				}
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), tok.Value())
			}
			noteSensitiveOutput(cmd.ErrOrStderr())
			return nil
		},
	}
	return cmd
}

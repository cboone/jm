package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var readCmd = &cobra.Command{
	Use:   "read <email-id>",
	Short: "Read the full content of an email",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return exitError("authentication_failed", err.Error(),
				"Check your token in JMAP_TOKEN or config file")
		}

		emailID := args[0]
		preferHTML, _ := cmd.Flags().GetBool("html")
		rawHeaders, _ := cmd.Flags().GetBool("raw-headers")
		showThread, _ := cmd.Flags().GetBool("thread")

		if showThread {
			tv, err := c.ReadThread(emailID, preferHTML, rawHeaders)
			if err != nil {
				return exitError("not_found", err.Error(), "")
			}
			return formatter().Format(os.Stdout, tv)
		}

		detail, err := c.ReadEmail(emailID, preferHTML, rawHeaders)
		if err != nil {
			return exitError("not_found", err.Error(), "")
		}

		return formatter().Format(os.Stdout, detail)
	},
}

func init() {
	readCmd.Flags().Bool("html", false, "prefer HTML body (default: plain text)")
	readCmd.Flags().Bool("raw-headers", false, "include all raw headers")
	readCmd.Flags().Bool("thread", false, "show all emails in the same thread")
	rootCmd.AddCommand(readCmd)
}

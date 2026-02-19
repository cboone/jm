package cmd

import "github.com/spf13/cobra"

var sieveCmd = &cobra.Command{
	Use:   "sieve",
	Short: "Manage sieve filtering scripts",
	Long: `Manage sieve filtering scripts on the server.

Sieve scripts control server-side email filtering. Only one script can be
active at a time. Use 'sieve list' to see all scripts, 'sieve create' to
add a new script, and 'sieve activate' to enable it.`,
}

func init() {
	rootCmd.AddCommand(sieveCmd)
}

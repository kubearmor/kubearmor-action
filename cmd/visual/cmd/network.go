package cmd

import (
	"github.com/spf13/cobra"
)

var networkCmd = &cobra.Command{
	Use:     "network",
	Short:   "network subcommand is a command to visualization network connection behaviors.",
	Example: "visual network -f [json file name] -o [png file name]",
	Run: func(cmd *cobra.Command, args []string) {

		// fmt.Fprint(os.Stdout, output)
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)
}

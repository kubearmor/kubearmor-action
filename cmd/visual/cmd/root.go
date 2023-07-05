package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "visual",
	Short:   "visual is a command to visualization system or network behaviors.",
	Example: "visual system -f [json file name] -o [png file name]\nvisual network -f [json file name] -o [png file name]",
}

func Execute() {
	rootCmd.Execute()
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/t-shah02/pacview/internal/ui"
	"github.com/t-shah02/pacview/internal/utils"
)

var rootCmd = &cobra.Command{
	Use:   "pacview",
	Short: "A visual representation of your pacman dependencies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		packages := utils.GetInstalledPacmanPackages()
		if packages == nil {
			fmt.Fprintln(os.Stderr, "pacview: could not list packages (is pacman installed?)")
			os.Exit(1)
		}
		if err := ui.Run(packages); err != nil {
			fmt.Fprintln(os.Stderr, "pacview:", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

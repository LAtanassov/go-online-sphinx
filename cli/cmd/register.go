package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers to Online SPHINX",
	Long:  `Registers to Online SPHINX`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("register called")
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)
}

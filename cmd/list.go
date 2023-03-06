package cmd

import (
	"fmt"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/spf13/cobra"
)

var (
	listConfig = &cobra.Command{
		Use:   "list",
		Short: "list algs",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			format := "- %s\n"
			fmt.Println("Kems:")
			for _, kem := range crypto.GetKemNames() {
				fmt.Printf(format, kem)
			}
			fmt.Println("Signatures:")
			for _, sign := range crypto.GetSignNames() {
				fmt.Printf(format, sign)
			}
		},
	}
)

func init() {
	configCmd.AddCommand(listConfig)
}

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
			fmt.Println("Kems:")
			for _, kem := range crypto.GetAllKems() {
				fmt.Printf("- %s\n", kem)
			}
			fmt.Println("Signatures:")
			for _, sign := range crypto.GetAllSigns() {
				fmt.Printf("- %s\n", sign)
			}
		},
	}
)

func init() {
	configCmd.AddCommand(listConfig)
}

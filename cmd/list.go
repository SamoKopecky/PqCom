package cmd

import (
	"fmt"

	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/spf13/cobra"
)

var (
	listConfig = &cobra.Command{
		Use:   "list",
		Short: "List available algorithms",
		Long:  `List available post-quantum algorithms that can be specified in the configuration file`,
		Run: func(cmd *cobra.Command, args []string) {
			format := "- %s\n"
			fmt.Println("Key encapsulation methods (kem_alg):")
			for _, kem := range crypto.GetKemNames() {
				fmt.Printf(format, kem)
			}
			fmt.Println("Digital signatures (sign_alg):")
			for _, sign := range crypto.GetSignNames() {
				fmt.Printf(format, sign)
			}
		},
	}
)

func init() {
	configCmd.AddCommand(listConfig)
}

package cmd

import (
	"github.com/SamoKopecky/pqcom/main/crypto"
	"github.com/spf13/cobra"
)

var (
	keygen = &cobra.Command{
		Use:   "keygen",
		Short: "Generate key pair",
		Long:  `Generate configured public/private key pair`,
		Run: func(cmd *cobra.Command, args []string) {
			crypto.WriteKeys()
		},
	}
)

func init() {
	rootCmd.AddCommand(keygen)
}

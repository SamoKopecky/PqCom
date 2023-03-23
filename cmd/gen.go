package cmd

import (
	"fmt"
	"os"

	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/crypto"
	kemAlg "github.com/SamoKopecky/pqcom/main/crypto/kem"
	signAlg "github.com/SamoKopecky/pqcom/main/crypto/sign"
	"github.com/spf13/cobra"
)

var (
	kem  string
	sign string

	genConfig = &cobra.Command{
		Use:   "gen",
		Short: "Generate configuration file",
		Long:  "Generate a configuration and also generate the public/private key pair for the given algorithms",
		Run: func(cmd *cobra.Command, args []string) {
			if !crypto.IsValidAlg(kem, crypto.GetKemNames) {
				invalidPrint("kem algorithm", cmd)
				os.Exit(1)
			}
			if !crypto.IsValidAlg(sign, crypto.GetSignNames) {
				invalidPrint("signature algorithm", cmd)
				os.Exit(1)
			}
			config.GenerateConfig(kem, sign)
		},
	}
)

func init() {
	configCmd.AddCommand(genConfig)

	genConfig.Flags().StringVarP(&kem, "kem", "k", crypto.GetKemById(kemAlg.PqComKyber1024{}.Id()), "specify a Key Encapsulation Method")
	genConfig.Flags().StringVarP(&sign, "sign", "s", crypto.GetSignById(signAlg.PqComDilithium5{}.Id()), "specify a digital signature algorithm")
}

func invalidPrint(what string, cmd *cobra.Command) {
	fmt.Printf("Error: invalid %s, check 'pqcom config list'\n\n", what)
	cmd.Help()
	fmt.Println()
}

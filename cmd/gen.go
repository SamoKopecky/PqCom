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
		Short: "generate config",
		Long:  `TODO`,
		Run: func(cmd *cobra.Command, args []string) {
			if !crypto.IsValidAlg(kem, crypto.GetAllKems) {
				invalidPrint("kem algorithm", cmd)
				os.Exit(1)
			}
			if !crypto.IsValidAlg(sign, crypto.GetAllSigns) {
				invalidPrint("signature algorithm", cmd)
				os.Exit(1)
			}
			config.GenerateConfig(kem, sign)
		},
	}
)

func init() {
	configCmd.AddCommand(genConfig)

	genConfig.Flags().StringVarP(&kem, "kem", "k", crypto.GetKemById(kemAlg.PqComKyber512{}.Id()), "Use kem")
	genConfig.Flags().StringVarP(&sign, "sign", "s", crypto.GetSignById(signAlg.PqComDilithium2{}.Id()), "Use sign")
}

func invalidPrint(what string, cmd *cobra.Command) {
	fmt.Printf("Error: invalid %s, check 'pqcom config list'\n\n", what)
	cmd.Help()
	fmt.Println()
}

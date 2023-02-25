package cmd

import (
	"github.com/SamoKopecky/pqcom/main/config"
	"github.com/SamoKopecky/pqcom/main/crypto"
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
			config.GenerateConfig(kem, sign)
			// TODO: validate given kem, sign
		},
	}
)

func init() {
	configCmd.AddCommand(genConfig)
	genConfig.Flags().StringVarP(&kem, "kem", "k", crypto.GetAllKems()[0], "Use kem")
	genConfig.Flags().StringVarP(&sign, "sign", "s", crypto.GetAllSigns()[0], "Use sign")
}

package cmd

import (
	"github.com/SamoKopecky/pqcom/main/benchmark"
	"github.com/spf13/cobra"
)

var (
	iterations int

	benchmarkCmd = &cobra.Command{
		Use:   "benchmark",
		Short: "Use benchmark mode",
		Long:  `Benchmark all available algorithms`,
		Run: func(cmd *cobra.Command, args []string) {
			benchmark.Run(iterations)
		},
	}
)

func init() {
	rootCmd.AddCommand(benchmarkCmd)
	benchmarkCmd.Flags().IntVarP(&iterations, "iternations", "i", 1000, "number of iterations")
}

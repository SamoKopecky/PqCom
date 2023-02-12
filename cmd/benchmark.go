package cmd

import (
	"github.com/SamoKopecky/pqcom/main/benchmark"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(benchmarkCmd)
	benchmarkCmd.Flags().IntVarP(&iterations, "iternations", "i", 1000, "Number of iterations")
}

var (
	iterations int

	benchmarkCmd = &cobra.Command{
		Use:   "benchmark",
		Short: "Use benchmark mode",
		Long:  `Benchmark Kyber and Dilithium to some other implementations`,
		Run: func(cmd *cobra.Command, args []string) {
			benchmark.Run(iterations)
		},
	}
)

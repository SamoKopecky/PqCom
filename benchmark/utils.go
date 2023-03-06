package benchmark

import (
	"fmt"
	"time"

	"github.com/SamoKopecky/pqcom/main/crypto"
)

func timeKem(toTest func(crypto.KemAlgorithm), arg crypto.KemAlgorithm, algName string, iterations int) {
	fmt.Printf("Benchmarking %s...\n", algName)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		toTest(arg)
	}
	print(parseElapsed(time.Since(start), algName, iterations))
}

func timeSign(toTest func(crypto.SignAlgorithm), arg crypto.SignAlgorithm, algName string, iterations int) {
	fmt.Printf("Benchmarking %s...\n", algName)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		toTest(arg)
	}
	print(parseElapsed(time.Since(start), algName, iterations))
}

func parseElapsed(elapsed time.Duration, functionName string, iterations int) string {
	average := (float64(elapsed.Nanoseconds()) / 1000) / float64(iterations)
	total := elapsed.Seconds()
	return fmt.Sprintf("Benchmark for %s took %.4f s, one iteration on average %.4f us\n", functionName, total, average)
}

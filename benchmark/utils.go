package benchmark

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func timeFunction(toTest func(), iterations int) {
	functionName := getFunctionName(toTest)
	fmt.Printf("Benchmarking %s...\n", functionName)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		toTest()
	}
	print(parseElapsed(time.Since(start), functionName, iterations))
}

func parseElapsed(elapsed time.Duration, functionName string, iterations int) string {
	average := (float64(elapsed.Nanoseconds()) / 1000) / float64(iterations)
	total := elapsed.Seconds()
	return fmt.Sprintf("Benchmark for %s took %.4f s, one iteration on average %.4f us.\n", functionName, total, average)
}

func getFunctionName(i interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return strings.Split(fullName, ".")[1]
}

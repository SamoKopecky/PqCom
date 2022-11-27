package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func timeFunction(toTest func(), iterations int) {
	functionName := GetFunctionName(toTest)
	fmt.Printf("Benchmarking %s...\n", functionName)
	start := time.Now()
	for i := 0; i < iterations; i++ {
		toTest()
	}
	print(parseElapsed(time.Since(start), functionName, iterations))
}

func parseElapsed(elapsed time.Duration, functionName string, iterations int) string {
	return fmt.Sprintf("Benchmark for %s took %.4f seconds.\n", functionName, elapsed.Seconds())
}

func GetFunctionName(i interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return strings.Split(fullName, ".")[1]
}

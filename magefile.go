//go:build mage

package main 

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	GOOS = getEnv("GOOS", runtime.GOOS)
	GOARCH = getEnv("GOARCH", runtime.GOARCH)
)

func getEnv(key string, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Build() error {
	mg.Deps(Clean) 

	fmt.Printf("Building all services for %s/%s...\n", GOOS, GOARCH)
	services := []string{"service-1", "service-2"}
	
	for _, svc := range services {
		if err := buildService(svc); err != nil {
			return fmt.Errorf("failed to build %s: %w", svc, err)
		}
	}

	mg.Deps(Package)
	return nil 
}

func buildService(svcName string) error {
	svcPath := svcName
	outputName := svcName
	if GOOS == "windows" {
		outputName += ".exe"
	}

	outputPath := filepath.Join("dist", outputName)

	env := map[string]string {
		"GOOS": GOOS, 
		"GOARCH": GOARCH, 
		"CGO_ENABLED": "0",
	}

	return sh.RunWith(env, "go", "build", "-o", outputPath, fmt.Sprintf("./%s", svcPath))
}


func Demo() error {
	mg.Deps(Build)
	
	fmt.Println("=== Service Integration Demo ===")
	
	fmt.Println("1. Testing Binary Calculator Service:")
	demos := [][]string{
		{"add", "15", "25"},
		{"sub", "50", "20"},
		{"mul", "7", "8"},
		{"div", "100", "5"},
	}
	
	for _, demo := range demos {
		fmt.Printf("   > ./service-1 %s\n", joinArgs(demo))
		if err := sh.Run("./dist/service-1", demo...); err != nil {
			fmt.Printf("     Error: %v\n", err)
		}
		fmt.Println()
	}
	
	fmt.Println("2. Testing String Reverser Service:")
	stringDemos := [][]string{
		{"reverse", "hello"},
		{"reverse", "GitHub Actions"},
		{"wordcount", "the quick brown fox"},
		{"wordcount", "microservices are cool"},
	}
	
	for _, demo := range stringDemos {
		fmt.Printf("   > ./service-2 %s\n", joinArgs(demo))
		if err := sh.Run("./dist/service-2", demo...); err != nil {
			fmt.Printf("     Error: %v\n", err)
		}
		fmt.Println()
	}
	
	fmt.Println("=== Demo Complete ===")
	return nil
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		// Add quotes if the argument contains spaces
		if len(arg) > 0 && (contains(arg, ' ') || contains(arg, '\t')) {
			result += `"` + arg + `"`
		} else {
			result += arg
		}
	}
	return result
}

func contains(s string, c rune) bool {
	for _, r := range s {
		if r == c {
			return true
		}
	}
	return false
}


func Test() error {
	mg.Deps(Build)
	
	fmt.Println("Running integration tests...")
	
	fmt.Println("Testing Binary Calculator Service...")
	tests := []struct {
		args     []string
		shouldSucceed bool
	}{
		{[]string{"add", "2", "3"}, true},
		{[]string{"mul", "4", "5"}, true},
		{[]string{"div", "10", "2"}, true},
		{[]string{"sub", "8", "3"}, true},
	}
	
	for _, test := range tests {
		fmt.Printf("  Testing: service-1 %s\n", joinArgs(test.args))
		err := sh.Run("./dist/service-1", test.args...)
		if test.shouldSucceed && err != nil {
			return fmt.Errorf("test failed for service-1 %v: %w", test.args, err)
		}
		fmt.Println()
	}
	
	fmt.Println("Testing String Reverser Service...")
	stringTests := []struct {
		args     []string
		shouldSucceed bool
	}{
		{[]string{"reverse", "test"}, true},
		{[]string{"reverse", "hello world"}, true},
		{[]string{"wordcount", "one two three"}, true},
		{[]string{"wordcount", "single"}, true},
	}
	
	for _, test := range stringTests {
		fmt.Printf("  Testing: service-2 %s\n", joinArgs(test.args))
		err := sh.Run("./dist/service-2", test.args...)
		if test.shouldSucceed && err != nil {
			return fmt.Errorf("test failed for service-2 %v: %w", test.args, err)
		}
		fmt.Println()
	}
	
	fmt.Println("All tests passed!")
	return nil
}

func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	return sh.Rm("dist")
}

func Package() error {
	fmt.Println("Creating packages for each platform")
	unixScript := `
	echo "_____Integration demo_____" 
	./service-1 add 10 20
	./service-1 mul 6 7	
	./service-2 reverse "hello"
	./service-2 wordcount "one two three"
	echo "__________________________"
	`

	if err := os.WriteFile("dist/unixScript.sh", []byte(unixScript), 0755); err != nil { 
		return fmt.Errorf("failed to create unix script artifact: %w", err)
	}

	winScript := `@echo off 
	echo _____Microservices Demo_____
	service-1.exe add 10 20
	service-1.exe mul 6 7
	service-2.exe reverse "hello"
	service-2.exe wordcount "one two three"
	echo __________________________
	pause
	`

	if err := os.WriteFile("dist/demo.bat", []byte(winScript), 0755); err != nil {
		return fmt.Errorf("failed to create Windows demo script: %w", err)
	}
	
	fmt.Println("Package created with demo scripts")
	return nil
}
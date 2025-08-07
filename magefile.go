//go:build mage

package main

import (
	"fmt"
	"os"
	
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Build() error {
	mg.Deps(Clean)
	
	services := []string{"service-1", "service-2"}
	
	for _, service := range services {
		fmt.Printf("Building %s...\n", service)
		if err := sh.Run("go", "build", "-o", "dist/"+service, "./"+service); err != nil {
			return fmt.Errorf("failed to build %s: %w", service, err)
		}
	}
	
	fmt.Println("Build complete!")
	return nil
}

func Clean() error {
	fmt.Println("Cleaning...")
	return sh.Rm("dist")
}

func Test() error {
	mg.Deps(Build)
	
	fmt.Println("Testing service-1...")
	if err := sh.Run("./dist/service-1", "add", "2", "3"); err != nil {
		return err
	}
	
	fmt.Println("Testing service-2...")
	if err := sh.Run("./dist/service-2", "reverse", "hello"); err != nil {
		return err
	}
	
	fmt.Println("All tests passed!")
	return nil
}

func Demo() error {
	mg.Deps(Build)
	
	fmt.Println("=== Demo ===")
	
	fmt.Println("Calculator:")
	sh.Run("./dist/service-1", "add", "15", "25")
	sh.Run("./dist/service-1", "mul", "7", "8")
	
	fmt.Println("\nString processor:")
	sh.Run("./dist/service-2", "reverse", "hello")
	sh.Run("./dist/service-2", "wordcount", "the quick brown fox")
	
	fmt.Println("\n=== Demo Complete ===")
	return nil
}

func Install() error {
	mg.Deps(Build)
	
	if err := os.MkdirAll("dist", 0755); err != nil {
		return err
	}
	
	unixScript := `#!/bin/bash
	echo "=== Service Demo ==="
	./service-1 add 10 20
	./service-1 mul 6 7
	./service-2 reverse "hello world"
	./service-2 wordcount "one two three"
	echo "===================="
	`
	if err := os.WriteFile("dist/demo.sh", []byte(unixScript), 0755); err != nil {
		return err
	}
	
	winScript := `@echo off
	echo === Service Demo ===
	service-1.exe add 10 20
	service-1.exe mul 6 7
	service-2.exe reverse "hello world"
	service-2.exe wordcount "one two three"
	echo ====================
	pause
	`
	if err := os.WriteFile("dist/demo.bat", []byte(winScript), 0644); err != nil {
		return err
	}
	
	fmt.Println("Demo scripts created in dist/")
	return nil
}
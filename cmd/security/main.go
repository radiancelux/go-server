package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	fmt.Println("🔒 Running Security Check for Go Dependencies...")

	// Check if govulncheck is available
	if err := checkGovulncheck(); err != nil {
		log.Fatalf("❌ govulncheck check failed: %v", err)
	}

	// Run vulnerability scan
	fmt.Println("🔍 Running vulnerability scan...")
	if err := runCommand("govulncheck", "./..."); err != nil {
		log.Fatalf("❌ Vulnerability scan failed: %v", err)
	}
	fmt.Println("✅ No vulnerabilities found!")

	// Run go mod tidy
	fmt.Println("🔍 Running go mod tidy...")
	if err := runCommand("go", "mod", "tidy"); err != nil {
		log.Fatalf("❌ go mod tidy failed: %v", err)
	}

	// Run go vet
	fmt.Println("🔍 Running go vet...")
	if err := runCommand("go", "vet", "./..."); err != nil {
		log.Fatalf("❌ go vet failed: %v", err)
	}
	fmt.Println("✅ Code analysis passed!")

	// Check code formatting
	fmt.Println("🔍 Checking code formatting...")
	if err := checkCodeFormatting(); err != nil {
		log.Fatalf("❌ Code formatting check failed: %v", err)
	}
	fmt.Println("✅ Code formatting is correct!")

	fmt.Println("🎉 Security check completed successfully!")
}

func checkGovulncheck() error {
	// Check if govulncheck is in PATH
	if _, err := exec.LookPath("govulncheck"); err != nil {
		fmt.Println("❌ govulncheck not found. Installing...")

		// Install govulncheck
		cmd := exec.Command("go", "install", "golang.org/x/vuln/cmd/govulncheck@latest")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to install govulncheck: %w", err)
		}

		// Add Go bin to PATH
		goPath := os.Getenv("GOPATH")
		if goPath == "" {
			// Default Go path
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			goPath = filepath.Join(homeDir, "go")
		}

		goBinPath := filepath.Join(goPath, "bin")
		if err := os.Setenv("PATH", os.Getenv("PATH")+string(os.PathListSeparator)+goBinPath); err != nil {
			return fmt.Errorf("failed to update PATH: %w", err)
		}
	}
	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func checkCodeFormatting() error {
	cmd := exec.Command("go", "fmt", "./...")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("go fmt failed: %w", err)
	}

	if len(output) > 0 {
		fmt.Printf("❌ Code formatting issues found:\n%s", string(output))
		fmt.Println("Run 'go fmt ./...' to fix formatting issues.")
		return fmt.Errorf("code formatting issues detected")
	}

	return nil
}

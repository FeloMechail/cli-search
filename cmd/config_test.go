package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	output := new(bytes.Buffer)
	rootCmd.SetOut(output)

	rootCmd.SetArgs([]string{"config", "--showpath"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Error reading config file %v", err)
	}
}

func TestLoadConfigError(t *testing.T) {
	output := new(bytes.Buffer)
	rootCmd.SetOut(output)

	// originalConfig
	originalConfig, err := os.ReadFile("config.yaml")
	if err != nil {
		t.Fatalf("Failed to read the original config file: %v", err)
	}

	defer func() {
		err := os.WriteFile("config.yaml", originalConfig, os.ModePerm)
		if err != nil {
			t.Fatalf("Failed to restore the original config file: %v", err)
		}
	}()

	err = os.Remove("config.yaml")
	if err != nil {
		t.Fatalf("Error removing config file for testing %v", err)
	}

	rootCmd.SetArgs([]string{"config"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Error: no error reported for config file %v", err)
	}
}

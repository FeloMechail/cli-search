package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SearchEngines  []SearchEngine `yaml:"search_engines"`
	DefaultSearch  string         `yaml:"default_search"`
	DefaultBrowser string         `yaml:"default_browser"`
}

type SearchEngine struct {
	Name     string `yaml:"name"`
	Shortcut string `yaml:"shortcut"`
	URL      string `yaml:"url"`
}

var configPath string = "cmd/config.yaml"

// loadConfig function  
func LoadConfig() (*Config, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	if config.DefaultBrowser == "" {
		browser, err := exec.Command("xdg-settings", "get", "default-web-browser").
			Output()
		if err != nil {
			log.Fatal(err)
		}

		config.DefaultBrowser = strings.TrimSpace(string(browser))

		updatedData, err := yaml.Marshal(&config)
		if err != nil {
			log.Fatalf("Error marshaling yaml: %v", err)
		}

		err = os.WriteFile(configPath, updatedData, os.ModePerm)
		if err != nil {
			log.Fatalf("Error writing default browser: %v\n", err)
		}

	}

	return &config, nil
}

func SetDefaultBrowser(browser string) error {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %v\n", err)
	}
	// TODO: clean browser string and check if browser exists
	config.DefaultBrowser = browser
	fmt.Printf("Changed %s to default browser\n", browser)

	updatedData, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("Error marshaling yaml: %v", err)
	}

	err = os.WriteFile(configPath, updatedData, os.ModePerm)
	if err != nil {
		log.Fatalf("Error writing default browser: %v\n", err)
	}

	return nil
}

func SetDefaultSearchEngine(engine string) error {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config file: %v\n", err)
	}

	for _, name := range config.SearchEngines {
		if name.Name == engine {
			config.DefaultSearch = name.Shortcut
			updatedData, err := yaml.Marshal(&config)
			if err != nil {
				log.Fatalf("Error marshaling yaml: %v", err)
			}

			err = os.WriteFile(configPath, updatedData, os.ModePerm)
			if err != nil {
				log.Fatalf("Error writing default browser: %v\n", err)
			}
			fmt.Printf("Changed %s to default search engine\n", engine)
			return nil

		}
	}
	return errors.New("Search Engine not in config file, add it using ..")
}

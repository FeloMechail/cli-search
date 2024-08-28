package cmd

import (
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

// loadConfig function  
func loadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
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

		err = os.WriteFile(filename, updatedData, os.ModePerm)
		if err != nil {
			log.Fatalf("Error writing default browser: %s\n", err)
		}

	}

	return &config, nil
}

// func main() {
// 	config, err := loadConfig("config.yaml")
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}
//
// 	for _, engine := range config.SearchEngines {
// 		fmt.Printf(
// 			"Name: %s, Shortcut: %s, URL: %s\n",
// 			engine.Name,
// 			engine.Shortcut,
// 			engine.URL,
// 		)
// 	}
// }

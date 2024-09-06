package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the structure of the configuration file.
type Config struct {
	SearchEngines  []SearchEngine `yaml:"search_engines"  mapstructure:"search_engines"`  // List of search engines
	DefaultSearch  string         `yaml:"default_search"  mapstructure:"default_search"`  // Shortcut for the default search engine
	DefaultBrowser string         `yaml:"default_browser" mapstructure:"default_browser"` // Default browser
}

// SearchEngine represents a search engine configuration.
type SearchEngine struct {
	Name     string `yaml:"name"     mapstructure:"name"`     // Name of the search engine
	Shortcut string `yaml:"shortcut" mapstructure:"shortcut"` // Shortcut for the search engine
	URL      string `yaml:"url"      mapstructure:"url"`      // URL template for the search engine
}

var (
	configPath   string
	config       Config
	urlMap       map[string]string
	configLoaded bool
)

func LoadConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get home directory: %w", err)
	}

	configPath = filepath.Join(homeDir, ".config", "cs", "cs.yaml")

	viper.SetConfigFile(configPath)

	// Attempt to read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Handle case where config file is not found
		if ok := os.IsNotExist(err); ok {
			fmt.Println(
				"Config file not found, creating a new one with default settings...",
			)
			if err := createDefaultConfig(); err != nil {
				return fmt.Errorf("failed to create default config: %w", err)
			}

			if err := viper.ReadInConfig(); err != nil {
				return fmt.Errorf("Error re-reading default config: %w", err)
			}

		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal the config into the struct
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to unmarshal config: %w", err)
	}

	// Populate the urlMap for easy access
	urlMap = make(map[string]string)
	for _, engine := range config.SearchEngines {
		urlMap[engine.Shortcut] = engine.URL
	}

	configLoaded = true
	return nil
}

func createDefaultConfig() error {
	// TODO: get system default browser
	defaultConfig := Config{
		SearchEngines: []SearchEngine{
			{
				Name:     "Google",
				Shortcut: "g",
				URL:      "https://www.google.com/search?q=",
			},
			{
				Name:     "Wikipedia",
				Shortcut: "wiki",
				URL:      "https://en.wikipedia.org/wiki/Special:Search?search=",
			},
		},
		DefaultSearch:  "g",
		DefaultBrowser: "",
	}

	// TODO: OS detection!!!!!
	if output, err := exec.Command("xdg-settings", "get", "default-web-browser").
		Output(); err != nil {
		err = fmt.Errorf(
			"Could not get default browser info.. add it manually, %v",
			err,
		)
		fmt.Println(err)
	} else {
		defaultConfig.DefaultBrowser = strings.TrimSpace(string(output))
	}

	fmt.Printf("Default configuration: %v\n\n", defaultConfig)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".config", "cs")
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("unable to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "cs.yaml")
	file, err := os.Create(configFile)
	if err != nil {
		return fmt.Errorf("unable to create config file: %w", err)
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(defaultConfig); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}

func SetDefaultBrowser(browser string) error {
	// Validate and clean browser string
	config.DefaultBrowser = browser
	fmt.Printf("Changed default browser to %s\n", browser)

	return saveConfig()
}

func SetDefaultSearchEngine(engine string) error {
	for _, se := range config.SearchEngines {
		if se.Name == engine {
			config.DefaultSearch = se.Shortcut
			fmt.Printf("Changed default search engine to %s\n", engine)
			return saveConfig()
		}
	}
	return fmt.Errorf("search engine '%s' not found in config", engine)
}

func saveConfig() error {
	updatedData, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file for writing: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(updatedData); err != nil {
		return fmt.Errorf("failed to write updated config to file: %w", err)
	}

	return nil
}

func PerformSearch(search string, flags []string) (string, error) {
	// TODO: search engine flag
	var url string
	if slices.Contains(flags, "u") {
		url = "https://" + search
	} else if slices.Contains(flags, "e") {
		if url = urlMap[engine]; url == "" {
			return "", errors.New("Engine not in config. add it")
		}

		url = url + search
	} else {
		url = urlMap[config.DefaultSearch]
		url = url + search
	}

	fmt.Printf("URL: %s\n", url)
	output, err := openBrowser(url)

	return string(output), err
}

func showConfigPath() {
	fmt.Println(configPath)
}

func showConfig() {
	for _, se := range config.SearchEngines {
		fmt.Printf(
			"Search engine: %s, Shortcut: %s, url: %s\n",
			se.Name,
			se.Shortcut,
			se.URL,
		)
	}
	fmt.Printf(
		"Default search: %s, Default browser: %s\n",
		config.DefaultSearch,
		config.DefaultBrowser,
	)
}

func openBrowser(query string) (output []byte, err error) {
	// TODO: change to whatever browser in config, OS detection
	cmd := "xdg-open"
	var args []string
	fmt.Print("Opening browser\n")

	args = append(args, query)
	output, err = exec.Command(cmd, args...).CombinedOutput()
	return output, err
}

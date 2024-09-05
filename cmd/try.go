package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
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

var (
	configPath   string = "cmd/config.yaml"
	config       Config
	urlMap       map[string]string
	configLoaded bool
)

// loadConfig function  
func LoadConfig() error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return err
	}

	if config.DefaultBrowser == "" {
		// TODO: execute command depending on OS
		browser, err := exec.Command("xdg-settings", "get", "default-web-browser").
			Output()
		if err != nil {
			return err
		}

		config.DefaultBrowser = strings.TrimSpace(string(browser))

		updatedData, err := yaml.Marshal(&config)
		if err != nil {
			return err
		}

		err = os.WriteFile(configPath, updatedData, os.ModePerm)
		if err != nil {
			return err
		}

	}

	urlMap = make(map[string]string)

	for _, enginess := range config.SearchEngines {
		urlMap[enginess.Shortcut] = enginess.URL
	}

	configLoaded = true

	return nil
}

func SetDefaultBrowser(browser string) error {
	// TODO: clean browser string and check if browser exists
	config.DefaultBrowser = browser
	fmt.Printf("Changed %s to default browser\n", browser)

	updatedData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, updatedData, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func SetDefaultSearchEngine(engine string) error {
	// TODO: hash table?
	for _, name := range config.SearchEngines {
		if name.Name == engine {
			config.DefaultSearch = name.Shortcut
			updatedData, err := yaml.Marshal(&config)
			if err != nil {
				return err
			}

			err = os.WriteFile(configPath, updatedData, os.ModePerm)
			if err != nil {
				return err
			}
			fmt.Printf("Changed %s to default search engine\n", engine)
			return nil

		}
	}
	return errors.New("Search Engine not in config file, add it using ..")
}

func PerformSearch(search string, flags []string) (string, error) {
	// TODO: search engine flag
	var url string
	if slices.Contains(flags, "u") {
		url = "https://" + search
	} else if slices.Contains(flags, "e") {
		if url = urlMap[engine]; url == "" {
			log.Fatal("Engine not in config. add it")
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

// TODO: show absolute path
func showConfigPath() {
	fmt.Print("PATH: cmd/config.yaml\n")
}

func showConfig() {
	for _, engine := range config.SearchEngines {
		fmt.Printf(
			"Name: %s, Shortcut: %s, URL: %s\n",
			engine.Name,
			engine.Shortcut,
			engine.URL,
		)
	}
	fmt.Printf(
		"Default Engine: %s, Default Browser: %s",
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

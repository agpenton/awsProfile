package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
)

var profileRegex = regexp.MustCompile(`\[profile .*]`)
var bracketsRemovalRegx = regexp.MustCompile(`(\[profile )|(\])`)
var defaultProfileChoice = "default"

func main() {
	fmt.Println("AWS Profile Switcher")
	homeDir, _ := os.UserHomeDir()

	data, err := ioutil.ReadFile(homeDir + "/.aws/config")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	matches := profileRegex.FindAllString(string(data), -1)

	if len(matches) == 0 {
		fmt.Println("No profiles found.")
		fmt.Println("Refer to this guide for help on setting up a new AWS profile:")
		fmt.Println("https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html")
		os.Exit(1)
	}

	profiles := make([]string, 0, len(matches))
	for _, match := range matches {
		profiles = append(profiles, bracketsRemovalRegx.ReplaceAllString(match, ""))
	}
	profiles = append(profiles, defaultProfileChoice)

	prompt := promptui.Select{
		Label:             "Choose a profile",
		Items:             profiles,
		StartInSearchMode: true,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile(homeDir+"/.awsp", []byte(result), 0644); err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}
}

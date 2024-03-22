package cmd

import (
	"fmt"
	"os"
	"os/user"

	"flag"

	"github.com/dhth/act3/ui"
)

func die(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func Execute() {
	currentUser, err := user.Current()
	var defaultConfigFilePath string
	if err == nil {
		defaultConfigFilePath = fmt.Sprintf("%s/.config/act3/act3.yml", currentUser.HomeDir)
	}
	configFilePath := flag.String("config-file", defaultConfigFilePath, "path of the config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *configFilePath == "" {
		die("config-file cannot be empty")
	}

	configFilePathExpanded := expandTilde(*configFilePath)

	_, err = os.Stat(configFilePathExpanded)
	if os.IsNotExist(err) {
		die(cfgErrSuggestion(fmt.Sprintf("Error: file doesn't exist at %q", configFilePathExpanded)))
	}

	workflows, err := ReadConfig(configFilePathExpanded)
	if err != nil {
		die(cfgErrSuggestion(fmt.Sprintf("Error reading config: %v", configFilePathExpanded)))
	}
	if len(workflows) == 0 {
		die(cfgErrSuggestion(fmt.Sprintf("No workflows found")))
	}

	ghClient, err := getGHClient()
	if err != nil {
		die("Error: %q", err.Error())
	}

	ui.RenderUI(ghClient, workflows)
}

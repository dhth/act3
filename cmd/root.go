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

var (
	format           = flag.String("format", "", "output format to use; possible values: html")
	htmlTemplateFile = flag.String("html-template-file", "", "path of the HTML template file to use")
)

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

	var outputFmt ui.OutputFmt
	if *format != "" {
		switch *format {
		case "html":
			outputFmt = ui.HTMLFmt
		default:
			die("unsupported value for format")
		}
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
		die(cfgErrSuggestion("No workflows found"))
	}

	var htmlTemplate string
	if *htmlTemplateFile != "" {
		_, err := os.Stat(*htmlTemplateFile)
		if os.IsNotExist(err) {
			die(fmt.Sprintf("Error: template file doesn't exist at %q", *htmlTemplateFile))
		}
		templateFileContents, err := os.ReadFile(*htmlTemplateFile)
		if err != nil {
			die(fmt.Sprintf("Error: couldn't read template file %q", *htmlTemplateFile))
		}
		htmlTemplate = string(templateFileContents)
	}

	ghClient, err := getGHClient()
	if err != nil {
		die("Error: %q", err.Error())
	}

	ui.RenderUI(ghClient, workflows, outputFmt, htmlTemplate)
}

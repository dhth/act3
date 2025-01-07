package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/types"
)

const (
	configPath      = "act3/act3.yml"
	author          = "@dhth"
	projectHomePage = "https://github.com/dhth/act3"
	issuesURL       = "https://github.com/dhth/act3/issues"
)

var (
	errFlagCombIncorrect       = errors.New("flag combination incorrect")
	errIncorrectRepoProvided   = errors.New("incorrect repo provided")
	errCouldntGetConfigDir     = errors.New("couldn't get your config directory")
	errConfigFilePathEmpty     = errors.New("config file path is empty")
	errIncorrectOutputFmt      = errors.New("incorrect value for output format provided")
	errConfigFileDoesntExit    = errors.New("config file doesn't exist")
	errCouldntReadConfig       = errors.New("couldn't read config")
	errCouldntGetGHClient      = errors.New("couldn't get a Github client")
	errNoWorkflows             = errors.New("no workflows found")
	errTemplateFileDoesntExit  = errors.New("template file doesn't exist")
	errCouldntReadTemplateFile = errors.New("couldn't read template file")
	errCouldntGetWorkflows     = errors.New("couldn't get workflows")
)

var (
	format           = flag.String("f", "", "output format to use; possible values: default, table, html")
	htmlTemplateFile = flag.String("t", "", "path of the HTML template file to use")
	global           = flag.Bool("g", false, "whether to use workflows defined globally via the config file")
	repo             = flag.String("r", "", "repo to fetch worflows for, in the format \"owner/repo\"")
)

func Execute() error {
	var defaultConfigDir string
	var configErr error
	switch runtime.GOOS {
	case "linux", "windows":
		defaultConfigDir, configErr = os.UserConfigDir()
	default:
		hd, configErr := os.UserHomeDir()
		if configErr != nil {
			break
		}
		defaultConfigDir = filepath.Join(hd, ".config")
	}
	if configErr != nil {
		fmt.Printf(`Couldn't get your default config directory. This is a fatal error;
use -c to specify config file path manually.
Let %s know about this via %s.
`, author, issuesURL)
		return fmt.Errorf("%w: %s", errCouldntGetConfigDir, configErr.Error())
	}
	defaultConfigFilePath := filepath.Join(defaultConfigDir, configPath)
	configFilePath := flag.String("c", defaultConfigFilePath, "path of the config file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n\nFlags:\n", helpText)
		flag.PrintDefaults()
	}

	flag.Parse()

	// flag validation
	if *configFilePath == "" {
		return fmt.Errorf("%w", errConfigFilePathEmpty)
	}

	if *global && *repo != "" {
		return fmt.Errorf("%w; -g and -r cannot both be provided at the same time", errFlagCombIncorrect)
	}

	if *repo != "" {
		repoEls := strings.Split(*repo, "/")
		if len(repoEls) != 2 {
			return fmt.Errorf("%w; repo needs to be in the format \"owner/repo\"", errIncorrectRepoProvided)
		}
	}

	var outputFmt types.OutputFmt
	if *format != "" {
		switch *format {
		case "default":
			outputFmt = types.DefaultFmt
		case "table":
			outputFmt = types.TableFmt
		case "html":
			outputFmt = types.HTMLFmt
		default:
			return fmt.Errorf("%w", errIncorrectOutputFmt)
		}
	}

	var htmlTemplate string
	if *htmlTemplateFile != "" {
		_, err := os.Stat(*htmlTemplateFile)
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: path: %s", errTemplateFileDoesntExit, *htmlTemplateFile)
		}
		templateFileContents, err := os.ReadFile(*htmlTemplateFile)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntReadTemplateFile, err.Error())
		}
		htmlTemplate = string(templateFileContents)
	}

	clientOpts := ghapi.ClientOptions{
		EnableCache: false,
		Timeout:     8 * time.Second,
	}

	var workflows []types.Workflow
	var currentRepo string
	var err error

	if *global {
		configFilePathExpanded := expandTilde(*configFilePath)

		_, err = os.Stat(configFilePathExpanded)
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: path: %s", errConfigFileDoesntExit, configFilePathExpanded)
		}

		workflows, err = ReadConfig(configFilePathExpanded)
		if err != nil {
			fmt.Print(configSampleFormat)
			return fmt.Errorf("%w: %s", errCouldntReadConfig, err.Error())
		}

	} else {
		if *repo != "" {
			currentRepo = *repo
		} else {
			currentRepo, err = getCurrentRepo()
			if err != nil {
				return err
			}
		}
		ghRClient, err := ghapi.NewRESTClient(clientOpts)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntGetGHClient, err.Error())
		}

		workflows, err = getWorkflowsForCurrentRepo(ghRClient, currentRepo)
		if err != nil {
			return fmt.Errorf("%w: %s", errCouldntGetWorkflows, err.Error())
		}
	}

	if len(workflows) == 0 {
		return fmt.Errorf("%w", errNoWorkflows)
	}

	ghClient, err := ghapi.NewGraphQLClient(clientOpts)
	if err != nil {
		return fmt.Errorf("%w: %s", errCouldntGetGHClient, err.Error())
	}

	var cr *string
	if !*global {
		cr = &currentRepo
	}
	config := types.Config{
		GHClient:     ghClient,
		CurrentRepo:  cr,
		Fmt:          outputFmt,
		HTMLTemplate: htmlTemplate,
	}

	err = render(workflows, config)
	if err != nil {
		return err
	}
	return nil
}

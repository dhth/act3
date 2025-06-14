package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ghapi "github.com/cli/go-gh/v2/pkg/api"
	"github.com/dhth/act3/internal/types"
	"github.com/dhth/act3/internal/utils"
	"github.com/spf13/cobra"
)

const (
	configFileName = "act3/act3.yml"
)

var (
	errCouldntGetUserHomeDir   = errors.New("couldn't get your home directory")
	errCouldntGetUserConfigDir = errors.New("couldn't get your config directory")
	errCouldntReadConfigFile   = errors.New("couldn't read config file")
	errFlagCombIncorrect       = errors.New("flag combination incorrect")
	errIncorrectOutputFmt      = errors.New("incorrect value for output format provided")
	ErrConfigFileDoesntExit    = errors.New("config file doesn't exist")
	errCouldntReadConfig       = errors.New("couldn't read config")
	errCouldntGetGHClient      = errors.New("couldn't get a Github client")
	errNoWorkflows             = errors.New("no workflows found")
	errTemplateFileDoesntExit  = errors.New("template file doesn't exist")
	errCouldntReadTemplateFile = errors.New("couldn't read template file")
	errCouldntGetWorkflows     = errors.New("couldn't get workflows")
	errInvalidRepoProvided     = errors.New("invalid repo provided")
	errReposProvided           = errors.New("invalid repos provided")
)

func Execute() error {
	rootCmd, err := NewRootCommand()
	if err != nil {
		return err
	}

	return rootCmd.Execute()
}

func NewRootCommand() (*cobra.Command, error) {
	var (
		configPath       string
		configBytes      []byte
		homeDir          string
		globalWorkflows  bool
		reposStr         string
		formatStr        string
		htmlTemplateFile string
		htmlTitle        string
		openFailed       bool
	)

	rootCmd := &cobra.Command{
		Use:          "act3",
		Short:        "Glance at the last 3 runs of your GitHub Actions workflows",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			if globalWorkflows && reposStr != "" {
				return fmt.Errorf("%w; -g and -r cannot both be provided at the same time", errFlagCombIncorrect)
			}

			if reposStr != "" {
				reposVal := strings.SplitSeq(reposStr, ",")
				var invalidRepos []string
				for r := range reposVal {
					if !utils.IsRepoNameValid(r) {
						invalidRepos = append(invalidRepos, r)
					}
				}
				if len(invalidRepos) == 1 {
					return fmt.Errorf(`%w: %q; value needs to be in the format "owner/repo"`, errInvalidRepoProvided, invalidRepos[0])
				}

				if len(invalidRepos) > 1 {
					return fmt.Errorf(`%w: %q; value needs to be in the format "owner/repo"`, errReposProvided, invalidRepos)
				}
			}

			var outputFmt types.OutputFmt
			if formatStr != "" {
				switch formatStr {
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
			if htmlTemplateFile != "" {
				_, err := os.Stat(htmlTemplateFile)
				if os.IsNotExist(err) {
					return fmt.Errorf("%w: path: %s", errTemplateFileDoesntExit, htmlTemplateFile)
				}
				templateFileContents, err := os.ReadFile(htmlTemplateFile)
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
			var reposToUse []string

			if globalWorkflows {
				configPathFull := utils.ExpandTilde(configPath, homeDir)

				var cfgErr error
				_, cfgErr = os.Stat(configPathFull)
				if os.IsNotExist(cfgErr) {
					return fmt.Errorf("%w, create one at %q", ErrConfigFileDoesntExit, configPathFull)
				}

				configPathFull = utils.ExpandTilde(configPath, homeDir)
				configBytes, cfgErr = os.ReadFile(configPathFull)
				if cfgErr != nil {
					return fmt.Errorf("%w: %w", errCouldntReadConfigFile, cfgErr)
				}
				workflows, cfgErr = ReadConfig(configBytes)
				if cfgErr != nil {
					fmt.Print(configSampleFormat)
					return fmt.Errorf("%w: %s", errCouldntReadConfig, cfgErr.Error())
				}

			} else {
				if reposStr != "" {
					reposToUse = strings.Split(reposStr, ",")
				} else {
					currentRepo, err := getCurrentRepo()
					if err != nil {
						return err
					}
					reposToUse = []string{currentRepo}
				}

				ghRClient, err := ghapi.NewRESTClient(clientOpts)
				if err != nil {
					return fmt.Errorf("%w: %s", errCouldntGetGHClient, err.Error())
				}

				var errors []WorkflowError
				workflows, errors = getWorkflowsForRepos(ghRClient, reposToUse)

				if len(errors) == 1 {
					return fmt.Errorf("%w:\n%s", errCouldntGetWorkflows, errors[0].Err.Error())
				}

				if len(errors) > 1 {
					errorStrs := make([]string, len(errors))
					for i, e := range errors {
						errorStrs[i] = fmt.Sprintf("- %s: %s", e.Repo, e.Err.Error())
					}

					return fmt.Errorf("%w:\n%s", errCouldntGetWorkflows, strings.Join(errorStrs, "\n"))
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
			if !globalWorkflows && len(reposToUse) == 1 {
				cr = &reposToUse[0]
			}
			config := types.Config{
				GHClient:     ghClient,
				CurrentRepo:  cr,
				Fmt:          outputFmt,
				HTMLTemplate: htmlTemplate,
				HTMLTitle:    htmlTitle,
			}

			results := getResults(workflows, config)

			err = render(results, config)
			if err != nil {
				return err
			}

			if openFailed {
				openFailedWorkflows(results)
			}
			return nil
		},
	}

	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserHomeDir, err.Error())
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errCouldntGetUserConfigDir, err.Error())
	}

	defaultConfigPath := filepath.Join(configDir, configFileName)

	rootCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, "location of ecsv's config file")
	rootCmd.Flags().StringVarP(&reposStr, "repos", "r", "", `comma delimited list of repos to fetch workflows for, in the format "owner/repo"`)
	rootCmd.Flags().BoolVarP(&globalWorkflows, "global", "g", false, "whether to use workflows defined globally via the config file")
	rootCmd.Flags().StringVarP(&formatStr, "format", "f", "default", "output format to use; possible values: default, table, html")
	rootCmd.Flags().StringVar(&htmlTemplateFile, "html-template-path", "", "path of the HTML template file to use")
	rootCmd.Flags().StringVar(&htmlTitle, "html-title", "", "title to use in the HTML output")
	rootCmd.Flags().BoolVarP(&openFailed, "open-failed", "o", false, "whether to open failed workflows")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Glance at the last 3 runs of your GitHub Actions workflows")
		flag.PrintDefaults()
	}

	return rootCmd, nil
}

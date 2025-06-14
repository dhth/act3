package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	errCouldntGetUserHomeDir      = errors.New("couldn't get your home directory")
	errCouldntGetUserConfigDir    = errors.New("couldn't get your config directory")
	errCouldntReadConfigFile      = errors.New("couldn't read config file")
	errFlagCombIncorrect          = errors.New("invalid flag combination")
	errIncorrectOutputFmt         = errors.New("incorrect value for output format provided")
	ErrConfigFileDoesntExit       = errors.New("config file doesn't exist")
	ErrCouldntGetConfig           = errors.New("couldn't get config")
	errCouldntGetGHClient         = errors.New("couldn't get a Github client")
	errNoWorkflows                = errors.New("no workflows found")
	errTemplateFileDoesntExit     = errors.New("template file doesn't exist")
	errCouldntReadTemplateFile    = errors.New("couldn't read template file")
	errCouldntGetWorkflows        = errors.New("couldn't get workflows")
	ErrCouldntMarshalConfigToYAML = errors.New("couldn't marshal workflows to YAML")
	errInvalidWorkflowFilterRegex = errors.New("workflow filter is invalid regex")
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
		configPath        string
		configBytes       []byte
		homeDir           string
		globalWorkflows   bool
		repos             []string
		workflowFilterStr string
		formatStr         string
		htmlTemplateFile  string
		htmlTitle         string
		openFailed        bool
	)

	rootCmd := &cobra.Command{
		Use:          "act3",
		Short:        "Glance at the last 3 runs of your GitHub Actions workflows",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			if globalWorkflows && len(repos) > 0 {
				return fmt.Errorf("%w; -g and -r cannot both be provided at the same time", errFlagCombIncorrect)
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
				workflows, cfgErr = getWorkflowsFromConfig(configBytes)
				if cfgErr != nil {
					return fmt.Errorf("%w: %s", ErrCouldntGetConfig, cfgErr.Error())
				}

			} else {
				if len(repos) > 0 {
					err := validateRepos(repos)
					if err != nil {
						return err
					}
					reposToUse = repos
				} else {
					currentRepo, err := getCurrentRepo()
					if err != nil {
						return err
					}
					reposToUse = []string{currentRepo}
				}

				ghRestClient, err := ghapi.NewRESTClient(clientOpts)
				if err != nil {
					return fmt.Errorf("%w: %s", errCouldntGetGHClient, err.Error())
				}

				var workflowFilter *regexp.Regexp
				if workflowFilterStr != "" {
					r, err := regexp.Compile(workflowFilterStr)
					if err != nil {
						return fmt.Errorf("%w: %s", errInvalidWorkflowFilterRegex, err.Error())
					}
					workflowFilter = r
				}

				var errors []WorkflowError
				workflows, errors = getWorkflowsForRepos(ghRestClient, reposToUse, workflowFilter)

				if len(errors) == 1 {
					return fmt.Errorf("%w: %s", errCouldntGetWorkflows, errors[0].Err.Error())
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

			ghGQLClient, err := ghapi.NewGraphQLClient(clientOpts)
			if err != nil {
				return fmt.Errorf("%w: %s", errCouldntGetGHClient, err.Error())
			}

			var cr *string
			if !globalWorkflows && len(reposToUse) == 1 {
				cr = &reposToUse[0]
			}
			config := types.RunConfig{
				GHClient:     ghGQLClient,
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

	configCmd := &cobra.Command{
		Use:          "config",
		Short:        "Interact with act3's config",
		SilenceUsage: true,
	}

	generateConfigCmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate act3's config",
		Long: `Generate act3's config.
You can either generate the config for the current repository or for the list of repos provided by you using the -r/--repos flag.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			var reposToUse []string
			if len(repos) > 0 {
				err := validateRepos(repos)
				if err != nil {
					return err
				}
				reposToUse = repos
			} else {
				currentRepo, err := getCurrentRepo()
				if err != nil {
					return err
				}
				reposToUse = []string{currentRepo}
			}

			var workflowFilter *regexp.Regexp
			if workflowFilterStr != "" {
				r, err := regexp.Compile(workflowFilterStr)
				if err != nil {
					return fmt.Errorf("%w: %s", errInvalidWorkflowFilterRegex, err.Error())
				}
				workflowFilter = r
			}

			clientOpts := ghapi.ClientOptions{
				EnableCache: false,
				Timeout:     8 * time.Second,
			}

			ghRestClient, err := ghapi.NewRESTClient(clientOpts)
			if err != nil {
				return fmt.Errorf("%w: %s", errCouldntGetGHClient, err.Error())
			}

			workflows, errors := getWorkflowsForRepos(ghRestClient, reposToUse, workflowFilter)

			if len(errors) == 1 {
				return fmt.Errorf("%w: %s", errCouldntGetWorkflows, errors[0].Err.Error())
			}

			if len(errors) > 1 {
				errorStrs := make([]string, len(errors))
				for i, e := range errors {
					errorStrs[i] = fmt.Sprintf("- %s: %s", e.Repo, e.Err.Error())
				}

				return fmt.Errorf("%w:\n%s", errCouldntGetWorkflows, strings.Join(errorStrs, "\n"))
			}

			if len(workflows) == 0 {
				fmt.Fprintln(os.Stderr, "no workflows found")
				return nil
			}

			config := config{
				Workflows: workflows,
			}

			configBytes, err := config.MarshalToYAML()
			if err != nil {
				return fmt.Errorf("%w: %s", ErrCouldntMarshalConfigToYAML, err.Error())
			}

			fmt.Printf("%s", configBytes)

			return nil
		},
	}

	validateConfigCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate act3's config",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			configPathFull := utils.ExpandTilde(configPath, homeDir)

			var cfgErr error
			_, cfgErr = os.Stat(configPathFull)
			if os.IsNotExist(cfgErr) {
				return fmt.Errorf("%w, path: %q", ErrConfigFileDoesntExit, configPathFull)
			}

			configPathFull = utils.ExpandTilde(configPath, homeDir)
			configBytes, cfgErr = os.ReadFile(configPathFull)
			if cfgErr != nil {
				return fmt.Errorf("%w: %w", errCouldntReadConfigFile, cfgErr)
			}
			_, cfgErr = getWorkflowsFromConfig(configBytes)
			if cfgErr != nil {
				return cfgErr
			}

			fmt.Fprintln(os.Stderr, "config looks good ✅")

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

	rootCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, "location of act3's config file")
	rootCmd.Flags().StringSliceVarP(&repos, "repos", "r", []string{}, `repos to fetch workflows for, in the format "owner/repo"`)
	rootCmd.Flags().StringVarP(&workflowFilterStr, "workflow-name-filter", "n", "", "regex expression to filter workflows by name")
	rootCmd.Flags().BoolVarP(&globalWorkflows, "global", "g", false, "whether to use workflows defined globally via the config file")
	rootCmd.Flags().StringVarP(&formatStr, "format", "f", "default", "output format to use; possible values: default, table, html")
	rootCmd.Flags().StringVar(&htmlTemplateFile, "html-template-path", "", "path of the HTML template file to use")
	rootCmd.Flags().StringVar(&htmlTitle, "html-title", "act3", "title to use in the HTML output")
	rootCmd.Flags().BoolVarP(&openFailed, "open-failed", "o", false, "whether to open failed workflows")

	generateConfigCmd.Flags().StringSliceVarP(&repos, "repos", "r", []string{}, `repos to generate the config for, in the format "owner/repo"`)
	generateConfigCmd.Flags().StringVarP(&workflowFilterStr, "workflow-name-filter", "n", "", "regex expression to filter workflows by name")

	validateConfigCmd.Flags().StringVarP(&configPath, "config-path", "c", defaultConfigPath, "location of act3's config file")

	configCmd.AddCommand(generateConfigCmd)
	configCmd.AddCommand(validateConfigCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd, nil
}

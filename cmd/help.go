package cmd

import "fmt"

var configSampleFormat = `
workflows:
- id: W_kwDOLafHJ84FQglU
  repo: dhth/outtasync
  name: release
  key: outtasync:release
- id: W_kwDOLghtl84FWTla
  repo: dhth/ecsv
  name: release
  key: ecsv:release
- id: W_kwDOLb3Pms4FRxjY
  repo: dhth/cueitup
  name: release
  key: cueitup:release
`

func getHelp(configFilePath string) string {
	return fmt.Sprintf(`Glance at the last 3 runs of your Github Actions.

Usage:
  act3 [flags]

Flags:
  -c string                       path of the config file (default "%s")
  -f string                       output format to use; possible values: default, table, html (default "default")
  -t string                       path of the HTML template file to use
  -r string                       repo to fetch workflows for, in the format "owner/repo"
  -g bool                         whether to use workflows defined globally via the config file (default false)
  -o bool                         whether to open failed workflows (via your OS's "open" command) (default false)
  -h, --help                      help for act3
`, configFilePath)
}

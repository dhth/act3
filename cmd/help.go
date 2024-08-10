package cmd

var (
	configSampleFormat = `
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
	helpText = `Glance at the last 3 runs of your Github Actions.

Usage: act3 [flags]`
)

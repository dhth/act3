package ui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/shurcooL/githubv4"
)

func RenderUI(ghClient *githubv4.Client, workflows []Workflow, outputFmt OutputFmt, htmlTemplate string) {

	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	var opts []tea.ProgramOption
	if outputFmt != UnspecifiedFmt {
		opts = append(opts, tea.WithoutRenderer())
		// TODO: this may be a hack, and will prevent using STDIN for
		// CLI mode, find a better way
		opts = append(opts, tea.WithInput(nil))
	}
	p := tea.NewProgram(InitialModel(ghClient, workflows, outputFmt, htmlTemplate), opts...)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Something went wrong %s", err)
	}
}

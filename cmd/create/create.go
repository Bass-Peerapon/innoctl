/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"fmt"

	"github.com/Bass-Peerapon/innoctl/cmd/create/crud"
	"github.com/Bass-Peerapon/innoctl/cmd/create/dbmlstruct"
	"github.com/Bass-Peerapon/innoctl/cmd/create/nwp"
	"github.com/Bass-Peerapon/innoctl/cmd/create/project"
	"github.com/Bass-Peerapon/innoctl/cmd/create/service"
	filepicker "github.com/Bass-Peerapon/innoctl/ui/file-picker"
	multiselect "github.com/Bass-Peerapon/innoctl/ui/multiSelect"
	progressbar "github.com/Bass-Peerapon/innoctl/ui/progress-bar"
	"github.com/Bass-Peerapon/innoctl/ui/selection"
	"github.com/Bass-Peerapon/innoctl/ui/textinput"
	"github.com/Bass-Peerapon/innoctl/ui/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const logo = `
  ___ _   _ _   _  _____     ___    ____ _____     _______ 
 |_ _| \ | | \ | |/ _ \ \   / / \  / ___|_ _\ \   / / ____|
  | ||  \| |  \| | | | \ \ / / _ \ \___ \| | \ \ / /|  _|  
  | || |\  | |\  | |_| |\ V / ___ \ ___) | |  \ V / | |___ 
 |___|_| \_|_| \_|\___/  \_/_/   \_\____/___|  \_/  |_____|
                                                           
`

var (
	logoStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("170")).
				Bold(true)
)

func CreateService() {
	serviceName := textinput.InitialTextInputModel("What is the name of your service?")
	p := tea.NewProgram(serviceName)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	if serviceName.IsExit() {
		return
	}
	drivers := []selection.Choice{
		{Title: "Postgres", Desc: "Go postgres driver for Go"},
		{Title: "MongoDB", Desc: "The MongoDB supported driver for Go."},
	}
	serviceModel := selection.NewModel("Select Repository Client", drivers)
	p = tea.NewProgram(serviceModel)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	if serviceModel.IsExit() {
		return
	}
	progressbarProgram := tea.NewProgram(progressbar.InitialProgressbarModel(7))
	c := make(chan struct{})
	go func() {
		go service.GenerateWithChannel(
			serviceName.GetOutput(),
			drivers[serviceModel.GetIndex()].Title,
			c,
		)
		for i := range c {
			progressbarProgram.Send(progressbar.ProgressMsg(i))
		}
		// end progress bar
		progressbarProgram.Send(progressbar.ProgressMsg{})
	}()
	if _, err := progressbarProgram.Run(); err != nil {
		panic(err)
	}
}

func CreateNewWithParams() {
	filepickerModel := filepicker.InitialFilepickerModel()
	p := tea.NewProgram(filepickerModel)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	if filepickerModel.IsExit() {
		return
	}

	fmt.Printf(
		"Selected file: " + selectedItemStyle.Render(filepickerModel.GetSelectedFile()) + "\n\n",
	)
	progressbarProgram := tea.NewProgram(progressbar.InitialProgressbarModel(5))
	c := make(chan struct{})
	go func() {
		go nwp.GeneratorWithChannel(filepickerModel.GetSelectedFile(), c)
		for i := range c {
			progressbarProgram.Send(progressbar.ProgressMsg(i))
		}
		// end progress bar
		progressbarProgram.Send(progressbar.ProgressMsg{})
	}()
	if _, err := progressbarProgram.Run(); err != nil {
		panic(err)
	}
}

func CreateCrud() {
	filepickerModel := filepicker.InitialFilepickerModel()
	filepickerModel.SetTitle("Model")
	p := tea.NewProgram(filepickerModel)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	if filepickerModel.IsExit() {
		return
	}
	fmt.Printf(
		"Selected model file: " + selectedItemStyle.Render(
			filepickerModel.GetSelectedFile(),
		) + "\n\n",
	)

	filepickerRepo := filepicker.InitialFilepickerModel()
	filepickerRepo.SetTitle("Repo")
	p = tea.NewProgram(filepickerRepo)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	if filepickerRepo.IsExit() {
		return
	}
	fmt.Printf(
		"Selected repo	file: " + selectedItemStyle.Render(
			filepickerModel.GetSelectedFile(),
		) + "\n\n",
	)

	drivers := []selection.Choice{
		{Title: "Postgres", Desc: "Go postgres driver for Go"},
		{Title: "MongoDB", Desc: "The MongoDB supported driver for Go."},
	}
	serviceModel := selection.NewModel("Select Repository Client", drivers)
	p = tea.NewProgram(serviceModel)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	if serviceModel.IsExit() {
		return
	}

	choices := []multiselect.Choice{
		{Title: "Create", Desc: "Generate Create Funtion"},
		{Title: "Read", Desc: "Generate Read Funtion"},
		{Title: "Update", Desc: "Generate Update Funtion"},
		{Title: "Delete", Desc: "Generate Delete Funtion"},
	}

	multiselectOpts := multiselect.NewModel("Select options", choices)
	p = tea.NewProgram(multiselectOpts)
	if _, err := p.Run(); err != nil {
		panic(err)
	}
	if multiselectOpts.IsExit() {
		return
	}

	progressbarProgram := tea.NewProgram(
		progressbar.InitialProgressbarModel(len(multiselectOpts.GetIndexs())),
	)
	c := make(chan struct{})
	go func() {
		go func() {
			switch drivers[serviceModel.GetIndex()].Title {
			case "Postgres":
				for _, v := range multiselectOpts.GetIndexs() {
					switch choices[v].Title {
					case "Create":
						if err := crud.GenCreatePostgres(filepickerModel.GetSelectedFile(), filepickerRepo.GetSelectedFile()); err != nil {
							fmt.Println(err)
						}
						c <- struct{}{}
					case "Read":
						if err := crud.GenReadPostgres(filepickerModel.GetSelectedFile(), filepickerRepo.GetSelectedFile()); err != nil {
							fmt.Println(err)
						}
						c <- struct{}{}
					case "Update":
						if err := crud.GenUpdatePostgres(filepickerModel.GetSelectedFile(), filepickerRepo.GetSelectedFile()); err != nil {
							fmt.Println(err)
						}
						c <- struct{}{}
					case "Delete":
						if err := crud.GenDeletePostgres(filepickerModel.GetSelectedFile(), filepickerRepo.GetSelectedFile()); err != nil {
							fmt.Println(err)
						}
						c <- struct{}{}
					}
				}
				close(c)
			case "MongoDB":
			}
		}()
		for i := range c {
			progressbarProgram.Send(progressbar.ProgressMsg(i))
		}
		// end progress bar
		progressbarProgram.Send(progressbar.ProgressMsg{})
	}()
	if _, err := progressbarProgram.Run(); err != nil {
		panic(err)
	}
}

func CreateDbml2struct() {
	dbmlView := viewport.InitModel()

	p := tea.NewProgram(dbmlView, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		panic(err)
	}

	fmt.Println(*dbmlView.Output)
}

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "innovasive control tool",
	Long:  `innovasive control tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", logoStyle.Render(logo))
		tools := []selection.Choice{
			{Title: "Service", Desc: "Generate New Service"},
			{Title: "NewWithParams", Desc: "Generate New With Params Function"},
			{Title: "CRUD", Desc: "Generate CRUD"},
			{Title: "Dbml2struct", Desc: "Generate Struct from Dbml"},
		}
		tool := selection.NewModel("Select Tool", tools)
		p := tea.NewProgram(tool)
		if _, err := p.Run(); err != nil {
			panic(err)
		}
		if tool.IsExit() {
			return
		}
		switch tools[tool.GetIndex()] {
		case tools[0]:
			CreateService()
		case tools[1]:
			CreateNewWithParams()
		case tools[2]:
			CreateCrud()
		case tools[3]:
			CreateDbml2struct()
		default:
			fmt.Println("invalid choice")
		}
	},
}

func init() {
	CreateCmd.AddCommand(service.ServiceCmd)
	CreateCmd.AddCommand(nwp.NewWithParamsCmd)
	CreateCmd.AddCommand(crud.CrudRepoCmd)
	CreateCmd.AddCommand(dbmlstruct.DbmlstructCmd)
	CreateCmd.AddCommand(project.ProjectCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package qdb

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
	"os"
)

func ManageConnections() error {
	for {
		fmt.Println("Current connections: ")
		if err := ListConnections(); err != nil {
			return nil
		}
		selectPrompt := promptui.Select{
			Label: "What do you want to do?",
			Items: []string{"Add a new connection", "Edit a connection", "Delete a connection", "Exit"},
		}

		result, _, err := selectPrompt.Run()
		if err != nil {
			return err
		}
		if result == 0 {
			namePrompt := promptui.Prompt{
				Label: "Connection Name",
			}
			name, err := namePrompt.Run()
			if err != nil {
				return err
			}
			urlPrompt := promptui.Prompt{
				Label: "URL",
			}
			url, err := urlPrompt.Run()
			if err != nil {
				return err
			}
			if err := AddConnection(name, url); err != nil {
				return err
			}
		} else if result == 1 {
			panic("edit not implemented")
		} else if result == 2 {
			panic("delete not implemented")
		} else if result == 3 {
			return nil
		}
	}
}

func ListConnections() error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := table.Row{"Name", "URL", "Default"}
	rows := make([]table.Row, 0)
	for _, p := range Profiles {
		rows = append(rows, table.Row{p.Name, p.Url, IsDefaultProfile(p.Name)})
	}
	t.AppendHeader(header)
	t.AppendRows(rows)

	t.AppendSeparator()
	t.Render()
	return nil
}

func AddConnection(name string, url string) error {
	return AddProfile(name, url)
}

func DeleteConnection(name string) error {
	return DeleteProfile(name)
}

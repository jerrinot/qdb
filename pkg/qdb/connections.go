package qdb

import (
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
	"net/url"
	"os"
	"qdb/pkg/qdb/config"
	"strings"
)

func ManageConnections() error {
	for {
		selectPrompt := promptui.Select{
			Label: "What do you want to do?",
		}
		if len(config.ConnectionDefs) == 0 {
			fmt.Println("No connection exists")
			selectPrompt.Items = []string{"Add a new connection", "Exit"}
		} else {
			fmt.Println("Current connections: ")
			selectPrompt.Items = []string{"Add a new connection", "Edit a connection", "Delete a connection", "Exit"}
			if err := ListConnections(); err != nil {
				return nil
			}
		}

		_, result, err := selectPrompt.Run()
		if err != nil {
			return err
		}
		if result == "Add a new connection" {
			if err := CreateNewConnection(); err != nil {
				return nil
			}
		} else if result == "Edit a connection" {
			panic("edit not implemented")
		} else if result == "Delete a connection" {
			if err := DeleteConnection(); err != nil {
				return nil
			}
		} else if result == "Exit" {
			return nil
		}
	}
}

func testConnection(serverUrl string) error {
	res, err := callGet(AddQueryPath(serverUrl) + url.QueryEscape("select now();"))
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if res.StatusCode == 200 {
		return nil
	}
	return errors.New(res.Status)
}

func CreateNewConnection() error {
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
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	if err := testConnection(url); err != nil {
		fmt.Printf("Cannot connect to the server %s\n", err)
		prompt := promptui.Prompt{
			Label:     "Do you really want to add it",
			IsConfirm: true,
		}

		_, err = prompt.Run()
		errOrig := err
		if err == nil {
			return config.AddConnection(name, url)
		}

		prompt = promptui.Prompt{
			Label:     "Do want to add a different connection",
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err != nil {
			return errOrig
		}
		return CreateNewConnection()
	}
	return config.AddConnection(name, url)
}

func ChooseConnection(canSetAsDefault bool) (config.ConnectionDef, error) {
	conns := make([]string, len(config.ConnectionDefs))
	for i, p := range config.ConnectionDefs {
		conns[i] = p.Name
	}
	selectPrompt := promptui.Select{
		Label: "Select Connection",
		Items: conns,
	}
	i, _, err := selectPrompt.Run()
	if err != nil {
		return config.ConnectionDef{}, err
	}

	if canSetAsDefault && config.DefaultConnectionName == "" {
		prompt := promptui.Prompt{
			Label:     "Do you want to set this connection as default",
			IsConfirm: true,
		}

		_, err = prompt.Run()
		if err != nil {
			return config.ConnectionDef{}, err
		}
		config.SetAsDefaultConnection(config.ConnectionDefs[i].Name)
	}
	return config.ConnectionDefs[i], err
}

func DeleteConnection() error {
	fmt.Println("Choose a connection to delete")
	connection, err := ChooseConnection(false)
	if err != nil {
		return err
	}
	return config.DeleteConnection(connection.Name)
}

func ListConnections() error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	header := table.Row{"Name", "URL", "Default"}
	rows := make([]table.Row, 0)
	for _, p := range config.ConnectionDefs {
		rows = append(rows, table.Row{p.Name, p.Url, config.IsDefaultConnection(p.Name)})
	}
	t.AppendHeader(header)
	t.AppendRows(rows)

	t.AppendSeparator()
	t.Render()
	return nil
}

package qdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"qdb/pkg/qdb/config"
	"strings"
)

type CsrfResponse struct {
	CsrfToken string `json:"csrfToken"`
}

type CloudInstance struct {
	Id                string `json:"id"`
	Name              string `json:"db_name"`
	Region            string `json:"region"`
	InstanceType      string `json:"instance_type"`
	HttpBasicUser     string `json:"http_basic_auth_user"`
	HttpBasicPassword string `json:"http_basic_auth_password"`
	RestEndpoint      string `json:"rest_endpoint"`
}

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
				return err
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

func CreateNewCloudConnection() error {
	resp, err := callGet("https://cloud.app.questdb.net/auth/signin?callbackUrl=https%3A%2F%2Fcloud.app.questdb.net")
	if err != nil {
		return nil
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	resp, err = callGet("https://cloud.app.questdb.net/api/auth/csrf")
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var csrfResponse CsrfResponse
	if err := json.Unmarshal(body, &csrfResponse); err != nil {
		return err
	}

	namePrompt := promptui.Prompt{
		Label: "email",
	}
	email, err := namePrompt.Run()
	if err != nil {
		return err
	}
	passwordPrompt := promptui.Prompt{
		Label: "password",
		Mask:  '*',
	}
	password, err := passwordPrompt.Run()
	if err != nil {
		return err
	}

	v := url.Values{}
	v.Add("email", email)
	v.Add("password", password)
	v.Add("csrfToken", csrfResponse.CsrfToken)
	v.Add("json", "true")

	req, err := http.NewRequest("POST", "https://cloud.app.questdb.net/api/auth/callback/credentials-email-login", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode == 200 {
		resp, err = callGet("https://api.app.questdb.net/client/v1/instances")
		if err != nil {
			return nil
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var instances []CloudInstance
		if err := json.Unmarshal(body, &instances); err != nil {
			return err
		}
		for _, i := range instances {
			fmt.Println(i)
		}
		return nil
	} else if resp.StatusCode == 401 {
		return errors.New("wrong QuestDB Cloud credentials")
	} else {
		return errors.New("unknown error " + resp.Status)
	}
}

func CreateNewConnection() error {
	selectPrompt := promptui.Select{
		Label: "Select Connection",
		Items: []string{"Localhost", "Cloud", "Demo", "Custom"},
	}
	resp, _, err := selectPrompt.Run()
	if err != nil {
		return err
	}
	if resp == 0 {
		panic("not implemented")
	}
	if resp == 1 {
		return CreateNewCloudConnection()
	}
	if resp == 2 {
		panic("not implemented")
	}
	if resp == 3 {
		return CreateNewHttpConnection()
	}
	panic("unexpected choice. should not happen")
}

func CreateNewHttpConnection() error {
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
		return CreateNewHttpConnection()
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

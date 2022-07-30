package qdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/knz/go-libedit"
	"github.com/manifoldco/promptui"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"qdb/pkg/qdb/config"
	"strings"
	"time"
)

const localhostPrefix = "http://localhost:9000/exec?count=true&query="
const demoPrefix = "https://demo.questdb.io/exec?count=true&query="

//const QdbServerPrefix = localhostPrefix

type ErrorResponse struct {
	Query   string `json:"query"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type ResultSet struct {
	Query   string   `json:"query"`
	Column  []Column `json:"columns"`
	Dataset [][]any  `json:"dataset"`
}

type Column struct {
	Name       string `json:"name"`
	ColumnType string `json:"type"`
}

var httpClient = http.Client{
	Timeout: time.Second * 1000,
}

type example struct{}

func (_ example) GetCompletions(word string) []string {
	if strings.HasPrefix(word, "sele") {
		return []string{"select "}
	}
	if strings.HasPrefix(word, "fro") {
		return []string{"from "}
	}
	if strings.HasPrefix(word, "wher") {
		return []string{"where "}
	}
	return nil
}

func resolveConnectionName(connectionName string) (string, error) {
	if connectionName == "" {
		connectionName = config.DefaultConnectionName
	} else if !config.ConnectionExists(connectionName) {
		return "", errors.New("connection name '" + connectionName + "' does not exist")
	}
	if connectionName == "" {
		if len(config.ConnectionDefs) == 0 {
			fmt.Println("No connection to QuestDB server found")
			prompt := promptui.Prompt{
				Label:     "Add a new connection",
				IsConfirm: true,
			}

			_, err := prompt.Run()
			if err != nil {
				return "", errors.New("no connection exists")
			} else {
				if err := CreateNewConnection(); err != nil {
					return "", err
				}
			}
		}

		if len(config.ConnectionDefs) == 1 {
			fmt.Println("Choosing the only existing connection: " + config.ConnectionDefs[0].Name)
			connectionName = config.ConnectionDefs[0].Name
		} else {
			c, err := ChooseConnection(true)
			if err != nil {
				return "", err
			}
			connectionName = c.Name
		}
	}
	return connectionName, nil
}

func AddQueryPath(baseUrl string) string {
	if !strings.HasSuffix(baseUrl, "/") {
		baseUrl = baseUrl + "/"
	}
	return baseUrl + "exec?count=true&query="
}

func toServerPrompt(conn config.ConnectionDef) string {
	return conn.Name + "> "
}

func RunSqlShell(query string, connectionName string) error {
	connectionName, err := resolveConnectionName(connectionName)
	if err != nil {
		return err
	}
	connection, err := config.ConnectionByName(connectionName)
	if err != nil {
		return err
	}

	if query != "" {
		return runAndPrintQuery(connection.Url, query)
	}

	// Open and immediately close a libedit instance to test that nonzero editor
	// IDs are tracked correctly.
	el, err := libedit.Init("example", true)
	if err != nil {
		return err
	}
	el.Close()

	el, err = libedit.Init("example", true)
	if err != nil {
		return err
	}
	defer el.Close()

	// RebindControlKeys ensures that Ctrl+C, Ctrl+Z, Ctrl+R and Tab are
	// properly bound even if the user's .editrc has used bind -e or
	// bind -v to load a predefined keymap.
	el.RebindControlKeys()

	el.UseHistory(-1, true)
	el.LoadHistory("hist")
	el.SetAutoSaveHistory("hist", true)
	el.SetCompleter(example{})
	el.SetLeftPrompt(toServerPrompt(connection))
	//el.SetRightPrompt("(-:")
	buff := ""
	for {
		s, err := el.GetLine()
		buff += s
		if err != nil {

			if err == io.EOF {
				break
			}
			if err == libedit.ErrInterrupted {
				break
			}
			log.Fatal(err)
		}
		// todo: deal with escaping and whitespaces after ;
		if strings.HasSuffix(s, ";\n") {
			if err := el.AddHistory(buff); err != nil {
				return err
			}
			err := runAndPrintQuery(connection.Url, buff)
			if err != nil {
				fmt.Printf("Error while running query %e\n", err)
			}
			el.SetLeftPrompt(toServerPrompt(connection))
			buff = ""
		} else {
			el.SetLeftPrompt(" > ")
		}
	}
	return nil
}

func callGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return httpClient.Do(req)
}

func runAndPrintQuery(serverUrl string, query string) error {
	qurl := AddQueryPath(serverUrl) + url.QueryEscape(query)
	res, err := callGet(qurl)
	if err != nil {
		return err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	respString := string(body)
	//fmt.Println(respString)
	if res.StatusCode == http.StatusOK {
		var rs ResultSet
		if err := json.Unmarshal(body, &rs); err != nil {
			fmt.Println("error while unmarshalling JSON: " + respString)
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		header := table.Row{}
		for _, arrColumn := range rs.Column {
			header = append(header, arrColumn.Name)
		}
		rows := make([]table.Row, 0)
		for _, arrRow := range rs.Dataset {
			rows = append(rows, arrRow)
		}
		t.AppendHeader(header)
		t.AppendRows(rows)

		t.AppendSeparator()
		t.Render()
	} else {
		var errResponse ErrorResponse
		if err := json.Unmarshal(body, &errResponse); err != nil {
			fmt.Println("error while unmarshalling JSON: " + respString)
		}
		if errResponse.Error != "" {
			fmt.Println(errResponse.Error)
		} else if errResponse.Message != "" {
			fmt.Println(errResponse.Message)
		} else {
			fmt.Println("Unexpected error JSON: " + respString)
		}
	}
	return nil
}

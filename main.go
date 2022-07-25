package main

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/knz/go-libedit"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type ResultSet struct {
	Query   string   `json:"query"`
	Column  []Column `json:"columns"`
	Error   string   `json:"error"`
	Dataset [][]any  `json:"dataset"`
}

type Column struct {
	Name       string `json:"name"`
	ColumnType string `json:"type"`
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

func main() {

	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	// Open and immediately close a libedit instance to test that nonzero editor
	// IDs are tracked correctly.
	el, err := libedit.Init("example", true)
	if err != nil {
		log.Fatal(err)
	}
	el.Close()

	el, err = libedit.Init("example", true)
	if err != nil {
		log.Fatal(err)
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
	el.SetLeftPrompt("qdb> ")
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
		if strings.HasSuffix(s, ";\n") {
			if err := el.AddHistory(s); err != nil {
				log.Fatal(err)
			}
			if err := el.AddHistory(buff); err != nil {
				log.Fatal(err)
			}
			//fmt.Println(buff)
			qurl := "http://localhost:9000/exec?count=true&query=" + url.QueryEscape(buff)
			req, err := http.NewRequest(http.MethodGet, qurl, nil)
			if err != nil {
				log.Fatal(err)
			}
			res, getErr := httpClient.Do(req)
			if getErr != nil {
				log.Fatal(getErr)
			}

			if res.Body != nil {
				defer res.Body.Close()
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			//respString := string(body)
			//fmt.Println(respString)
			var rs ResultSet
			if err := json.Unmarshal(body, &rs); err != nil {
				log.Fatal(err)
			}
			el.SetLeftPrompt("qdb> ")
			buff = ""
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
			el.SetLeftPrompt(" > ")
		}
	}
	fmt.Println("goodbye!")
}

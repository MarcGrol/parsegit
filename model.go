package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type Author struct {
	Date      int64     `json:"date,omitempty"` // dropped after parse
	Timestamp time.Time `json:"timestamp"`      // enriched based on Date
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone"`
}

type Committer struct {
	Date      int64     `json:"date,omitempty"` // dropped after parse
	Timestamp time.Time `json:"timestamp"`      // enriched based on Date
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Timezone  string    `json:"timezone"`
}

type Change struct {
	LinesAdded   int    `json:"linesAdded"`
	LinesRemoved int    `json:"linesRemoved"`
	Filename     string `json:"filename"`
}

type ChangeSet struct {
	IsMerge         bool     `json:"isMerge"`
	NumFilesChanged int      `json:"numFilesChanged"`
	LinesAdded      int      `json:"linesAdded"`
	LinesRemoved    int      `json:"linesRemoved"`
	Changes         []Change `json:"changes"`
}

type Commit struct {
	Author    Author          `json:"author"`
	Changes   [][]interface{} `json:"changes,omitempty"` // dropped after parse
	ChangeSet ChangeSet       `json:"changeSet"`         // enriched based on Changes
	Commit    string          `json:"commit"`
	Committer Committer       `json:"committer"`
	Message   string          `json:"message"`
	Parents   []string        `json:"parents"`
	Tree      string          `json:"tree"`
}

type Commits []Commit

func (c Commits) Len() int      { return len(c) }
func (c Commits) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

type ByTimestamp struct{ Commits }

func (c ByTimestamp) Less(i, j int) bool {
	return c.Commits[i].Author.Timestamp.Before(c.Commits[j].Author.Timestamp)
}

func parse(filename string, callback func(c Commit)) error {
	// parse file
	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading file %s: %s", filename, err)
	}

	// parse json
	commits := []Commit{}
	err = json.Unmarshal(jsonBlob, &commits)
	if err != nil {
		return fmt.Errorf("Error parsing file %s: %s", filename, err)
	}

	for _, c := range commits {
		// parse unix timestamp into time
		c.Committer.Timestamp = time.Unix(c.Committer.Date, 0)
		c.Committer.Date = 0 // remove old raw value

		// parse unix timestamp into time
		c.Author.Timestamp = time.Unix(c.Author.Date, 0)
		c.Author.Date = 0 // remove old raw value

		// parse array of heterogenous array into something meaningfull
		changes := []Change{}
		numAdds := 0
		numDels := 0
		for _, s := range c.Changes {
			change := Change{}
			for i, p := range s {
				switch i {
				case 0:
					val, ok := p.(float64)
					if ok {
						change.LinesAdded = int(val)
						numAdds += change.LinesAdded
					}
				case 1:
					val, ok := p.(float64)
					if ok {
						change.LinesRemoved = int(val)
						numDels += change.LinesRemoved
					}
				case 2:
					change.Filename, _ = p.(string)
				}
			}
			changes = append(changes, change)
		}
		c.ChangeSet = ChangeSet{
			Changes:         changes,
			NumFilesChanged: len(changes),
			LinesAdded:      numAdds,
			LinesRemoved:    numDels,
		}

		c.Changes = [][]interface{}{} // remove old raw value

		// mark as merge
		if strings.HasPrefix(c.Message, "Merge branch") {
			c.ChangeSet.IsMerge = true
		}

		// call logic
		callback(c)
	}

	return nil
}

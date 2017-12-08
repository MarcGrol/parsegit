package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

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

		// parse array of array with heterogenous elements into something typestrong
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
		if strings.HasPrefix(c.Message, "Merge") {
			c.ChangeSet.IsMerge = true
		}

		// call logic
		callback(c)
	}

	return nil
}

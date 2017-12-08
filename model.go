package main

import (
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

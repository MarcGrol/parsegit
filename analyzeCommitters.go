package main

import (
	"fmt"
	"os"
)

type committerSummary struct {
	NumCommits   int
	LinesAdded   int
	LinesRemoved int
	FilesChanged int
}

func analyseCommitters(commits Commits) {
	// analyze work done by commiter
	fmt.Fprintf(os.Stdout, "committer-name; #commits; #additions; #removals; #files\n")
	committers := map[string]committerSummary{}
	for _, c := range commits {
		commiterInfo(committers, c)
	}
	for name, summ := range committers {
		fmt.Fprintf(os.Stdout, "%s; %d; %d; %d; %d\n",
			name, summ.NumCommits, summ.LinesAdded, summ.LinesRemoved, summ.FilesChanged)
	}
}

func commiterInfo(committers map[string]committerSummary, commit Commit) {
	name := parseName(commit.Author.Name)
	sum, found := committers[name]
	if !found {
		sum = committerSummary{}
	}
	sum.LinesAdded += commit.ChangeSet.LinesAdded
	sum.LinesRemoved += commit.ChangeSet.LinesRemoved
	sum.FilesChanged += commit.ChangeSet.NumFilesChanged
	sum.NumCommits++
	committers[name] = sum
}

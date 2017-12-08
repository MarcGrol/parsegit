package main

import (
	"fmt"
	"os"
	"strings"
)

type committerSummary struct {
	NumCommits   int
	LinesAdded   int
	LinesRemoved int
	FilesChanged int
}

func doAnalyseCommitters(commits Commits) {
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

type fileSummary struct {
	NumCommits   int
	LinesAdded   int
	LinesRemoved int
	Commiters    map[string]int
}

func doAnalyseFiles(commits Commits) {
	// analyze work done on files
	fmt.Fprintf(os.Stdout, "filename; #commits; #additions; #removals; #commiters\n")
	files := map[string]fileSummary{}
	for _, c := range commits {
		fileInfo(files, c)
	}
	// print as csv
	for fn, summ := range files {
		fmt.Fprintf(os.Stdout, "%s;%d;%d;%d;%d\n",
			fn, summ.NumCommits, summ.LinesAdded, summ.LinesRemoved, len(summ.Commiters))
	}
}

func fileInfo(files map[string]fileSummary, commit Commit) {
	committerName := parseName(commit.Author.Name)
	for _, cs := range commit.ChangeSet.Changes {
		filename := cs.Filename
		sum, found := files[filename]
		if !found {
			sum = fileSummary{
				Commiters: map[string]int{},
			}
		}
		sum.LinesAdded += commit.ChangeSet.LinesAdded
		sum.LinesRemoved += commit.ChangeSet.LinesRemoved
		sum.NumCommits++
		cm := sum.Commiters[cs.Filename]
		cm++
		sum.Commiters[committerName] = cm

		files[filename] = sum
	}
}

func parseName(in string) string {
	parts := strings.Split(in, " ")
	return strings.ToLower(parts[0])
}

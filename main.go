package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type committerSummary struct {
	NumCommits   int
	LinesAdded   int
	LinesRemoved int
	FilesChanged int
}

type fileSummary struct {
	NumCommits   int
	LinesAdded   int
	LinesRemoved int
	Commiters    map[string]int
}

func main() {
	// read config
	filename, analyzeCommiters, analyzeFiles, dumpAsJson := parseArgs()

	// parse and post-process git2json file
	commits := Commits{}
	parse(filename, func(commit Commit) {
		commits = append(commits, commit)
	})
	sort.Sort(ByTimestamp{commits})

	if dumpAsJson {
		// dump to debug
		for _, c := range commits {
			logCommit(c)
		}
	}

	if analyzeCommiters {
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

	if analyzeFiles {
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
}

func logCommit(commit Commit) {
	fmt.Fprintf(os.Stdout, "%s - %s - #files: %d, #added: %d, #removed: %d\n",
		commit.Author.Timestamp, commit.Author.Name,
		commit.ChangeSet.NumFilesChanged, commit.ChangeSet.LinesAdded, commit.ChangeSet.LinesRemoved)

	jsonBlob, _ := json.MarshalIndent(commit, "", "\t")
	fmt.Fprintf(os.Stdout, "[\n%s,\n]\n", jsonBlob)
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

func parseArgs() (string, bool, bool, bool) {

	help := flag.Bool("help", false, "This help text")
	filename := flag.String("filename", "git_history.json", "Json file with git history as created by git2json")
	analyzeCommitters := flag.Bool("analyze-committers", false, "Analyse commiters of project")
	analyzeFiles := flag.Bool("analyze-files", false, "Analyse files of project")
	dumpAsJson := flag.Bool("dumpAsJson-json", false, "Dump commits of project as json")

	flag.Parse()

	if help != nil && *help {
		printHelp()
	}

	if *analyzeCommitters == false && *analyzeFiles == false && *dumpAsJson == false {
		printHelp()
	}

	return *filename, *analyzeCommitters, *analyzeFiles, *dumpAsJson
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, " %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(-1)
}

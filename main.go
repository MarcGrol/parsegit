package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

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
		doAnalyseCommitters(commits)
	}

	if analyzeFiles {
		doAnalyseFiles(commits)
	}
}

func logCommit(commit Commit) {
	fmt.Fprintf(os.Stdout, "%s - %s - #files: %d, #added: %d, #removed: %d\n",
		commit.Author.Timestamp, commit.Author.Name,
		commit.ChangeSet.NumFilesChanged, commit.ChangeSet.LinesAdded, commit.ChangeSet.LinesRemoved)

	//jsonBlob, _ := json.MarshalIndent(commit, "", "\t")
	//fmt.Fprintf(os.Stdout, "[\n%s,\n]\n", jsonBlob)
}

func parseArgs() (string, bool, bool, bool) {

	help := flag.Bool("help", false, "This help text")
	filename := flag.String("filename", "git_history.json", "Json file with git history as created by git2json")
	analyzeCommitters := flag.Bool("analyze-committers", false, "Analyse commiters of project")
	analyzeFiles := flag.Bool("analyze-files", false, "Analyse files of project")
	dumpAsJson := flag.Bool("dump", false, "Dump commits of project as json")

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

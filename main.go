package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"text/template"
)

func main() {
	// read config
	filename, logTemplate, needToAnalyzeCommiters, needToAnalyzeFiles, analyzeFilesOverTime, debug := parseArgs()

	// parse and post-process git2json file
	commits := Commits{}
	parse(filename, func(commit Commit) {
		commits = append(commits, commit)
	})
	sort.Sort(ByTimestamp{commits})

	// dump to debug
	if debug {
		t, err := template.New("log-template").Parse(logTemplate)
		if err != nil {
			log.Fatalf("Error parsing template: %s", err)
		}
		for _, c := range commits {
			applyTemplate(c, t)
		}
	}

	if needToAnalyzeCommiters {
		analyseCommitters(commits)
	}

	if needToAnalyzeFiles {
		analyseFilesAsCsv(commits)
	}

	if analyzeFilesOverTime {
		result := analyseFilesOverTime(commits, 20)
		blob, _ := json.MarshalIndent(result, "", "\t")

		fmt.Fprintf(os.Stdout, "%s", string(blob))
	}
}

func applyTemplate(commit Commit, template *template.Template) {
	err := template.Execute(os.Stdout, commit)
	if err != nil {
		log.Fatalf("Error applying template: %s", err)
	}
}

func parseArgs() (string, string, bool, bool, bool, bool) {

	help := flag.Bool("help", false, "This help text")
	filename := flag.String("filename", "git_history.json", "Json file with git history as created by git2json")
	logTemplate := flag.String("template", "{{.Author.Timestamp}},{{.Author.Name}},{{.ChangeSet.NumFilesChanged}},{{.ChangeSet.LinesAdded}},{{.ChangeSet.LinesRemoved}}\n", "Logging template")
	analyzeCommitters := flag.Bool("analyze-committers", false, "Analyse commiters of project")
	analyzeFiles := flag.Bool("analyze-files", false, "Analyse files of project")
	analyzeFilesOverTime := flag.Bool("analyze-files-over-time", false, "Analyse files of project over time")
	dumpAsJson := flag.Bool("debug", false, "Dump details of all commits of project")

	flag.Parse()

	if help != nil && *help {
		printHelp()
	}

	if *analyzeCommitters == false && *analyzeFiles == false && *dumpAsJson == false && *analyzeFilesOverTime == false {
		printHelp()
	}

	return *filename, *logTemplate, *analyzeCommitters, *analyzeFiles, *analyzeFilesOverTime, *dumpAsJson
}

func printHelp() {
	fmt.Fprintf(os.Stderr, "\nUsage:\n")
	fmt.Fprintf(os.Stderr, " %s [flags]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(-1)
}

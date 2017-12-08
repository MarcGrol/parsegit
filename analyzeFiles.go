package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type fileSummaryForInterval struct {
	Start         time.Time
	End           time.Time
	FileSummaries []fileSummary
}

type fileSummariesForIntervals []fileSummaryForInterval

type fileSummary struct {
	Filename     string
	Extention    string
	NumCommits   int
	LinesAdded   int
	LinesRemoved int
	Commiters    map[string]int
}

type fileSummaries []fileSummary

func (fs fileSummaries) Len() int      { return len(fs) }
func (fs fileSummaries) Swap(i, j int) { fs[i], fs[j] = fs[j], fs[i] }

type ByNumCommits struct{ fileSummaries }

func (fs ByNumCommits) Less(i, j int) bool {
	return fs.fileSummaries[i].NumCommits > fs.fileSummaries[j].NumCommits
}

type comm struct {
	NumCommits int `json:"commits"`
}

func analyseFilesOverTime(commits Commits, limit int) []map[string]interface{} {
	allResults := []map[string]interface{}{}
	results := doAnalyseFilesOverTime(commits, limit)
	for _, r := range results {
		intervalResult := map[string]interface{}{}
		intervalResult["date"] = r.Start.Format("2006-01-02")
		for _, i := range r.FileSummaries {
			intervalResult[i.Filename] = comm{NumCommits: i.NumCommits}
		}
		allResults = append(allResults, intervalResult)
	}
	return allResults
}

const numTimeSlices = 100

func doAnalyseFilesOverTime(commits Commits, limit int) fileSummariesForIntervals {
	if len(commits) == 0 {
		return fileSummariesForIntervals{}
	}

	oldest := commits[0].Author.Timestamp
	newest := commits[len(commits)-1].Author.Timestamp
	diff := newest.Sub(oldest) / numTimeSlices

	intervals := fileSummariesForIntervals{}
	for i := 0; i < numTimeSlices-1; i++ {
		start := oldest.Add(time.Duration(i) * diff)
		end := oldest.Add((time.Duration(i + 1)) * diff)
		subset := commitsBetween(commits, start, end)
		fmt.Fprintf(os.Stderr, "*** Found %d commits from %s to %s\n", len(subset), start, end)

		topFiles := analyseFilesTop(subset, limit)
		intervals = append(intervals,
			fileSummaryForInterval{
				Start:         start,
				End:           end,
				FileSummaries: topFiles,
			})

	}

	return intervals
}

func commitsBetween(commits Commits, lowerBound, upperBound time.Time) Commits {
	result := Commits{}
	for _, c := range commits {
		if c.Author.Timestamp.After(lowerBound) && c.Author.Timestamp.Before(upperBound) {
			result = append(result, c)
		}
	}

	return result
}

func analyseFilesAsCsv(commits Commits) map[string]fileSummary {
	files := map[string]fileSummary{}
	// analyze work done on files
	for _, c := range commits {
		fileInfo(files, c)
	}
	// print as csv
	fmt.Fprintf(os.Stdout, "filename; #extension, #commits; #additions; #removals; #commiters\n")
	for fn, summ := range files {
		extention := ""
		parts := strings.Split(fn, ".")
		if len(parts) > 0 {
			extention = parts[len(parts)-1]
		}

		fmt.Fprintf(os.Stdout, "%s;%s;%d;%d;%d;%d\n",
			fn, extention, summ.NumCommits, summ.LinesAdded, summ.LinesRemoved, len(summ.Commiters))
	}
	return files
}

func analyseFilesTop(commits Commits, limit int) fileSummaries {
	files := map[string]fileSummary{}
	// analyze work done on files
	for _, c := range commits {
		fileInfo(files, c)
	}

	summaries := fileSummaries{}
	for _, f := range files {
		summaries = append(summaries, f)
	}
	sort.Sort(ByNumCommits{summaries})

	if len(summaries) < 10 {
		return summaries
	}

	return summaries[0:limit]
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
		sum.Filename = filename
		parts := strings.Split(filename, ".")
		if len(parts) > 0 {
			sum.Extention = parts[len(parts)-1]
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

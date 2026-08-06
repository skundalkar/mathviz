// Command digest prints a short daily "what got built and what it teaches"
// report in Markdown. The scheduled build task runs this after a cycle and
// sends the output to you — it is the tutoring channel. Run it yourself any
// time with `make digest` or `go run ./cmd/digest`.
//
// Usage:
//
//	go run ./cmd/digest [--since=24h]
//
// It is intentionally dependency-free: it shells out to git and reads two
// Markdown files, so it works anywhere the repo is checked out.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	since := flag.String("since", "24 hours ago", "git --since window for the commit list")
	flag.Parse()

	fmt.Println("# MathViz — daily digest")
	fmt.Println()

	commits := gitLog(*since)
	if len(commits) == 0 {
		fmt.Printf("No commits in the last window (%s). The loop may be idle or the backlog empty.\n", *since)
	} else {
		fmt.Printf("**%d commits** in the last window (%s):\n\n", len(commits), *since)
		for _, c := range commits {
			fmt.Printf("- %s\n", c)
		}
	}
	fmt.Println()

	done, total := backlogProgress()
	fmt.Printf("**Progress:** %d of %d concepts built.\n\n", done, total)

	if lesson := latestLearning(); lesson != "" {
		fmt.Println("## Today's lesson")
		fmt.Println()
		fmt.Println(lesson)
	}

	fmt.Println()
	fmt.Println("---")
	fmt.Println("_Open the gallery with `make serve` and turn the knobs. Review the diff, then merge if it looks right._")
}

func gitLog(since string) []string {
	out, err := exec.Command("git", "log", "--since="+since, "--pretty=format:%h  %s").Output()
	if err != nil {
		return nil
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

func backlogProgress() (done, total int) {
	b, err := os.ReadFile("BACKLOG.md")
	if err != nil {
		return 0, 0
	}
	for _, line := range strings.Split(string(b), "\n") {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "- [x]") {
			done++
			total++
		} else if strings.HasPrefix(l, "- [ ]") {
			total++
		}
	}
	return
}

// latestLearning returns the first "## ..." section from LEARNINGS.md (newest,
// since entries are prepended).
func latestLearning() string {
	b, err := os.ReadFile("LEARNINGS.md")
	if err != nil {
		return ""
	}
	lines := strings.Split(string(b), "\n")
	start := -1
	for i, l := range lines {
		if strings.HasPrefix(l, "## ") {
			start = i
			break
		}
	}
	if start < 0 {
		return ""
	}
	var out []string
	for i := start; i < len(lines); i++ {
		if i > start && (strings.HasPrefix(lines[i], "## ") || strings.TrimSpace(lines[i]) == "---") {
			break
		}
		out = append(out, lines[i])
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

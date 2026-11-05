package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	cc "github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
)

var (
	optionalLink       = `(fix\w*\s|close\w*\s|resolve\w*\s)?` // Optional issue linking prefix, see https://docs.github.com/en/issues/tracking-your-work-with-issues/linking-a-pull-request-to-an-issue.
	issueRegexFull     = regexp.MustCompile(`^` + optionalLink + `https://github\.com/piplabs/story/issues/\d+$`)
	issueRegexShort    = regexp.MustCompile(`^` + optionalLink + `#\d+$`)                       // e.g. "#1334"
	issueRegexCrossRef = regexp.MustCompile(`^` + optionalLink + `piplabs\/[a-zA-Z0-9-]+#\d+$`) // e.g. "piplabs/story-geth#1559"
	issueLineRegex     = regexp.MustCompile(`(?i)^\s*[-*]?\s*issue:\s*(.+)$`)                   // e.g. "issue: #1334" or "- issue: #1334"
)

// run runs the verification.
func run() error {
	pr, err := prFromEnv()
	if err != nil {
		return err
	}

	// Skip dependabot PRs.
	if strings.Contains(pr.Title, "deps") && strings.Contains(pr.Body, "dependabot") {
		return nil
	}

	log.Printf("Verifying story PR against template\n")
	log.Printf("PR Title: %s\n", pr.Title)
	log.Printf("## PR Body:\n%s\n####\n", pr.Body)

	return verify(pr.Title, pr.Body)
}

type PR struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	ID    string `json:"node_id"`
}

// prFromEnv returns the PR by parsing it from "GITHUB_PR" env var or an error.
func prFromEnv() (PR, error) {
	const prEnv = "GITHUB_PR"

	prJSON, ok := os.LookupEnv(prEnv)
	if !ok || strings.TrimSpace(prJSON) == "" {
		return PR{}, errors.New("env variable not set")
	}

	var pr PR

	if err := json.Unmarshal([]byte(prJSON), &pr); err != nil {
		return PR{}, err
	}

	if pr.Title == "" || pr.Body == "" || pr.ID == "" {
		return PR{}, errors.New("pr field not set")
	}

	return pr, nil
}

// verify returns an error if the PR title isn't a valid conventional commit
// or if the PR body doesn't reference a github issue.
func verify(title, body string) error {
	if err := verifyTitle(title); err != nil {
		return err
	}

	return verifyIssue(body)
}

// verifyTitle ensures the PR title follows the conventional commit style
// (e.g. "feat: ...", "fix(scope): ..."). Casing, punctuation and scope
// are intentionally not restricted beyond what the conventional commit spec requires.
func verifyTitle(title string) error {
	// Fix line endings, since conventional commit parser doesn't support CRLF.
	title = strings.ReplaceAll(title, "\r\n", "\n")

	const maxLen = 100
	if len(title) > maxLen {
		return errors.New("title too long")
	}

	m := parser.NewMachine()
	m.WithTypes(cc.TypesConventional)

	msg, err := m.Parse([]byte(title))
	if err != nil {
		return fmt.Errorf("title is not a conventional commit: %v", err)
	}

	commit, ok := msg.(*cc.ConventionalCommit)
	if !ok {
		return errors.New("title is not a conventional commit")
	}

	if !commit.Ok() {
		return errors.New("title is not a valid conventional commit")
	}

	return nil
}

// verifyIssue ensures the PR body footer contains a single `issue:` line referencing a
// github issue. The footer is the last paragraph of the body (the block after the final
// blank line), so the body above it may contain anything (e.g. bullet lists).
func verifyIssue(body string) error {
	body = strings.ReplaceAll(body, "\r\n", "\n")

	// Footer = last paragraph, i.e. the block after the final blank line.
	footer := strings.TrimSpace(body)
	if idx := strings.LastIndex(footer, "\n\n"); idx >= 0 {
		footer = strings.TrimSpace(footer[idx+len("\n\n"):])
	}

	var issues []string
	for _, line := range strings.Split(footer, "\n") {
		if matches := issueLineRegex.FindStringSubmatch(line); matches != nil {
			issues = append(issues, strings.TrimSpace(matches[1]))
		}
	}

	if len(issues) == 0 {
		return errors.New("missing `issue` section in footer. Please add an `issue:` line at the end of the PR body")
	}

	if len(issues) != 1 {
		return errors.New("invalid number of issue sections, only one allowed")
	}

	// The issue value is never empty: issueLineRegex only captures non-empty content.
	issue := issues[0]
	//nolint:nestif // nested ifs readability
	if issue == "none" {
		// None is fine
	} else if issueRegexFull.MatchString(issue) {
		// Full issue URL
	} else if issueRegexShort.MatchString(issue) {
		// Short issue URL
	} else if issueRegexCrossRef.MatchString(issue) {
		// Cross-repo (same org) issue URL
	} else {
		return errors.New("invalid issue section")
	}

	return nil
}

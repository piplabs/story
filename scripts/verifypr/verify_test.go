package main

import (
	"os"
	"strings"
	"testing"
)

func TestVerify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		title   string
		body    string
		wantErr bool
	}{
		{
			name:  "conventional title with caps, punctuation and long description",
			title: "feat: Add Foo Bar! (with CAPS) & symbols, a fairly long description that goes past eighty chars",
			body:  "Some description.\n\nissue: #1234",
		},
		{
			name:  "no scope is allowed",
			title: "fix: do the thing",
			body:  "issue: #1",
		},
		{
			name:  "scope is allowed",
			title: "fix(dkg): do the thing",
			body:  "issue: #1",
		},
		{
			name:    "non-conventional title rejected",
			title:   "Random title without a type",
			body:    "issue: #1",
			wantErr: true,
		},
		{
			name:    "title over 100 chars rejected",
			title:   "feat: " + strings.Repeat("a", 100),
			body:    "issue: #1",
			wantErr: true,
		},
		{
			name:  "issue as bullet point is allowed",
			title: "chore: cleanup",
			body:  "- some change\n- issue: #1234",
		},
		{
			name:  "issue in footer after bullet body",
			title: "chore: cleanup",
			body:  "- did a thing\n- fixed a bug\n\nissue: #1234",
		},
		{
			name:  "issue in footer after markdown body",
			title: "chore: cleanup",
			body:  "## Description\n\n- foo\n- bar\n\nissue: #1234",
		},
		{
			name:    "issue not in footer (content below) rejected",
			title:   "chore: cleanup",
			body:    "issue: #1234\n\n- a follow up note",
			wantErr: true,
		},
		{
			// Only the last paragraph (footer) is inspected; an `issue:` line
			// elsewhere in the body is ignored, so the footer's single line wins.
			name:  "issue line outside footer is ignored",
			title: "chore: cleanup",
			body:  "issue: #1 referenced in prose\n\nissue: #2",
		},
		{
			name:  "issue with linking prefix is allowed",
			title: "chore: cleanup",
			body:  "issue: fixes #1234",
		},
		{
			name:  "cross-repo issue is allowed",
			title: "chore: cleanup",
			body:  "issue: piplabs/story-geth#1559",
		},
		{
			name:  "full issue url is allowed",
			title: "chore: cleanup",
			body:  "issue: https://github.com/piplabs/story/issues/1234",
		},
		{
			name:  "issue none is allowed",
			title: "chore: cleanup",
			body:  "issue: none",
		},
		{
			name:    "missing issue rejected",
			title:   "chore: cleanup",
			body:    "just a body without an issue line",
			wantErr: true,
		},
		{
			name:    "two issue lines rejected",
			title:   "chore: cleanup",
			body:    "issue: #1\nissue: #2",
			wantErr: true,
		},
		{
			name:    "invalid issue value rejected",
			title:   "chore: cleanup",
			body:    "issue: not-an-issue",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := verify(tt.title, tt.body)
			if tt.wantErr && err == nil {
				t.Fatalf("verify(%q, %q) = nil, want error", tt.title, tt.body)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("verify(%q, %q) = %v, want nil", tt.title, tt.body, err)
			}
		})
	}
}

func TestPrFromEnv(t *testing.T) {
	// Note: cannot use t.Parallel with t.Setenv.
	tests := []struct {
		name    string
		set     bool
		val     string
		want    PR
		wantErr bool
	}{
		{
			name:    "env not set",
			set:     false,
			wantErr: true,
		},
		{
			name:    "env blank",
			set:     true,
			val:     "   ",
			wantErr: true,
		},
		{
			name:    "invalid json",
			set:     true,
			val:     "{not json}",
			wantErr: true,
		},
		{
			name:    "missing required field",
			set:     true,
			val:     `{"title":"feat: x","body":"issue: #1"}`,
			wantErr: true,
		},
		{
			name: "valid pr",
			set:  true,
			val:  `{"title":"feat: x","body":"issue: #1","node_id":"abc"}`,
			want: PR{Title: "feat: x", Body: "issue: #1", ID: "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				t.Setenv("GITHUB_PR", tt.val)
			} else {
				os.Unsetenv("GITHUB_PR")
			}

			got, err := prFromEnv()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("prFromEnv() = %+v, want error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("prFromEnv() error = %v, want nil", err)
			}
			if got != tt.want {
				t.Fatalf("prFromEnv() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	// Note: cannot use t.Parallel with t.Setenv.
	tests := []struct {
		name    string
		val     string
		wantErr bool
	}{
		{
			name: "valid pr passes",
			val:  `{"title":"feat: x","body":"issue: #1","node_id":"abc"}`,
		},
		{
			name: "dependabot pr skipped",
			val:  `{"title":"chore(deps): bump x","body":"made by dependabot","node_id":"abc"}`,
		},
		{
			name:    "env error propagates",
			val:     "",
			wantErr: true,
		},
		{
			name:    "verification failure propagates",
			val:     `{"title":"not conventional","body":"issue: #1","node_id":"abc"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.val == "" {
				os.Unsetenv("GITHUB_PR")
			} else {
				t.Setenv("GITHUB_PR", tt.val)
			}

			err := run()
			if tt.wantErr && err == nil {
				t.Fatalf("run() = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("run() = %v, want nil", err)
			}
		})
	}
}

package parser

import (
	"regexp"
	"strings"
)

type Event struct {
	SrcIP     string
	Username  string
	EventType string
	Raw       string
}

var (
	ipPattern = regexp.MustCompile(`from (\d{1,3}(?:\.\d{1,3}){3})`)

	// Extracts username from common auth.log failure patterns
	usernamePatterns = []*regexp.Regexp{
		regexp.MustCompile(`Failed password for (?:invalid user )?(\S+) from`),
		regexp.MustCompile(`Invalid user (\S+) from`),
		regexp.MustCompile(`Connection closed by invalid user (\S+)\s`),
	}

	failPatterns = []string{
		"Failed password",
		"Invalid user",
		"Connection closed by invalid user",
		"authentication failure",
		"BREAK-IN ATTEMPT",
	}
)

// Parse extracts a login failure event from an auth.log line.
// Returns nil if the line is not a failure event.
func Parse(line string) *Event {
	eventType := matchEventType(line)
	if eventType == "" {
		return nil
	}

	m := ipPattern.FindStringSubmatch(line)
	if m == nil {
		return nil
	}

	return &Event{
		SrcIP:     m[1],
		Username:  extractUsername(line),
		EventType: eventType,
		Raw:       line,
	}
}

func matchEventType(line string) string {
	for _, p := range failPatterns {
		if strings.Contains(line, p) {
			return p
		}
	}
	return ""
}

func extractUsername(line string) string {
	for _, re := range usernamePatterns {
		if m := re.FindStringSubmatch(line); m != nil {
			return m[1]
		}
	}
	return "unknown"
}

package git

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/adevinta/maiao/pkg/log"
)

const fixupPrefix = "fixup! "
const changeIDHeader = "Change-Id"

// Message defines the model of a commit message
type Message struct {
	Title   string
	Body    string
	Headers map[string]string
}

var headerRe = regexp.MustCompile("^([^:]+):([^:/]+)$")

// Parse parses a commit message string into a structured object
func Parse(message string) *Message {
	scanner := bufio.NewScanner(strings.NewReader(message))
	scanner.Scan()
	m := &Message{
		Title:   scanner.Text(),
		Headers: map[string]string{},
	}
	l := log.Logger.WithField("commit message", message)
	sep := ""
	for scanner.Scan() {
		line := scanner.Text()
		submatch := headerRe.FindSubmatch([]byte(line))
		if len(submatch) == 3 {
			name := strings.Trim(string(submatch[1]), " ")
			value := strings.Trim(string(submatch[2]), " ")
			l.WithField("line", line).WithField("header", name).WithField("value", value).Trace("header found in commit message")
			m.Headers[name] = value
		} else {
			m.Body += sep + line
		}
		sep = "\n"
	}
	m.Body = strings.Trim(m.Body, "\n")
	strings.Split(message, "\n")
	return m
}

func (m *Message) String() string {
	if m == nil {
		return ""
	}
	s := m.Title
	if m.Body != "" || len(m.Headers) > 0 {
		s += "\n\n"
	}
	if m.Body != "" {
		s += m.Body
		if len(m.Headers) > 0 {
			s += "\n\n"
		}
	}
	for k, v := range m.Headers {
		s += fmt.Sprintf("%s:%s", k, v)
	}
	return s
}

// IsFixup returns if a commit is a fixup of another commit
func (m *Message) IsFixup() bool {
	if m == nil {
		return false
	}
	return isFixupTitle(m.Title)
}

// GetTitle returns the commit title after stripping all fixup prefixes
func (m *Message) GetTitle() string {
	if m == nil {
		return ""
	}
	t := m.Title
	for isFixupTitle(t) {
		t = t[len(fixupPrefix):]
	}
	return t
}

// GetChangeID returns the changeID in the commit message and if it has been found
func (m *Message) GetChangeID() (changeID string, ok bool) {
	if m == nil {
		return "", false
	}
	if m.Headers == nil {
		return "", false
	}
	changeID, ok = m.Headers[changeIDHeader]
	return
}

func isFixupTitle(title string) bool {
	return strings.HasPrefix(strings.ToLower(title), fixupPrefix)
}

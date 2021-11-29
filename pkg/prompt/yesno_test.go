package prompt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setStdInOut(in, out *os.File) {
	stdin = in
	stdout = out
}
func TestPromptOutputsOnStderr(t *testing.T) {
	assert.Equal(t, os.Stderr, stdout)
}

func TestPromptOutput(t *testing.T) {
	defer setStdInOut(stdin, stdout)
	inr, inw, _ := os.Pipe()
	outr, outw, _ := os.Pipe()
	setStdInOut(inr, outw)
	go YesNo("This is a question?")
	q := make([]byte, 100)
	// wait until the YesNo routine has written something on the output
	for c := 0; c < 10; c, _ = outr.Read(q) {
	}
	t.Run("output contains the original question", func(t *testing.T) {
		assert.Contains(t, string(q), "This is a question?")
	})
	t.Run("output contains the [y/N] selector", func(t *testing.T) {
		assert.Contains(t, string(q), "[y/N]")
	})
	// unlock the YesNo routine
	inw.Write([]byte("y\n"))
}

func TestPromptYesInputReturnsTrue(t *testing.T) {
	testPromptReturnsValue(t, "y\n", true)
	testPromptReturnsValue(t, "Y\n", true)
	testPromptReturnsValue(t, "yes\n", true)
	testPromptReturnsValue(t, "Yes\n", true)
	testPromptReturnsValue(t, "YeS\n", true)
	testPromptReturnsValue(t, "YES\n", true)
	testPromptReturnsValue(t, "yES\n", true)
}
func TestPromptNoInputReturnsFalse(t *testing.T) {
	testPromptReturnsValue(t, "n\n", false)
	testPromptReturnsValue(t, "N\n", false)
	testPromptReturnsValue(t, "No\n", false)
	testPromptReturnsValue(t, "nO\n", false)
	testPromptReturnsValue(t, "NO\n", false)
	testPromptReturnsValue(t, "no\n", false)
	testPromptReturnsValue(t, "kamoulox\n", false)
}

func testPromptReturnsValue(t testing.TB, enteredValue string, expected bool) {
	defer setStdInOut(stdin, stdout)
	inr, inw, _ := os.Pipe()
	_, outw, _ := os.Pipe()
	setStdInOut(inr, outw)
	go func() {
		for {
			_, err := inw.Write([]byte(enteredValue))
			if err != nil {
				break
			}
		}
	}()
	assert.Equal(t, expected, YesNo("This is a question?"))
}

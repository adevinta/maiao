package prompt

import (
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

var stdin, stdout *os.File

// YesNo prompts your question and awaits user response
func YesNo(question string) bool {
	prompt := promptui.Prompt{
		Label:     question,
		IsConfirm: true,
		Stdin:     stdin,
		Stdout:    stdout,
	}
	result, _ := prompt.Run()
	result = strings.ToLower(result)
	return result == "y" || result == "yes"
}

func init() {
	stdin = os.Stdin
	stdout = os.Stderr
}

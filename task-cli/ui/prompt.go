package ui

import (
	"bufio"
	"fmt"
	"strings"
)

func PromptChoice(prompt string, reader *bufio.Reader) string {
	fmt.Println(prompt)
	fmt.Print(">>")

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func PromptRaw(prompt string, reader *bufio.Reader) string {
	fmt.Println(prompt)
	fmt.Print(">>")

	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func ConfirmRemoval(taskName string, reader *bufio.Reader) bool {
	confirmation := PromptChoice(fmt.Sprintf("Are you sure you want to remove task \"%s\"? (y/n)", taskName), reader)
	return confirmation == "y" || confirmation == "yes"
}
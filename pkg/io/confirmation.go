package io

import (
	"fmt"
	libio "io"
	"os"
	"strings"
)

func ReadOsConfirmation(prompt string) (bool, error) {
	return ReadConfirmation(prompt, os.Stdin, os.Stdout)
}

func ReadConfirmation(prompt string, src libio.Reader, dst libio.Writer) (choice bool, err error) {
	var response string

	for {
		fmt.Fprint(dst, prompt, " [y/n]: ")

		_, err = fmt.Fscanln(src, &response)
		if err != nil {
			break
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			choice = true
			break
		} else if response == "n" || response == "no" {
			break
		}
	}

	return
}

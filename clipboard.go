package main

import (
	"github.com/atotto/clipboard"
)

// ClipboardCopy puts the given text into the system clipboard.
// It works on Windows, macOS and Linux (Linux needs xclip or xsel).
func ClipboardCopy(text string) error {
	return clipboard.WriteAll(text)
}

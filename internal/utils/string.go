package utils

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func PadStringToWidth(s string, targetWidth int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= targetWidth {
		return s
	}
	padding := targetWidth - currentWidth
	return s + strings.Repeat(" ", padding)
}

func TruncateStringWithWidth(s string, maxWidth int) string {
	if runewidth.StringWidth(s) <= maxWidth {
		return s
	}

	var result []rune
	currentWidth := 0

	for _, r := range s {
		runeWidth := runewidth.RuneWidth(r)
		if currentWidth+runeWidth > maxWidth-3 {
			break
		}
		result = append(result, r)
		currentWidth += runeWidth
	}

	return string(result) + "..."
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

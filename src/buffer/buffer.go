package buffer

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type Buffer struct {
	Contents string
}

func (b *Buffer) splitContentsByLineEnd() []string {
	return strings.Split(b.Contents, "\n")
}

func (b *Buffer) updateContentsFromSplits(splits []string) {
	b.Contents = strings.Join(splits, "\n")
}

func (b *Buffer) WriteToBuf(ch rune, x, y int) error {
	splits := b.splitContentsByLineEnd()
	if y >= len(splits) {
		return errors.New("internal error")
	}
	stringToEdit := splits[y]
	if x > len(stringToEdit) {
		return errors.New("internal error")
	}
	if x == len(stringToEdit) {
		stringToEdit += string(ch)
	} else {
		stringToEdit = stringToEdit[0:x] + string(ch) + stringToEdit[x:]
	}
	splits[y] = stringToEdit
	b.updateContentsFromSplits(splits)
	return nil
}

func (b *Buffer) DeleteFromBuf(x, y int) error {
	if b.Contents == "" {
		return nil
	}
	splits := b.splitContentsByLineEnd()
	if y >= len(splits) {
		return errors.New("internal error")
	}
	stringToEdit := splits[y]
	if x > len(stringToEdit) {
		return errors.New("internal error")
	}
	if x == 0 {
		if y-1 >= 0 {
			splits[y-1] += stringToEdit
			splits = append(splits[0:y], splits[y+1:]...)
		}
	} else {
		if x == len(stringToEdit) {
			stringToEdit = stringToEdit[0 : x-1]
		} else {
			stringToEdit = stringToEdit[0:x-1] + stringToEdit[x:]
		}
		splits[y] = stringToEdit
	}
	b.updateContentsFromSplits(splits)
	return nil
}

func (b *Buffer) GetLineEndX(y int) (int, error) {
	splits := b.splitContentsByLineEnd()
	if y < 0 {
		return -1, errors.New("internal error")
	}
	if y >= len(splits) {
		return -1, errors.New("internal error")
	}

	return len(splits[y]), nil
}

func (b *Buffer) GetLineStartX(y int) (int, error) {
	splits := b.splitContentsByLineEnd()
	if y >= len(splits) {
		return -1, errors.New("internal error")
	}
	var idx int
	for i, ch := range splits[y] {
		idx = i
		if ch != '\t' {
			break
		}
	}
	return idx, nil
}

func (b *Buffer) GetNextWordPos(x, y int) (int, int, error) {
	lines := b.splitContentsByLineEnd()

	if y < 0 || y >= len(lines) {
		return -1, -1, errors.New("internal error")
	}

	for lineIdx := y; lineIdx < len(lines); lineIdx++ {
		line := lines[lineIdx]
		startX := x
		if lineIdx != y {
			x, err := b.GetLineStartX(lineIdx)
			if err != nil {
				return -1, -1, fmt.Errorf("get next word pos: %s", err)
			}
			return x, lineIdx, nil
		}

		i := startX
		n := len(line)

		// Skip over current word or non-word characters
		for i < n && !unicode.IsLetter(rune(line[i])) && !unicode.IsDigit(rune(line[i])) {
			i++
		}
		for i < n && (unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i]))) {
			i++
		}
		// Now skip any spaces or punctuation to get to the next word
		for i < n && !unicode.IsLetter(rune(line[i])) && !unicode.IsDigit(rune(line[i])) {
			i++
		}
		if i < n {
			return i, lineIdx, nil
		}
	}

	// If no next word found
	return x, y, nil
}

func (b *Buffer) GetNextWordEndPos(x, y int) (int, int, error) {
	lines := b.splitContentsByLineEnd()

	if y < 0 || y >= len(lines) {
		return 0, 0, errors.New("internal error")
	}

	for lineIdx := y; lineIdx < len(lines); lineIdx++ {
		line := lines[lineIdx]
		startX := x
		if lineIdx != y {
			x, err := b.GetLineStartX(lineIdx)
			if err != nil {
				return -1, -1, fmt.Errorf("get next word pos: %s", err)
			}
			startX = x
		}

		i := startX
		n := len(line)

		// If we're in the middle of a word, move to its end
		if i < n && unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i])) {
			for i < n && (unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i]))) {
				i++
			}
			if i-1 != startX {
				return i - 1, lineIdx, nil
			}
			i++
		}

		// Otherwise, skip non-word characters to find next word
		for i < n && !unicode.IsLetter(rune(line[i])) && !unicode.IsDigit(rune(line[i])) {
			i++
		}
		// Now go to the end of this word
		start := i
		for i < n && (unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i]))) {
			i++
		}
		if start < i {
			return i - 1, lineIdx, nil
		}
	}

	// No more words found
	if len(lines) == 0 {
		return x, y, nil
	}

	lastLine := len(lines) - 1
	lastLen := len(lines[lastLine])
	if lastLen == 0 {
		return -1, lastLine, nil
	}
	return lastLen - 1, lastLine, nil
}

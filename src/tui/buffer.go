package tui

import (
	"errors"
	"strings"
)

type Buffer struct {
	contents string
}

func (b *Buffer) splitContentsByLineEnd() []string {
	return strings.Split(b.contents, "\n")
}

func (b *Buffer) updateContentsFromSplits(splits []string) {
	b.contents = strings.Join(splits, "\n")
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

func (b *Buffer) GetLineEndX(y int) (int, error) {
	splits := b.splitContentsByLineEnd()
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

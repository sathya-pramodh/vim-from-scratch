package buffer

import (
	"errors"
	"strings"
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

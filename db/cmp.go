package db

import (
	"strconv"

	"github.com/syndtr/goleveldb/leveldb/comparer"
)

type MyComparator struct {
}

func (cmp MyComparator) Name() string {
	return "my diy cmp"
}
func (cmp MyComparator) Compare(a, b []byte) int {
	i, err1 := strconv.Atoi(string(a))
	j, err2 := strconv.Atoi(string(b))
	if err1 != nil || err2 != nil {
		return comparer.DefaultComparer.Compare(a, b)
	}
	if i == j {
		return 0
	} else if i < j {
		return -1
	}
	return 1
}
func (cmp MyComparator) Separator(dst, a, b []byte) []byte {
	return nil
}

func (cmp MyComparator) Successor(dst, b []byte) []byte {
	return nil
}

package enumerators

import (
	"sync"

	"github.com/deckarep/golang-set"
)

type EnumError struct {
	err        error
	enumerator string
}

func (enumerr EnumError) Error() string {
	return "[ERROR]: " + enumerr.enumerator + " " + enumerr.err.Error()
}

type Enumerator func(string, chan<- mapset.Set, chan<- EnumError, *sync.WaitGroup)

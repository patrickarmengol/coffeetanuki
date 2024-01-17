package data

import (
	"fmt"
	"strings"

	"github.com/patrickarmengol/coffeetanuki/internal/validator"
)

type SearchQuery struct {
	Term            string
	Sort            string
	SortableColumns []string // TODO: change to hashset?
}

func (sq SearchQuery) Validate(v *validator.Validator) {
	sortByAsc, okAsc := strings.CutSuffix(sq.Sort, "_asc")
	sortByDesc, okDesc := strings.CutSuffix(sq.Sort, "_desc")
	var sortBy string
	if okAsc {
		sortBy = sortByAsc
	} else {
		sortBy = sortByDesc
	}
	v.CheckField(okAsc || okDesc, "sort", "invalid sort format")
	v.CheckField(validator.PermittedValue(sortBy, sq.SortableColumns...), "sort", "invalid sort column")
}

func (sq SearchQuery) termWords() []string {
	return strings.Fields(sq.Term)
}

func (sq SearchQuery) termWordsWrapped() []string {
	wrappedWords := []string{}
	for _, w := range sq.termWords() {
		wrappedWords = append(wrappedWords, fmt.Sprintf("%%%s%%", w))
	}
	return wrappedWords
}

func (sq SearchQuery) sortBy() string {
	sortByAsc, okAsc := strings.CutSuffix(sq.Sort, "_asc")
	sortByDesc, okDesc := strings.CutSuffix(sq.Sort, "_desc")
	var sortBy string
	if okAsc {
		sortBy = sortByAsc
	} else if okDesc {
		sortBy = sortByDesc
	} else {
		// sort should have been validated
		panic("invalid sort format: " + sq.Sort)
	}
	for _, col := range sq.SortableColumns {
		if sortBy == col {
			return sortBy
		}
	}

	// sort should have been validated
	panic("invalid sort column: " + sq.Sort)
}

func (sq SearchQuery) sortDir() string {
	if strings.HasSuffix(sq.Sort, "_asc") {
		return "ASC"
	} else if strings.HasSuffix(sq.Sort, "_desc") {
		return "DESC"
	} else {
		// sort should have been validated
		panic("invalid sort format: " + sq.Sort)
	}
}

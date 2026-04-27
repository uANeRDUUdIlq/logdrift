package jsonpivot

import "errors"

var errEmptyField = errors.New("jsonpivot: array, keyField and valField must all be non-empty")

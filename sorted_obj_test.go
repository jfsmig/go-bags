// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"sort"
	"testing"
)

type Obj struct {
	pk int64
}

func (o *Obj) PK() int64 { return o.pk }

func TestObj_Unsorted(T *testing.T) {
	var bag SortedObj[int64, *Obj]
	bag.Add(&Obj{3})
	bag.Add(&Obj{1})
	bag.Add(&Obj{0})
	bag.Add(&Obj{2})
	T.Log("bag", bag)
	bag.Assert()
}

func TestObj_Sorted(T *testing.T) {
	var bag SortedObj[int64, *Obj]
	bag.Add(&Obj{0})
	bag.Add(&Obj{1})
	bag.Add(&Obj{2})
	bag.Add(&Obj{3})
	T.Log("bag", bag)
	bag.Assert()

	s := bag.Slice(1, 2)
	T.Log("slice", s)
	if len(s) != 2 {
		panic("bad length")
	} else if s[0].PK() != 2 && s[1].PK() != 3 {
		panic("bad content")
	}
}

// Assert panics if Check returns an error
func (s SortedObj[int64, T]) Assert() {
	if err := s.Check(); err != nil {
		panic(err)
	}
}

// Check validates the ordering and the unicity of the elements in the array
func (s SortedObj[int64, T]) Check() error {
	if !sort.IsSorted(s) {
		return ErrUnsorted
	}
	if !s.areItemsUnique() {
		return ErrDuplicates
	}
	return nil
}

// areItemsUnique validates the unicity of the elements in the array
func (s SortedObj[int64, T]) areItemsUnique() bool {
	lastId := s[0].PK()
	for _, a := range s[1:] {
		if lastId == a.PK() {
			return false
		}
		lastId = a.PK()
	}
	return true
}

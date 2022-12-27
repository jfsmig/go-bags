// Copyright (c) 2018-2023 Jean-Francois SMIGIELSKI
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

func TestObj_RemoveSorted(T *testing.T) {
	bag := SortedObj[int64, *Obj]{&Obj{0}, &Obj{1}, &Obj{2}, &Obj{3}}
	for _, v := range []int64{-1, 4} {
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
	for _, v := range []int64{0, 1, 2, 3} {
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
}

func TestObj_RemoveUnsorted(T *testing.T) {
	bag := SortedObj[int64, *Obj]{&Obj{0}, &Obj{1}, &Obj{2}, &Obj{3}}
	for _, v := range []int64{-1, 4} {
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
	for _, v := range []int64{3, 2, 1, 0} {
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
}

func TestObj_Lookup(T *testing.T) {
	bag := SortedObj[int64, *Obj]{&Obj{0}, &Obj{1}, &Obj{2}, &Obj{3}}
	for idx, v := range []int64{0, 1, 2, 3} {
		if !bag.Has(v) {
			T.Fatal()
		}
		if x, ok := bag.Get(v); !ok {
			T.Fatal()
		} else if x.PK() != v {
			T.Fatal()
		}
		if idx != bag.GetIndex(v) {
			T.Fatal()
		}
	}
	for _, v := range []int64{-1, -2, 5, 6} {
		if bag.Has(v) {
			T.Fatal()
		}
		if _, ok := bag.Get(v); ok {
			T.Fatal()
		}
		if -1 != bag.GetIndex(v) {
			T.Fatal()
		}
	}
}

func TestObj_Slice(T *testing.T) {
	bag := SortedObj[int64, *Obj]{&Obj{0}, &Obj{1}, &Obj{2}, &Obj{3}}
	testSlice := func(marker int64, max uint32, expectations ...*Obj) {
		slice := bag.Slice(marker, max)
		if len(slice) != len(expectations) {
			T.Fatal()
		}
		for i, v := range slice {
			if v.PK() != expectations[i].PK() {
				T.Fatal()
			}
		}
	}
	testSlice(bag[0].PK(), 2, bag[1], bag[2])
	testSlice(bag[0].PK()-1, 2, bag[0], bag[1])
	testSlice(bag[len(bag)-1].PK(), 1)
	testSlice(bag[0].PK()-1, MinSliceSize-1, bag[:MinSliceSize]...)
	testSlice(bag[0].PK()-1, MaxSliceSize+1, bag...)
}

func TestObj_SliceLarge(T *testing.T) {
	const total = 3 * MaxSliceSize
	var bag SortedObj[int64, *Obj]
	for i := int64(0); i < total; i++ {
		bag.Add(&Obj{i})
	}
	bag.Assert()
	slice := bag.Slice(bag[0].PK(), total)
	if len(slice) < MinSliceSize || len(slice) > MaxSliceSize {
		T.Fatal()
	}
}

func TestObj_AppendUnsorted(T *testing.T) {
	var bag SortedObj[int64, *Obj]
	bag.Append(&Obj{3}, &Obj{1}, &Obj{0}, &Obj{2})
	bag.Assert()
	if bag.Len() != 4 {
		T.Fatalf("invalid length")
	}
}

func TestObj_AppendSorted(T *testing.T) {
	var bag SortedObj[int64, *Obj]
	bag.Append(&Obj{0}, &Obj{1}, &Obj{2}, &Obj{3})
	bag.Assert()
	if bag.Len() != 4 {
		T.Fatal()
	}
}

func TestObj_AddUnsorted(T *testing.T) {
	var bag SortedObj[int64, *Obj]
	bag.Add(&Obj{3})
	bag.Add(&Obj{1})
	bag.Add(&Obj{0})
	bag.Add(&Obj{2})
	bag.Assert()
}

func TestObj_AddSorted(T *testing.T) {
	var bag SortedObj[int64, *Obj]
	bag.Add(&Obj{0})
	bag.Add(&Obj{1})
	bag.Add(&Obj{2})
	bag.Add(&Obj{3})
	bag.Assert()
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
		return errUnsorted
	}
	if !s.areItemsUnique() {
		return errDuplicates
	}
	return nil
}

// areItemsUnique validates the unicity of the elements in the array
func (s SortedObj[int64, T]) areItemsUnique() bool {
	if s.Len() < 2 {
		return true
	}
	lastId := s[0].PK()
	for _, a := range s[1:] {
		if lastId == a.PK() {
			return false
		}
		lastId = a.PK()
	}
	return true
}

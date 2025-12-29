// Copyright (c) 2018-2023 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"cmp"
	"sort"
	"testing"
)

type CmpInt int

func (i0 CmpInt) Compare(i1 CmpInt) int {
	return int(i0) - int(i1)
}

func TestCmp_RemoveSorted(T *testing.T) {
	bag := SortedCmp[CmpInt]{0, 1, 2, 3}
	for _, v := range []CmpInt{-1, 4} {
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
	for _, v := range []CmpInt{0, 1, 2, 3} {
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
}

func TestCmp_RemoveUnsorted(T *testing.T) {
	bag := SortedCmp[CmpInt]{0, 1, 2, 3}
	for _, v := range []CmpInt{-1, 4} {
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
	for _, v := range []CmpInt{3, 2, 1, 0} {
		bag.Remove(v)
		if bag.Has(v) {
			T.Fatal()
		}
		bag.Assert()
	}
}

func TestCmp_Lookup(T *testing.T) {
	bag := SortedCmp[CmpInt]{0, 1, 2, 3}
	for idx, v := range []CmpInt{0, 1, 2, 3} {
		if !bag.Has(v) {
			T.Fatal()
		}
		if x, ok := bag.Get(v); !ok {
			T.Fatal()
		} else if x != v {
			T.Fatal()
		}
		if idx != bag.GetIndex(v) {
			T.Fatal()
		}
	}
	for _, v := range []CmpInt{-1, -2, 5, 6} {
		if bag.Has(v) {
			T.Fatal()
		}
		if x, ok := bag.Get(v); ok {
			T.Fatal()
		} else if x != 0 {
			T.Fatal()
		}
		if -1 != bag.GetIndex(v) {
			T.Fatal()
		}
	}
}

func TestCmp_Slice(T *testing.T) {
	bag := SortedCmp[CmpInt]{0, 1, 2, 3}
	testSlice := func(marker CmpInt, max uint32, expectations ...CmpInt) {
		slice := bag.Slice(marker, max)
		if len(slice) != len(expectations) {
			T.Fatal()
		}
		for i, v := range slice {
			if v != expectations[i] {
				T.Fatal()
			}
		}
	}
	testSlice(bag[0], 2, bag[1], bag[2])
	testSlice(bag[0]-1, 2, bag[0], bag[1])
	testSlice(bag[len(bag)-1], 1)
	testSlice(bag[0]-1, MinSliceSize-1, bag[:MinSliceSize]...)
	testSlice(bag[0]-1, MaxSliceSize+1, bag...)
}

func TestCmp_SliceLarge(T *testing.T) {
	const total = 3 * MaxSliceSize
	var bag SortedCmp[CmpInt]
	for i := 0; i < total; i++ {
		bag.Add(CmpInt(i))
	}
	bag.Assert()
	slice := bag.Slice(bag[0], total)
	if len(slice) < MinSliceSize || len(slice) > MaxSliceSize {
		T.Fatal()
	}
}

func TestCmp_AppendUnsorted(T *testing.T) {
	var bag SortedCmp[CmpInt]
	bag.Append(3, 1, 0, 2)
	T.Log("bag", bag)
	bag.Assert()
	if bag.Len() != 4 {
		T.Fatalf("invalid length")
	}
}

func TestCmp_AppendSorted(T *testing.T) {
	var bag SortedCmp[CmpInt]
	bag.Append(0, 1, 2, 3)
	T.Log("bag", bag)
	bag.Assert()
	if bag.Len() != 4 {
		T.Fatal()
	}
}

func TestCmp_AddUnsorted(T *testing.T) {
	var bag SortedCmp[CmpInt]
	bag.Add(3)
	bag.Add(1)
	bag.Add(0)
	bag.Add(2)
	T.Log("bag", bag)
	bag.Assert()
}

func TestCmp_AddSorted(T *testing.T) {
	var bag SortedCmp[CmpInt]
	bag.Add(0)
	bag.Add(1)
	bag.Add(2)
	bag.Add(3)
	T.Log("bag", bag)
	bag.Assert()

	s := bag.Slice(1, 2)
	T.Log("slice", s)
	if len(s) != 2 {
		panic("bad length")
	} else if s[0] != 2 && s[1] != 3 {
		panic("bad content")
	}
}

// Assert panics if Check returns an error
func (s SortedCmp[T]) Assert() {
	if err := s.Check(); err != nil {
		panic(err)
	}
}

// Check validates the ordering and the unicity of the elements in the array
func (s SortedCmp[T]) Check() error {
	if !sort.IsSorted(s) {
		return errUnsorted
	}
	if !s.areItemsUnique() {
		return errDuplicates
	}
	return nil
}

// areItemsUnique validates the unicity of the elements in the array
func (s SortedCmp[T]) areItemsUnique() bool {
	if s.Len() < 2 {
		return true
	}
	lastId := s[0]
	for _, a := range s[1:] {
		if lastId.Compare(a) == 0 {
			return false
		}
		lastId = a
	}
	return true
}

type Cmp2Int struct {
	A, B int
}

func (x Cmp2Int) Compare(o Cmp2Int) int {
	if x.A == o.A {
		return cmp.Compare(x.B, o.B)
	}
	return cmp.Compare(x.A, o.A)
}

func TestCmd_SearchItem(T *testing.T) {
	bag := make(SortedCmp[Cmp2Int], 0)
	bag.Add(Cmp2Int{1, 1})
	bag.Add(Cmp2Int{1, 2})
	bag.Add(Cmp2Int{1, 3})
	bag.Add(Cmp2Int{3, 1})
	bag.Add(Cmp2Int{3, 2})
	bag.Add(Cmp2Int{3, 3})

	if idx := bag.SearchPredicate(func(x *Cmp2Int) bool { return x.A >= 1 }); idx != 0 {
		T.Fatal("idx", idx, "bag", bag)
	}
	if idx := bag.SearchPredicate(func(x *Cmp2Int) bool { return x.A >= 3 }); idx != 3 {
		T.Fatal("idx", idx, "bag", bag)
	}
	if idx := bag.SearchPredicate(func(x *Cmp2Int) bool { return x.A == 3 }); idx != 3 {
		T.Fatal("idx", idx, "bag", bag)
	}
	if idx := bag.SearchPredicate(func(x *Cmp2Int) bool { return x.A == 4 }); idx != -1 {
		T.Fatal("idx", idx, "bag", bag)
	}

	if idx := bag.SearchGreater(Cmp2Int{A: 1, B: 0}); idx != 0 {
		T.Fatal("idx", idx, "bag", bag)
	}
	if idx := bag.SearchGreater(Cmp2Int{A: 2, B: 0}); idx != 3 {
		T.Fatal("idx", idx, "bag", bag)
	}
	if idx := bag.SearchGreater(Cmp2Int{A: 3, B: 0}); idx != 3 {
		T.Fatal("idx", idx, "bag", bag)
	}
}

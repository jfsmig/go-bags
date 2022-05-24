// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"errors"
)

const (
	// MinSliceSize is the absolute minimal size of a slice return by the Slice methods.
	MinSliceSize = 1

	// MaxSliceSize is the absolute upper limit to the number of elements returned by the Slice methods.
	MaxSliceSize = 1000
)

var (
	errUnsorted   = errors.New("unsorted")
	errDuplicates = errors.New("duplicates")
)

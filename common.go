// Copyright (c) 2018-2022 Jean-Francois SMIGIELSKI
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bags

import (
	"errors"
)

const (
	minSliceSize = 1
	maxSliceSize = 1000
)

var (
	ErrUnsorted   = errors.New("unsorted")
	ErrDuplicates = errors.New("duplicates")
)

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"
)

func BenchmarkLoad(b *testing.B) {
	if err := generateLoad(b.N); err != nil {
		b.Errorf("generateLoad got err %v want nil", err)
	}
}

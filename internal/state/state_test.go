// Copyright © 2018 Jesse Lang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package state_test

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/jesselang/kcn/internal/state"
)

func TestStateClearWithoutPathFails(t *testing.T) {
	st := state.NewState(state.WithFs(afero.NewMemMapFs()))
	err := st.Clear()
	if err != state.ErrNoPath {
		t.Fatal(err)
	}
}

// kcn - when last context is empty
// kcn . - when last namespace is empty
// kcn . - when last namespace doesn't exist in current context
// kcn <contaxt> - when last namespace doesn't exist in <context>

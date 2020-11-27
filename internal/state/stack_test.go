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

package state

import (
	"encoding/json"
	"reflect"
	"testing"
)

func elementFixture() []Element {
	return []Element{
		{
			Context:   "alpha-dev",
			Namespace: "kube-system",
		},
		{
			Context:   "bravo-stage",
			Namespace: "default",
		},
		{
			Context:   "delta-prod",
			Namespace: "app-prod",
		},
	}
}

func TestEmptyStack(t *testing.T) {
	var s stack
	_, err := s.Pop()
	if err != errStackEmpty {
		t.Error("pop on empty stack should fail")
	}

	_, err = s.Peek()
	if err != errStackEmpty {
		t.Error("peek on empty stack should fail")
	}

	_, err = s.PeekPrev()
	if err != errStackEmpty {
		t.Error("peek prev on empty stack should fail")
	}

	l := s.Length()
	if l != 0 {
		t.Errorf("empty stack Length() should be zero, got %d", l)
	}
}

func TestStackOrdering(t *testing.T) {
	fixture := elementFixture()

	var s stack

	// insert fixture data in reverse order
	for i := len(fixture) - 1; i >= 0; i-- {
		s.Push(fixture[i])
	}

	// ensure that popped data matches fixture
	for _, v := range fixture {
		popped, err := s.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if *popped != v {
			t.Errorf("%v popped from stack does not match fixture%v",
				*popped, v)
		}
	}
}

func TestStackPeekXAndLength(t *testing.T) {
	fixture := elementFixture()

	var s stack

	// insert fixture data in reverse order
	for i := len(fixture) - 1; i >= 0; i-- {
		s.Push(fixture[i])
	}

	l := s.Length()
	if l != len(fixture) {
		t.Fatalf("stack length should match fixture length, got %d", l)
	}

	curr, err := s.Peek()
	if err != nil {
		t.Fatal(err)
	}

	if *curr != fixture[0] {
		t.Fatalf("stack peek should return first element, got %+v", *curr)
	}

	// stack peek should not change stack length or contents
	l = s.Length()
	if l != len(fixture) {
		t.Fatalf("stack length should match fixture length, got %d", l)
	}

	prev, err := s.PeekPrev()
	if err != nil {
		t.Fatal(err)
	}

	if *prev != fixture[1] {
		t.Fatalf("stack peek prev should return second element, got %+v", *prev)
	}

	// stack peek should not change stack length or contents
	l = s.Length()
	if l != len(fixture) {
		t.Fatalf("stack length should match fixture length, got %d", l)
	}

	// ensure that popped data matches fixture
	for _, v := range fixture {
		popped, err := s.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if *popped != v {
			t.Errorf("%v popped from stack does not match fixture%v",
				*popped, v)
		}
	}
}

func TestStackClear(t *testing.T) {
	fixture := elementFixture()

	var s stack

	// insert fixture data in reverse order
	for i := len(fixture) - 1; i >= 0; i-- {
		s.Push(fixture[i])
	}

	s.Clear()

	if s.Length() != 0 {
		t.Error("stack contains items after clearing")
	}
}

func TestStackSwap(t *testing.T) {
	input := elementFixture()
	expected := []Element{input[1], input[0], input[2]}

	var s stack

	// insert input data in reverse order
	for i := len(input) - 1; i >= 0; i-- {
		s.Push(input[i])
	}

	err := s.Swap()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range expected {
		popped, err := s.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if *popped != v {
			t.Errorf("%v popped from stack does not match expected %v",
				*popped, v)
		}
	}
}

func TestStackSwapInsufficient(t *testing.T) {
	input := elementFixture()[:1]

	var s stack

	// insert input data in reverse order
	for i := len(input) - 1; i >= 0; i-- {
		s.Push(input[i])

		err := s.Swap()
		if err != errStackInsufficient {
			t.Errorf("stack swap should fail due to insufficient elements %v", err)
		}
	}
}

func TestStackPeekPrevInsufficient(t *testing.T) {
	input := elementFixture()[:1]

	var s stack

	// insert input data in reverse order
	for i := len(input) - 1; i >= 0; i-- {
		s.Push(input[i])

		_, err := s.PeekPrev()
		if err != errStackInsufficient {
			t.Errorf("stack peek prev should fail due to insufficient elements %v", err)
		}
	}
}

func TestStackMarshaling(t *testing.T) {
	fixture := elementFixture()

	var s stack

	// insert fixture data in reverse order
	for i := len(fixture) - 1; i >= 0; i-- {
		s.Push(fixture[i])
	}

	fixture_json := []byte(
		`[{"context":"alpha-dev","namespace":"kube-system"},` +
			`{"context":"bravo-stage","namespace":"default"},` +
			`{"context":"delta-prod","namespace":"app-prod"}]`)

	b, err := json.Marshal(s)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b, fixture_json) {
		t.Errorf("%s does not match fixture %s", b, fixture_json)
	}
}

func TestStackUnmarshaling(t *testing.T) {
	fixture := elementFixture()

	fixture_json := []byte(
		`[{"context":"alpha-dev","namespace":"kube-system"},` +
			`{"context":"bravo-stage","namespace":"default"},` +
			`{"context":"delta-prod","namespace":"app-prod"}]`)

	var s stack
	err := json.Unmarshal(fixture_json, &s)
	if err != nil {
		t.Error(err)
	}

	// ensure that popped data matches fixture
	for _, v := range fixture {
		popped, err := s.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if *popped != v {
			t.Errorf("%v popped from stack does not match fixture%v",
				*popped, v)
		}
	}
}

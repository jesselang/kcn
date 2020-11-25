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
	"errors"
)

// element to be used in the stack
type Element struct {
	Context   string `json:"context"`
	Namespace string `json:"namespace"`
}

// basic stack that serializes to and from a JSON array
type stack struct {
	data []Element
}

func (s *stack) Length() int {
	return len(s.data)
}

func (s *stack) Clear() {
	s.data = nil
}

func (s *stack) Push(item Element) {
	s.data = append([]Element{item}, s.data...)
}

func (s *stack) Pop() (*Element, error) {
	if len(s.data) == 0 {
		return nil, errors.New("stack is empty")
	}

	popped := s.data[0]
	s.data = s.data[1:]

	return &popped, nil
}

func (s *stack) Peek() (*Element, error) {
	if len(s.data) == 0 {
		return nil, errors.New("stack is empty")
	}

	return &s.data[0], nil
}

func (s *stack) PeekPrev() (*Element, error) {
	if len(s.data) < 2 {
		return nil, errors.New("stack has less than 2 elements")
	}

	return &s.data[1], nil
}

func (s *stack) Swap() {
	if len(s.data) < 2 {
		return
	} else {
		swapped := s.data[0]
		s.data[0] = s.data[1]
		s.data[1] = swapped
	}
}

func (s stack) MarshalJSON() ([]byte, error) {
	jsonValue, err := json.Marshal(s.data)
	if err != nil {
		return nil, err
	}

	return jsonValue, nil
}

func (s *stack) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &s.data)
}

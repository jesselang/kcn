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
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/afero"

	"github.com/jesselang/kcn/internal/k8s"
)

var (
	ErrNoPath   = errors.New("no path set")
	ErrNotExist = errors.New("does not exist")
)

type State struct {
	Stack stack `json:"stack"`

	path string

	kc k8s.Client
	fs afero.Fs
}

type StateOption func(*State)

func WithPath(path string) StateOption {
	return func(s *State) {
		s.path = path
	}
}

func WithK8sClient(kc k8s.Client) StateOption {
	return func(s *State) {
		s.kc = kc
	}
}

func WithFs(fs afero.Fs) StateOption {
	return func(s *State) {
		s.fs = fs
	}
}

func NewState(opts ...StateOption) *State {
	s := &State{
		fs: afero.NewOsFs(),
		kc: k8s.NewKubectlClient(),
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (st *State) Read() error {
	if st.path == "" {
		return ErrNoPath
	}

	file, err := os.Open(st.path)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &st.Stack)
}

func (s State) Path() string {
	return s.path
}

func (s *State) Clear() error {
	if s.path == "" {
		return ErrNoPath
	}

	s.Stack.Clear()

	return s.Write()
}

func (s *State) Write() error {
	if s.path == "" {
		cache, err := os.UserCacheDir()
		if err != nil {
			return err
		}

		s.path = fmt.Sprintf("%s/kcn-%d-%s", cache, os.Getppid(), randString(6))
	}

	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Truncate(0)

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}

	n, err := file.Write(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return fmt.Errorf("could not fully write to state file")
	}

	return nil
}

func (st *State) Update(args ...string) error {
	if len(args) == 0 {
		return errors.New("no args given")
	}

	var context, namespace string
	if len(args) > 1 {
		namespace = args[1]
	}
	context = args[0]

	ctxList, err := st.kc.Contexts()
	if err != nil {
		return errors.New("could not get context list")
	}

	next := &Element{}

	if context == "." {
		if st.Stack.Length() == 0 {
			next.Context, err = st.kc.CurrentContext()
			if err != nil {
				return errors.New("could not get current context")
			}

			//next.Namespace = k8s.DefaultNamespace
		} else {
			curr, err := st.Stack.Peek()
			if err != nil {
				return err
			}
			next.Context = curr.Context
		}
	} else if context == "-" {
		if st.Stack.Length() == 0 {
			return errors.New("no previous state, try `kcn .`")
		} else {
			if len(namespace) == 0 {
				st.Stack.Swap()
				return st.Write()
			}

			next, err = st.Stack.Peek()
			if err != nil {
				return err
			}
		}
	} else {
		found := false
		for _, v := range ctxList {
			if context == v {
				next.Context = context
				found = true
			}
		}

		if !found {
			return fmt.Errorf("context %s not found", context)
		}
	}

	nsList, err := st.kc.Namespaces(next.Context)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"kcn: could not get namespace list for context %s,"+
				" falling back to %s\n",
			next.Context, k8s.DefaultNamespace)
		next.Namespace = k8s.DefaultNamespace
	}

	if len(namespace) == 0 || st.Stack.Length() == 0 {
		next.Namespace = k8s.DefaultNamespace
	} else {
		if namespace == "-" {
			prev, err := st.Stack.PeekPrev()
			if err != nil {
				return err
			}
			namespace = prev.Namespace
		}

		found := false
		for _, v := range nsList {
			if namespace == v {
				next.Namespace = namespace
				found = true
			}
		}

		if !found {
			return fmt.Errorf("namespace %s not found in context %s",
				namespace, next.Context)
		}
	}

	st.Stack.Push(*next)

	return st.Write()
}

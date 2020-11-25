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

	"github.com/jesselang/kcn/internal/kubectl"
)

const (
	EnvStatePath = "KCN_STATE_PATH"
)

type State struct {
	Stack stack `json:"stack"`

	path string
	k    kubectl.Kubectl
}

func NewState(k kubectl.Kubectl) (*State, error) {
	cache, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	initial := State{
		path: fmt.Sprintf("%s/kcn-%d-%s", cache, os.Getppid(), randString(6)),
	}

	if k == nil {
		initial.k = kubectl.NewKubectl()
	}

	return &initial, initial.Write()
}

func ReadState(path string) (*State, error) {
	if len(path) == 0 {
		return nil, errors.New("no state path given")
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var s State
	err = json.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}

	s.path = path
	s.k = kubectl.NewKubectl()
	return &s, nil
}

func (s *State) Path() string {
	return s.path
}

func (s *State) Clear() error {
	s.Stack.Clear()

	return s.Write()
}

func (s *State) Write() error {
	if len(s.path) == 0 {
		return fmt.Errorf("state path not set")
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

	ctxList, err := st.k.GetContextList()
	if err != nil {
		return errors.New("could not get context list")
	}

	next := &Element{}

	if context == "." {
		if st.Stack.Length() == 0 {
			next.Context, err = st.k.GetCurrentContext()
			if err != nil {
				return errors.New("could not get current context")
			}

			//next.Namespace = kubectl.DefaultNamespace
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

	nsList, err := st.k.GetNamespaceList(next.Context)
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"kcn: could not get namespace list for context %s,"+
				" falling back to %s\n",
			next.Context, kubectl.DefaultNamespace)
		next.Namespace = kubectl.DefaultNamespace
	}

	if len(namespace) == 0 || st.Stack.Length() == 0 {
		next.Namespace = kubectl.DefaultNamespace
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

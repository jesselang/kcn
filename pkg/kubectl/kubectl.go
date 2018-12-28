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

package kubectl

import (
	"os/exec"
	"strings"
)

const (
	DefaultNamespace = "default"
)

type Kubectl interface {
	GetContextList() ([]string, error)
	GetCurrentContext() (string, error)
	GetNamespaceList(context string) ([]string, error)
}

type Command struct{}

func NewKubectl() Kubectl {
	return &Command{}
}

func (k *Command) GetContextList() ([]string, error) {
	out, err := exec.Command("kubectl", "config", "get-contexts", "-o", "name").Output()
	return strings.Split(strings.TrimSpace(string(out)), "\n"), err
}

func (k *Command) GetCurrentContext() (string, error) {
	out, err := exec.Command("kubectl", "config", "current-context").Output()
	return strings.TrimSpace(string(out)), err
}

func (k *Command) GetNamespaceList(context string) ([]string, error) {
	out, err := exec.Command("kubectl", "--context", context,
		"get", "namespaces", "-o", "template",
		"--template={{range .items}}{{.metadata.name}} {{end}}").Output()
	return strings.Split(strings.TrimSpace(string(out)), " "), err
}

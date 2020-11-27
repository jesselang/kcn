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

type Stub struct {
	contexts       []string
	currentContext string
	namespaces     map[string][]string
}

func NewStub() Kubectl {
	return &Stub{
		contexts: []string{
			"alpha-dev",
			"bravo-stage",
			"delta-prod",
		},
		currentContext: "bravo-stage",
		namespaces: map[string][]string{
			"alpha-dev": []string{
				"app-a",
				"app-b",
				"app-c",
				DefaultNamespace,
				"kube-system",
			},
			"bravo-stage": []string{
				"app-d",
				"app-e",
				"app-f",
				DefaultNamespace,
				"kube-system",
			},
			"delta-prod": []string{
				"app-x",
				"app-y",
				"app-z",
				DefaultNamespace,
				"kube-system",
			},
		},
	}
}

func (k *Stub) Contexts() ([]string, error) {
	return k.contexts, nil
}

func (k *Stub) CurrentContext() (string, error) {
	return k.currentContext, nil
}

func (k *Stub) Namespaces(context string) ([]string, error) {
	if v, ok := k.namespaces[context]; ok {
		return v, nil
	}

	return []string{}, nil
}

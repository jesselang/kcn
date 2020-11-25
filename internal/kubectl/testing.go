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

type Mock struct {
	contextList    []string
	currentContext string
	namespaceList  map[string][]string
}

func NewMock() Kubectl {
	return &Mock{
		contextList: []string{
			"alpha-dev",
			"bravo-stage",
			"delta-prod",
		},
		currentContext: "bravo-stage",
		namespaceList: map[string][]string{
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

func (k *Mock) GetContextList() ([]string, error) {
	return k.contextList, nil
}

func (k *Mock) GetCurrentContext() (string, error) {
	return k.currentContext, nil
}

func (k *Mock) GetNamespaceList(context string) ([]string, error) {
	if v, ok := k.namespaceList[context]; ok {
		return v, nil
	} else {
		return []string{}, nil
	}
}

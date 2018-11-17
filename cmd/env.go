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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jesselang/kcn/pkg/state"
)

var (
	envInit bool
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		st, err := state.ReadState(os.Getenv(state.EnvStatePath))
		if !envInit {
			// This branch is likely being executed using shell's process
			// substitution. Non-zero exit codes won't propagate through
			// process substitution to the user's session.

			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}

			var curr *state.Element
			if st.Stack.Length() <= 0 {
				curr = &state.Element{}
			} else {
				curr, err = st.Stack.Peek()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s\n", err)
					os.Exit(1)
				}
			}

			fmt.Printf("%s=%s\n", envContext, curr.Context)
			fmt.Printf("%s=%s\n", envNamespace, curr.Namespace)
		} else {
			// XXX: won't work on non-bash shells or windows
			if err != nil {
				st, err = state.NewState(nil)

				if err != nil {
					fmt.Fprintf(os.Stderr, "error: %s\n", err)
					os.Exit(1)
				}

				fmt.Printf("export %s= %s=\n", envContext, envNamespace)
			}

			fmt.Printf("export %s=%s\n", state.EnvStatePath, st.Path())
			fmt.Println(
				`
kcn() {
	command kcn "$@"
	kcn_code=$?
	source <(command kcn env)
	[[ $kcn_code -eq 0 ]] || return $kcn_code
};`)
			// https://github.com/kubernetes/kubernetes/issues/27308#issuecomment-309207951
			fmt.Println(`alias kubectl="kubectl \
\${KCN_CONTEXT/[[:alnum:]-]*/--context=\${KCN_CONTEXT}} \
\${KCN_NAMESPACE/[[:alnum:]-]*/--namespace=\${KCN_NAMESPACE}}"`)
		}
	},
}

func init() {
	RootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	envCmd.Flags().BoolVarP(&envInit, "init", "i", false, "Initialize state (source from .*shrc)")

}

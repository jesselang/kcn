# kcn - Kubernetes Context and Namespace Switcher

Manage Kubernetes context and namespace within each shell session. Works with
zsh and probably bash.

## Installation

```
go install github.com/jesselang/kcn

# source kcn's environment from your shell's rc file
echo 'source <(kcn env --init)' >> $HOME/.zshrc

# use KCN_CONTEXT and KCN_NAMESPACE variables in your zsh theme
if [[ -n $KCN_CONTEXT || -n $KCN_NAMESPACE ]]; then
  prompt_segment cyan black "☸ $KCN_CONTEXT/$KCN_NAMESPACE"
fi

# source kcn's environment to your existing shell session
source <(kcn env --init)

# select a context and namespace
# kcn <context> [ <namespace> ]

# select alpha-dev context, default namespace
kcn alpha-dev

# same context, different namespace
kcn . kube-system

# different context, different namespace
kcn bravo-stage app-stage

# return to previous context and namespace, alpha-dev and kube-system
kcn -

# different context, different namespace
kcn delta-prod app-prod

# same context, previous namespace, delta-prod and kube-system
kcn . -

# clear context and namespace for this session
kcn clear
```

## Building

Requires golang 1.11.

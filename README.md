# kubectl ver

![test](https://github.com/dty1er/kubectl-ver/workflows/test/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/dty1er/kubectl-ver)](https://goreportcard.com/report/github.com/dty1er/kubectl-ver)

kubectl plugin for better "kubectl version"

![skew-s](https://user-images.githubusercontent.com/60682957/105186576-35fce700-5b75-11eb-8e3a-cd4c0fdabff4.png)

## What's this?

With `kubectl-ver` , you can check if your kubernetes usage meets the __version skew policy__.

In kubernetes, [version skew policy](https://kubernetes.io/docs/setup/release/version-skew-policy/) is a bit confusing, especially beginners.
It is important to understand you are always following the policy because using unsupported cluster/kubectl is problematic and even dangerous.
To know if your kubernetes usage is met with it, you need to know the cluster version, client version, and current latest version. Of course, you have to understand the detail of the policy.
`kubectl ver skew` command helps this situation. When you run it, it automatically fetches the cluster, client, and latest version and judges if it's following the policy.
By using this, it will be easy for you to understand it your kubernetes usage is met the policy.

* Check if kubectl update is available

Because kubernetes development is fast, it is not easy to follow the latest kubectl.
`kubectl ver check` cheks if kubectl update is available.


Also, `kubectl-ver` can do all "kubectl version" can do. You can use `kubectl ver` command as an complete alternative to `kubectl version`.

## Installation

Currently "krew" kubectl plugin manager installation is unsupported.

### Manually via go get

```sh
go install github.com/dty1er/kubecolor/cmd/kubecolor
```

## Usage

`kubectl ver` is designed to be an alternative to `kubectl version`.
It means, you can pass any options which `kubectl version` accepts are accepted by `kubectl ver`.

```shell
# all of below works, because it works with `kubectl version`
$ kubectl ver
$ kubectl ver --short
$ kubectl ver --client
$ kubectl ver -o json
```

Additionally, `kubectl ver` accepts some subcommands to show more information about kubernetes cluster and kubectl versions.

### Supported commands

* `check`

`kubectl ver check` is a subcommand to check if there is the kubectl available update.
When update is available, the message is shown.
![check](https://user-images.githubusercontent.com/60682957/105186764-60e73b00-5b75-11eb-86a9-98c55743ea5c.png)

* `skew`

`kubectl ver skew` is a subcommand to show if there is the kubernetes cluster and kubectl versions skew.

* cluster version problem

![skew-s](https://user-images.githubusercontent.com/60682957/105186576-35fce700-5b75-11eb-8e3a-cd4c0fdabff4.png)

* kubectl version problem

![skew-c](https://user-images.githubusercontent.com/60682957/105186580-36957d80-5b75-11eb-8686-44742c0605b9.png)

* following version skew policy

![skew](https://user-images.githubusercontent.com/60682957/105186753-5d53b400-5b75-11eb-9f6c-5149511e3bbf.png)


## Contributions

Always welcome. Just opening an issue should be also greatful.

## LICENSE

MIT


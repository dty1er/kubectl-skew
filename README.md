# kubectl skew

![test](https://github.com/dty1er/kubectl-skew/workflows/test/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/dty1er/kubectl-skew)](https://goreportcard.com/report/github.com/dty1er/kubectl-skew)

A simple kubectl plugin to show if your kubernetes/kubectl version is "skewed"

![skew-s](https://user-images.githubusercontent.com/60682957/105196269-cb50a900-5b7e-11eb-9505-4d0f14a4ca84.png)

## What's this?

With `kubectl skew` , you can check if your kubernetes usage meets the __version skew policy__.

In kubernetes, [version skew policy](https://kubernetes.io/docs/setup/release/version-skew-policy/) is a bit confusing, especially for beginners.<br>
However, it is important to make sure you are always following the policy because using unsupported cluster/kubectl is problematic and even dangerous.<br>
In order to know if your kubernetes usage is met with it, you need to know the cluster version, client version, and current latest version. Of course, you have to understand the detail of the policy.<br>
`kubectl ver skew` command helps this situation. When you run it, it automatically fetches the cluster, client, and latest version and judges if it's following the policy.<br>
By using `kubectl skew`, it will be easy for you to understand if your kubernetes usage meets the policy.

## Installation

You can install `kubectl-skew` by [krew](https://github.com/kubernetes-sigs/krew) (kubectl plugin manager).
Run below to install on your machine.

```sh
kubectl krew install skew
```

To use `krew`, first you need to install it. Follow the [krew installation guide](https://krew.sigs.k8s.io/docs/user-guide/setup/install/).

## Usage

You simply need to run `kubectl skew`, which shows if there is the kubernetes cluster and kubectl versions skew.

* cluster version problem

![skew-s](https://user-images.githubusercontent.com/60682957/105196269-cb50a900-5b7e-11eb-9505-4d0f14a4ca84.png)

* kubectl version problem

![skew-c](https://user-images.githubusercontent.com/60682957/105197817-5d0ce600-5b80-11eb-8505-f47afad7dad3.png)

* following version skew policy

![skew](https://user-images.githubusercontent.com/60682957/105196273-cc81d600-5b7e-11eb-99d9-31ef0213b9bb.png)

## Upcoming releases

* Support output option (e.g. `-o json`)

## Contributions

Always welcome. Just opening an issue should be also greatful.

## LICENSE

MIT


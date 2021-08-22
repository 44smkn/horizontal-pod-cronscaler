# Horizontal Pod Cronscaler

[![Go Reference](https://pkg.go.dev/badge/github.com/44smkn/horizontal-pod-cronscaler.svg)](https://pkg.go.dev/github.com/44smkn/horizontal-pod-cronscaler)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

`Horizontal Pod Cronscaler` scales the workload based on the specified Cron Schedule.

## Installation

Latest `Horizontal Pod Cronscaler` release can be installed by running:

```sh
kubectl apply -f https://github.com/44smkn/horizontal-pod-cronscaler/latest/download/components.yaml
```

Installation instructions for previous releases can be found in [Horizontal Pod Cronscaler releases](https://github.com/44smkn/horizontal-pod-cronscaler/releases).

## Example

```yaml
apiVersion: autoscaling.44smkn.github.io/v1beta1
kind: HorizontalPodCronscaler
metadata:
  name: horizontalpodcronscaler-sample
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: nginx-deployment
  replicas: 2
  schedule: "*/3 * * * *"
```

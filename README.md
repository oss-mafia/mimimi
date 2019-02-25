# Mimimi Slack bot

A Slack bot that mimimizes all messages sent to the channels where the bot is present.

![mimimi](mimimi.jpg)

# Build

```bash
$ make help

Usage:
  make <target>

Build targets
  build            Build the bot
  build-static     Build the statically linked Linux binary
  clean            Clean all binary articats

Packaging and distribution
  docker-build     Build the Docker image
  docker-push      Push the Docker image to the configured registry

Others
  help             Display this help
```

# Deploy

Convenience service and deployment manifests are provided to deploy the bot in a Kubernetes that has
[Istio](https://istio.io/) installed.

First of all update the `config/mimimi-secrets.yaml` file with appropriate values, and then apply all the provided
manifests:

```bash
kubectl create namespace mimimi
kubectl label namespace mimimi istio-injection=enabled

kubectl apply -n mimimi -f config/  
```

# License

This software is licensed under the Apache License 2.0. See [LICENSE](LICENSE) file for details.

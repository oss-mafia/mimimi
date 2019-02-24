# Mimimi Slack bot

A Slack bot that mimimizes all messages sent to the channels where the bot is present.

![mimimi](mimimi.jpg)

# Building

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
  docker-run       Run the Docker image

Others
  help             Display this help
```

# License

This software is licensed under the Apache License 2.0. See [LICENSE](LICENSE) file for details.

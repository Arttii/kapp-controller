# kapp-controller

- Slack: [#k14s in Kubernetes slack](https://slack.kubernetes.io)
- [Docs](docs/README.md) with topics about config, etc.
- Install: see below section

kapp controller provides a way to specify which applications should run on your K8s cluster. It will install, and continiously apply updates.

Features:
- supports fetching Helm charts, git repos, Docker images, inline content within resource
- supports templating of Helm Charts, ytt configuration
- installs and syncs resources with [kapp](https://get-kapp.io)

## Development & Deploy

Install ytt, kbld, kapp beforehand (https://k14s.io).

```
./hack/build.sh  # to build locally
./hack/deploy.sh # to deploy
```
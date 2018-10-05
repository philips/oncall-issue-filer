# On-Call Issue Filer

[![Docker Repository on Quay](https://quay.io/repository/philips/oncall-issue-filer/status "Docker Repository on Quay")](https://quay.io/repository/philips/oncall-issue-filer)

File issues on GitHub based on acknowledged alerts.

On-call rotations where an incident requires multiple stakeholders coming
together over a period of days or weeks need a coordination point. For many
teams that is GitHub Issues.

This tool will assist an on-call person coordinate by automatically filing a
GitHub issue once an alert has been acknowledged.

| Tool       | Supported?  |
| ---------- |------------ |
| OpsGenie   | Yes         |
| Pager Duty | Help Wanted |
| ???        | Help Wanted |

## Kubernetes Usage

Deploy the CronJob:

```
git clone https://github.com/philips/oncall-issue-filer
kubectl create -f kube/cronjob.yaml
```

Add your secrets based on `kube/secret.yaml`. Keep in mind the [fields are base64 encoded](https://kubernetes.io/docs/concepts/configuration/secret/).

## Hacking

Make sure you have Go 1.11+ installed then:

```
git clone https://github.com/philips/oncall-issue-filer
export GO111MODULE=on
go run main.go
```

## TODO

- Refactor into cleaner packages
- Introduce dry-run via refactor

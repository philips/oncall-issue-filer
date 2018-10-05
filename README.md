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

## OpsGenie Configuration

### User Tag Configuration

Every user must add a `github=GITHUB_USERNAME` tag to their OpsGenie profile for the tool to work.

![User Configuration](/docs/opsgenie/user-tag.png)

### Alert Configuration

**Ignore Replies**

Add Ignore Filter for `(?:(?i)re: |fwd: )?(.*)`.

![Ignore Filter Configuration](/docs/opsgenie/ignore-filter.png)

**Alert Fields**

Configure as shown below.

![Alert Fields Configuration](/docs/opsgenie/alert-fields.png)

## TODO

- Refactor into cleaner packages
- Introduce dry-run via refactor

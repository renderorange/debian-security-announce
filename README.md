# debian-security-announce

A CLI tool that fetches Debian Security Advisories from the official RSS feed and posts them to a Slack channel via incoming webhook.

## Build

```bash
go build -o dsa-slack ./cmd/dsa-slack/
```

## Usage

```bash
DSA_SLACK_WEBHOOK_URL=https://hooks.slack.com/services/... ./dsa-slack
```

## Configuration

| Environment Variable | Required | Default | Description |
|---------------------|----------|---------|-------------|
| `DSA_SLACK_WEBHOOK_URL` | Yes | - | Slack incoming webhook URL |
| `DSA_STATE_FILE` | No | `~/.config/debian-security-announce/state.json` | Path to state file |
| `DSA_RSS_FEED_URL` | No | `https://www.debian.org/security/dsa-long` | RSS feed URL |

## Cron

```bash
*/15 * * * * DSA_SLACK_WEBHOOK_URL=https://hooks.slack.com/services/... /usr/local/bin/dsa-slack
```

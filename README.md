# GOTRACK

gotrack - is a little helper to track time in [YouTrack](https://www.jetbrains.com/youtrack/)

## Usage

Create timetrack table in format by `example/format.csv` via tables

Create access token in `Account/Security`

Create config in `${HOME}/.gotrack/config`

```json
{
    "youTrackHost": "https://youtrack.instance",
    "token": "perm:token"
}
```

```shell
> gotrack track -tt timetrack.csv --dry-run # run in dry-run
2023/07/29 12:15:22 Tracked time for MEMES-2, spent time 2h15m0s, dry-run
2023/07/29 12:15:22 TimeTrack record for Issue MEMES-2 marked as 'tracked', skip
2023/07/29 12:15:22 Tracked time for MEMES-4, spent time 1h10m0s, dry-run
2023/07/29 12:15:22 Tracked time for MEMES-3, spent time 30m0s, dry-run
2023/07/29 12:15:23 Tracked time for MEMES-2, spent time 1h0m0s, dry-run
2023/07/29 12:15:23 Total tracked time: 4h55m0s

> gotrack track -tt timetrack.csv
2023/07/29 12:14:56 Tracked time for MEMES-2, spent time 2h15m0s
2023/07/29 12:14:56 TimeTrack record for Issue MEMES-2 marked as 'tracked', skip
2023/07/29 12:14:56 Tracked time for MEMES-4, spent time 1h10m0s
2023/07/29 12:14:57 Tracked time for MEMES-3, spent time 30m0s
2023/07/29 12:14:58 Tracked time for MEMES-2, spent time 1h0m0s
2023/07/29 12:14:58 Total tracked time: 4h55m0s
```

## Build

```shell
go build -v -o ./bin/gotrack ./cmd/gotrack
```
# Speed test metrics

This small tool is useful for process automatic speedtest and scrape it from other tool.

## Usage

### Run in standalone

```sh
speed-test
```

**For customize test interval you can simply specify CRON expression with `-interval`. Eg.**

```sh
speed-test -interval="@hourly"
```

### Run in Docker

#### Without arguments

```sh
docker run orblazer/speed-test:latest
```

#### With arguments

```sh
docker run --entrypoint "speed-test -interval=\"@hourly\"" orblazer/speed-test:latest
```

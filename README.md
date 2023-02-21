# WheelOfFortuneEveryday

This is a Go REST API that connects to your locally running Roku TV and runs a specific set of instructions.

I use this to turn on my TV, launch the ABC App, and select "Live TV" so I can watch Jeopardy and Wheel of Fortune every day that it's on.

I use the iOS Shortcuts app to schedule a request be sent to the locally running API everyday at 7pm.

## Startup script
```sh
#!/bin/sh
export DEVICE_NAME="enter your TV's device name here"
export TV_ADDRESS="http://192.168.1.18:8060/" # your TV's local IP address here, must be 8060 bc that's what roku expects
go run main.go
```

The REST API is exposed on port 8787.

Outputs logs to a `wofed.log` file

### Examples

```sh
curl localhost:8787/healthcheck
```

```sh
curl -X POST localhost:8787/api/v1/jobs?force=true -H Content-Type:application/json -d '{"appName": "ABC"}'
```


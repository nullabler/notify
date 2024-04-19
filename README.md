## Initialization
Create telegram bot - [@BotFather](https://t.me/BotFather)

Create config/dev.yaml looks like config/example.yaml and add telegram token

## Start
```bash
make up build=1 watch=1
```

## Setting
For setting notification you need asking your bot:
```
>>> ChatID
```
Add you ChatID in config/dev.yaml for templateToChats block

## Test notification
```bash
curl -X POST http://localhost:8081/send/pipeline-stage -d '{"username": "MyUsername", "state": "Start", "build-number": "123"}'
```

# TG Notify Service  

## About
This is a simple notification service (not a daemon) that sends notifications to Telegram at the end of a running command.  

The service reads its configuration from the `NOTIFY_CONFIG` environment variable. If the variable is not set, the service stops and does not execute any commands. The command's output is redirected to stdout.  

The configuration must include:  
```
chat_id: CHAT_ID
tg_bot_token: TG_BOT_TOKEN
instance: YOUR_INSTANCE
thread_id: THREAD_ID
```
where:  
- `chat_id` is the ID of the chat where the message will be sent.  
- `tg_bot_token` is the Telegram bot token.  
- `instance` is the instance where the service is running.  
- `thread_id` is the thread ID in a Telegram group. It can be set to `-1`, meaning there is no thread.  

**IMPORTANT:**  
- If sending to a forum chat, the `chat_id` in the configuration must start with `-100` followed by the chat ID.  
- If sending to a group chat, the `chat_id` in the configuration should match the Telegram chat ID.  

## How to Run

### Bare Metal
To run it on bare metal, first compile the binary:

```sh
go mod download
go build -o tgnotify
```

Then run it with:

```sh
./tgnotify ping google.com
```

### Docker Compose
You can also run it using Docker Compose:

```sh
docker compose up -d
```

Before doing that, make sure to specify the command you want to run in your `docker-compose.yml` file:

```yml
command: ["ping", "google.com"]
```

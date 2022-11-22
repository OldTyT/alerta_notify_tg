[![Go Report Card](https://goreportcard.com/badge/github.com/OldTyT/alerta_notify_tg)](https://goreportcard.com/report/OldTyT/alerta_notify_tg)
[![GolangCI](https://golangci.com/badges/github.com/OldTyT/alerta_notify_tg.svg)](https://golangci.com/r/github.com/OldTyT/alerta_notify_tg)


# Alerta notify TG

---

This is simple desktop notify on alerts from [Alerta](https://github.com/alerta/alerta)

Default path to config - `$HOME/.config/alerta_notify_tg.json`

### How to start

```
./alerta_notify_tg -config="path/to/config"
```

### Config

* `alerta_username` - `STRING` - User name from [Alerta](https://github.com/alerta/alerta)
* `alerta_password` - `STRING` -User password from [Alerta](https://github.com/alerta/alerta)
* `alerta_url` - `STRING` - URL Address [Alerta](https://github.com/alerta/alerta)
* `alert_query` - `STRING` - Request to [Alerta](https://github.com/alerta/alerta) by which alerts will be received
* `time_sleep` - `INT` - Sleep time between iterations, in seconds
* `telegram_token` - `STRING` - Token Telegram bot
* `telegram_chat` - `INT` - Telegram chat id

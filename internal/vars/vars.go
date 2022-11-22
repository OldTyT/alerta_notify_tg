package vars

type NotifierCFG struct {
	AlertaUsername string `json:"alerta_username"`
	AlertaPassword string `json:"alerta_password"`
	AlertaURL      string `json:"alerta_url"`
	AlertaQuery    string `json:"alert_query"`
	TimeSleep      int    `json:"time_sleep"`
	TGToken        string `json:"telegram_token"`
	TGChat         int    `json:"telegram_chat"`
}

type OtherCFG struct {
	AlertaToken string
}

var (
	Notifier NotifierCFG
	Other    OtherCFG
	Version  string = "v0.0.2"
)

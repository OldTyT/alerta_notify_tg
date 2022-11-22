package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/OldTyT/alerta_notify/internal/log"
	"github.com/OldTyT/alerta_notify/internal/vars"
)

func main() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	ConfFile := flag.String("config", homedir+"/.config/alerta_notify_tg.json", "Path to conf file.")
	flag.Parse()
	file, err := os.Open(*ConfFile)
	if err != nil {
		ErrorMsg := err.Error()
		fmt.Println(ErrorMsg)
		log.Error.Println(ErrorMsg)
		ErrorMsg = "Error open conf file -" + *ConfFile
		fmt.Println(ErrorMsg)
		log.Error.Println(ErrorMsg)
		os.Exit(1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&vars.Notifier)
	if err != nil {
		ErrorMsg := "error:" + err.Error()
		fmt.Println(ErrorMsg)
		log.Error.Println(ErrorMsg)
		os.Exit(1)
	}
	log.LogFunc()
	go SendSummary()
	LoginAlerta()
	for {
		go UpdateAlerts()
		time.Sleep(time.Duration(vars.Notifier.TimeSleep) * time.Second)
	}
}

func SendMessage(text string) {
	URL, err := url.Parse("https://api.telegram.org/bot" + vars.Notifier.TGToken + "/sendMessage")
	if err != nil {
		ErrorExiting(err.Error())
	}
	query := URL.Query()
	query.Add("chat_id", strconv.Itoa(vars.Notifier.TGChat))
	query.Add("text", text)
	URL.RawQuery = query.Encode()
	resp, err := http.Get(URL.String())
	if err != nil {
		ErrorExiting("Error send message in telegram. Error: " + err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ErrorExiting("Response status != 200, when send message in telegram.\nExit.")
	}
}

func SendSummary() {
	URL := vars.Notifier.AlertaURL + vars.Notifier.AlertaQuery
	SendMessage("Alerta query: " + URL + "\nSleep time: " + strconv.Itoa(vars.Notifier.TimeSleep) + "sec" + "\nVersion: " + vars.Version)
}

func ErrorExiting(ErrorMsg string) {
	log.Warn.Println(ErrorMsg)
	log.Error.Println(ErrorMsg)
	go SendMessage(ErrorMsg)
	go SendMessage("Exiting")
	os.Exit(1)
}

func LoginAlerta() {
	type AlertaAuth struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	type AlertaToken struct {
		Token string `json:"token"`
	}
	var (
		AlertaAuthData AlertaAuth
		TokenLocal     AlertaToken
	)
	AlertaAuthData.Password = vars.Notifier.AlertaPassword
	AlertaAuthData.UserName = vars.Notifier.AlertaUsername
	JsonData, err := json.Marshal(AlertaAuthData)
	if err != nil {
		fmt.Println(err)
	}
	URL := vars.Notifier.AlertaURL + "/auth/login"
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(JsonData))
	if err != nil {
		ErrorExiting("Error auth in alerta. " + err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		ErrorExiting("Response status != 200, when authorization in alerta.\nExit.")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorExiting("Can't read JSON: " + err.Error())
	}
	if err := json.Unmarshal(body, &TokenLocal); err != nil {
		ErrorExiting("Can't unmarshal JSON: " + err.Error())
	}
	vars.Other.AlertaToken = TokenLocal.Token
	go SendMessage("Successful authorization in Alerta")
}

func UpdateAlerts() {
	type AlertaAlertList struct {
		Alerts []map[string]interface{} `json:"alerts"`
		Total  int                      `json:"total"`
	}
	type AlertSummary struct {
		AlertName string `json:"event"`
		Resource  string `json:"resource"`
		ENV       string `json:"environment"`
		Severity  string `json:"severity"`
		ID        string `json:"id"`
	}
	var (
		AlertsSummary AlertaAlertList
		Alert         []AlertSummary
	)
	URL := vars.Notifier.AlertaURL + vars.Notifier.AlertaQuery
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		go SendMessage("Error when receiving alerts.")
	}
	req.Header.Set("Authorization", "Bearer "+vars.Other.AlertaToken)
	resp, err := client.Do(req)
	if err != nil {
		go SendMessage("Error when receiving alerts.")
	}
	if resp.StatusCode != 200 {
		go SendMessage("Response status != 200, when update alerts.")
	} else {
		// defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			go SendMessage("Can't read JSON, when update alerts. " + err.Error())
		}
		if err = json.Unmarshal(body, &AlertsSummary); err != nil {
			go SendMessage("Can't unmarshal JSON, when update alerts. " + err.Error())
		}
		if AlertsSummary.Total != 0 {
			b, err := json.Marshal(AlertsSummary.Alerts)
			if err != nil {
				go SendMessage("Can't marshal JSON. " + err.Error())
			}
			if err = json.Unmarshal(b, &Alert); err != nil {
				go SendMessage("Can't unmarshal JSON. " + err.Error())
			}
			for key := range Alert {
				AlertTxt := "................................................................................\nAlert name: " + Alert[key].AlertName + "\nResource: " + Alert[key].Resource + "\nENV: " + Alert[key].ENV + "\nSeverity: " + Alert[key].Severity + "\nURL: " + vars.Notifier.AlertaURL + "/#/alert/" + Alert[key].ID + " \n................................................................................"
				go SendMessage(AlertTxt)
			}
		}
	}
}

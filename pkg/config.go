package pkg

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Log struct {
	MaxSize   int  `json:"max_size"`
	MaxBackup int  `json:"max_backup"`
	Debug     bool `json:"debug"`
}

type Database struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	User    string `json:"user"`
	Pass    string `json:"password"`
	Name    string `json:"name"`
	Ssl     string `json:"ssl"`
	Timeout int    `json:"timeout"`
}

type Queue struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Name    string `json:"name"`
	Expired int    `json:"expired"`
}

type Chatbot struct {
	Host              string `json:"host"`
	HostPush          string `json:"host_push"`
	PhoneID           string `json:"phone_id"`
	AccessToken       string `json:"access_token"`
	WabaAccountNumber string `json:"waba_account_number"`
	AccessTokenPush   string `json:"access_token_push"`
	DivisionID        string `json:"division_id"`
	AccountID         string `json:"account_id"`
}

type Service struct {
	Log            Log      `json:"log"`
	Database       Database `json:"database"`
	Queue          Queue    `json:"queue"`
	Chatbot        Chatbot  `json:"chatbot"`
	DownloadFolder string   `json:"download_folder"`
}

func LoadServiceConfig(configFilePath string) (cfg *Service, err error) {

	if len(configFilePath) == 0 {
		err = errors.New("can't load config file")
		return
	}

	cfg, err = loadConfigFile(configFilePath)

	return
}

func loadConfigFile(f string) (c *Service, err error) {

	var content []byte
	content, err = ioutil.ReadFile(f)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &c)
	return

}

package info

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	Version = "0.12"
)

type version struct {
	LatestVersion string `json:"Message"`
	Error         string `json:"Err"`
	Info          info   `json:"Info"`
}
type info struct {
	Detail string `json:"details"`
	Date   string `json:"date"`
}

func CheckVersion() string {
	u := "https://hmwri.com:8080/latestVersion"
	request, _ := http.NewRequest("GET", u, nil)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{ServerName: "hmwri.com", InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(request)
	if err != nil {
		return fmt.Sprintln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintln(err)
	}
	var versionStruct version
	if err := json.Unmarshal(body, &versionStruct); err != nil {
		return fmt.Sprintln(err)
	}
	if versionStruct.LatestVersion != Version {
		return fmt.Sprintln("\n注意!最新バージョン(Ver"+versionStruct.LatestVersion+")があります。\n(https:hmwri.com/peridot/download.html)で最新版をダウンロードしてください") +
			fmt.Sprintf(">>最新版:Version %s (%s) \n %s \n", versionStruct.LatestVersion, versionStruct.Info.Date, versionStruct.Info.Detail)
	}
	return "(最新バージョン)\n"
}

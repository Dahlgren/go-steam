package depot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Server struct {
	CellID                 int64   `json:"cell_id"`
	Host                   string  `json:"host"`
	HTTPSSupport           string  `json:"https_support"`
	Load                   int64   `json:"load"`
	NumEntriesInClientList int64   `json:"num_entries_in_client_list"`
	PreferredServer        bool    `json:"preferred_server"`
	SourceID               int64   `json:"source_id"`
	Type                   string  `json:"type"`
	Vhost                  string  `json:"vhost"`
	WeightedLoad           float64 `json:"weighted_load"`
}

type steamPipeServersResponse struct {
	Response struct {
		Servers []Server `json:"servers"`
	} `json:"response"`
}

func GetContentServers(cellId int64) ([]Server, error) {
	getServerInfo := NewSteamMethod("IContentServerDirectoryService", "GetServersForSteamPipe", 1)

	vals := url.Values{}
	vals.Add("cell_id", strconv.FormatInt(cellId, 10))
	//vals.Add("max_servers", strconv.FormatInt(maxServers, 10))

	var resp steamPipeServersResponse
	err := getServerInfo.Request(vals, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Response.Servers, nil
}

// BaseSteamAPIURLProduction is the steam url used to do requests in prod
const BaseSteamAPIURLProduction = "https://api.steampowered.com"

// BaseSteamAPIURL is the url used to do requests, defaulted to prod
var BaseSteamAPIURL = BaseSteamAPIURLProduction

// A SteamMethod represents a Steam Web API method.
type SteamMethod string

// NewSteamMethod creates a new SteamMethod.
func NewSteamMethod(interf, method string, version int) SteamMethod {
	m := fmt.Sprintf("%v/%v/%v/v%v/", BaseSteamAPIURL, interf, method, strconv.Itoa(version))
	return SteamMethod(m)
}

func (s SteamMethod) Request(data url.Values, v interface{}) error {
	url := string(s)
	if data != nil {
		url += "?" + data.Encode()
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("steamapi %s Status code %d", s, resp.StatusCode)
	}

	d := json.NewDecoder(resp.Body)

	return d.Decode(&v)
}

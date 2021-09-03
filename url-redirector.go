package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"net/http"
)

type Settings struct {
	Global Global	`json:"global"`
	Sites json.RawMessage `json:"redirects"`
	Aliases json.RawMessage `json:"aliases"`
}

type Global struct {
	Port string
}

type Site struct {
	FullURL string
	Description string
}

func readSettings() ( Global, map[string]*Site, error ) {
	settingsFile, err := ioutil.ReadFile("settings.json")
	var globalSettings Global
	var settings Settings
	sites := make( map[string]*Site )
	var tempSites map[string]interface{}
	var tempAliases map[string]interface{}

	if err != nil {
		return globalSettings, sites, err
	}
	json.Unmarshal( settingsFile, &settings )
	port :=  strings.Trim( settings.Global.Port , ":" )
	globalSettings.Port = port

	json.Unmarshal( settings.Sites, &tempSites )

	for name, value := range tempSites{
		temp := value.( map[ string ]interface{})
		sites[ name ] = &Site{
			FullURL: temp["fullURL"].(string),
			Description: temp["Description"].(string)}
	}

	json.Unmarshal(settings.Aliases, &tempAliases )
	for name, value := range tempAliases{
		sites[ name ] = sites[value.(string)]
	}
	return globalSettings, sites, nil
}

func redirect( w http.ResponseWriter, r * http.Request ){
	_, sites, err := readSettings()
	if err == nil {
		path := strings.Trim( r.URL.Path, "/")
		value, ok := sites[path]
		if ok {
			http.Redirect(w, r, value.FullURL, http.StatusTemporaryRedirect)
		} else {
			// display our own page
			fmt.Fprintln( w, "Uh oh" )
		}
	} 
}

func ignore( w http.ResponseWriter, r * http.Request ){

}

func main() {
	globalSettings, _, err := readSettings()
	if err == nil {
		http.HandleFunc( "/favicon.ico", ignore )
		http.HandleFunc( "/", redirect )
		err = http.ListenAndServe(":" + globalSettings.Port, nil )	
	} else {
		fmt.Println(err)
	}
}
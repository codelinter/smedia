package main

import (
	"context"
	"fmt"
	"net/http"

	fb "github.com/huandu/facebook"
	"github.com/therecipe/qt/core"
	"golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook"
)

var srvr *http.Server

func handler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	codeChan <- vals.Get("code")
	fmt.Fprint(w, "Please go back to the app")
}

func localhost(cert, key string) {
	srvr = &http.Server{Addr: "localhost:12345"}
	http.HandleFunc("/", handler)
	err := srvr.ListenAndServeTLS(cert, key)
	if err != nil {
		e := fmt.Sprintf("Server startup ERR ", err.Error())
		mainWindow.StatusBar().ShowMessage(core.QDir_ToNativeSeparators(e), 0)
	}
}

var codeChan = make(chan string)

// startServerAndRegister return Perm Token and PageID
func startServer(cert, key string) (string, string, error) {
	if clientID == "" || clientSecret == "" {
		return "", "", fmt.Errorf("ClientID and ClientSecret required")
	}
	println("CLIENT", clientID, clientSecret)
	//return "", "", fmt.Errorf("CUSTOM")
	go localhost(cert, key)
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://localhost:12345/",
		Scopes:       []string{"manage_pages,publish_pages"},
		Endpoint:     oauth2fb.Endpoint,
	}

	url := conf.AuthCodeURL("deb6bfg1b75f0b1", oauth2.AccessTypeOnline)
	//fmt.Printf("Visit the URL for the auth dialog: \n%v\n", url)
	goToURLLineEdit.SetText(url)
	var code string
	code = <-codeChan

	ephemralToken, err := conf.Exchange(ctx, code)
	if err != nil {
		return "", "", fmt.Errorf("Code Exchange ERR %v", err)
	}

	session := &fb.Session{}
	session.SetAccessToken(ephemralToken.AccessToken)
	session.Version = "v3.0"

	res, err := session.Get("/oauth/access_token", fb.Params{
		"grant_type":        "fb_exchange_token",
		"client_id":         clientID,
		"client_secret":     clientSecret,
		"fb_exchange_token": ephemralToken.AccessToken,
	})
	if err != nil {
		return "", "", fmt.Errorf("fb_exchange_token -> %v", err)
	}

	var pToken interface{}
	var ok bool

	if pToken, ok = res["access_token"]; !ok {
		return "", "", fmt.Errorf(" Exchange access_token' not found")
	}
	permMeToken := pToken.(string)
	pid, err := getPageID(permMeToken)
	if err != nil {
		return "", "", err
	}
	pageID := "/" + pid
	res, err = session.Get(pageID, fb.Params{
		"fields":       "access_token",
		"access_token": permMeToken,
	})
	if pToken, ok = res["access_token"]; !ok {
		return "", "", fmt.Errorf("Permanent access_token not found")
	}

	return pToken.(string), pid, nil
}

func getPageID(permMeToken string) (string, error) {
	//permMeToken = `EAACj7ixscaABAAHMrhoAdJavH5V3BdKBYAqBCsNlt78QpiyllDOLvuRL80gGY6NZBNAXB6TYjZCLAqZBLG7rydVGX8sq4ZBZCeLQelUnI5FWB7Cezd5KEm81NcUqPzzT68JpdpYr19mfC73IHd9WcZBK8CuAF30SZCnQXJDm5ybMbL5J6ajDG1B`
	session := &fb.Session{}
	session.SetAccessToken(permMeToken)
	session.Version = "v3.0"
	uri := "/" + clientFacebookPageURL
	res, err := session.Get(uri, fb.Params{
		"fields":       "id",
		"access_token": permMeToken,
	})
	if err != nil {
		return "", err
	}
	return res["id"].(string), nil
}

func getCredsFromFile() {}

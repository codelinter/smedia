package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const (
	certFileName       = "example.com.pem"
	keyFileName        = "example.com-key.pem"
	fbIDSecretFileName = "fbcreds"
)

var (
	clientFacebookPageURL, gotoURL string
	certificateDir                 string //  directory to find the cert,key pair

	// https://developers.facebook.com/apps/<your-app-id>/settings/basic/
	// clientID and clientSecret are aliases for App ID and App Secret
	clientID     string
	clientSecret string
)

func checkForEmptyFBCred() error {
	if clientFacebookPageURL != "" && certificateDir != "" {
		return nil
	}
	return fmt.Errorf("Missing Items")
}

func checkForEmptyInstaCred() error {
	if appInstaUsername != "" && appInstaPassword != "" {
		return nil
	}
	return fmt.Errorf("Missing Instagram credentials")
}

func loadInstaProfile() {
	prof.Insta.Username = appInstaUsername
	prof.Insta.Password = appInstaPassword
}

func splitOnEqual(str string) []string {
	subStr := strings.Split(str, "=")
	return subStr
}

func getFileContent(filename string) ([]byte, error) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func scanCreds(str []byte) (id string, secret string, err error) {
	//var id, secret string
	scanner := bufio.NewScanner(bytes.NewReader(str))
	for scanner.Scan() {
		txt := strings.TrimSpace(scanner.Text())
		if txt == "" {
			continue
		}
		subStr := splitOnEqual(txt)
		if len(subStr) == 2 {
			x := strings.ToLower(strings.TrimSpace(subStr[0]))
			if x == "clientid" {
				id = strings.TrimSpace(subStr[1])
			}
			if x == "clientsecret" {
				secret = strings.TrimSpace(subStr[1])
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return "", "", scanner.Err()
	}
	return
}

// getCertFileName takes typ=0 or 1 => (0 = cert, 1 = key, 2 = credsFile)
func getCertFileName(typ int) string {
	switch typ {
	case 0:
		return filepath.Join(certificateDir, certFileName)
	case 1:
		return filepath.Join(certificateDir, keyFileName)
	case 2:
		return filepath.Join(certificateDir, fbIDSecretFileName)
	}
	return ""
}

func settingsTab() widgets.QWidget_ITF {
	facebookGroup := getFacebookGroupWidget()
	instaGroup := getInstaGroupWidget()
	importGroup := getImportGroupWidget()

	var layout = widgets.NewQGridLayout2()
	layout.AddWidget(instaGroup, 0, 0, 0)
	layout.AddWidget(facebookGroup, 1, 0, 0)
	layout.AddWidget(importGroup, 2, 0, 0)

	var secondTabWrapper = widgets.NewQWidget(tabWidget, 0)
	secondTabWrapper.SetLayout(layout)
	layout.SetAlign(core.Qt__AlignTop)
	return secondTabWrapper
}

var goToURLLineEdit *widgets.QLineEdit

func getFacebookGroupWidget() widgets.QWidget_ITF {
	var (
		facebookGroup     = widgets.NewQGroupBox2("Facebook Page Token Setup", nil)
		fbPageURLLabel    = widgets.NewQLabel2("FB Page link", nil, 0)
		fbPageURLLineEdit = getCustomLineEdit("https://www.facebook.com/<your-publicly-available-page>/")
		TokenURLLabel     = widgets.NewQLabel2("Token URL", nil, 0)

		certDirLabel    = widgets.NewQLabel2("Cert Directory", nil, 0)
		certDirLineEdit = getCustomLineEdit("Path to directory containing certificates")
		fbRegisterBtn   = getCustomPushButton("Register to get permanent token")
		credErrors      = widgets.NewQLabel2("", nil, 0)
	)
	goToURLLineEdit = getCustomLineEdit("Go to this URL when populated")
	certDirLineEdit.SetText(certificateDir)

	setLineEditColor(certDirLineEdit, func(txt string) {
		certificateDir = txt
	})
	credErrors.SetStyleSheet("color:red")
	credErrors.SetHidden(true)
	setLineEditColor(fbPageURLLineEdit, func(txt string) {
		clientFacebookPageURL = txt
	})

	setLineEditColor(goToURLLineEdit, func(txt string) {
		gotoURL = txt
	})

	fbRegisterBtn.ConnectClicked(func(c bool) {
		defer func() {
			//srvr.Shutdown(context.Background())
			fbRegisterBtn.SetDisabled(false)
			fbRegisterBtn.SetText("Register to get permanent token")
		}()
		err := checkForEmptyFBCred()
		if err != nil {
			setErroronLineEdit(fbRegisterBtn, credErrors, err)
			return
		}
		certFile := getCertFileName(0)
		keyFile := getCertFileName(1)
		credsFile := getCertFileName(2)
		if certFile == "" || keyFile == "" || credsFile == "" {
			setErroronLineEdit(fbRegisterBtn, credErrors, fmt.Errorf("Cert/fbcreds File(s) for local server not found in %v", certificateDir))
			return
		}
		contents, err := getFileContent(credsFile)
		if err != nil {
			setErroronLineEdit(fbRegisterBtn, credErrors, err)
			return
		}
		clientID, clientSecret, err = scanCreds(contents)
		if err != nil {
			setErroronLineEdit(fbRegisterBtn, credErrors, err)
			return
		}
		go func() {
			token, pageid, err := startServer(certFile, keyFile)
			if err != nil {
				setErroronLineEdit(fbRegisterBtn, credErrors, err)
				return
			}
			prof.FB.PageID = pageid
			prof.FB.PermToken = token
			// If we reach this far, its time to save prof into the file
			b, err := json.Marshal(prof)
			if err != nil {
				setErroronLineEdit(fbRegisterBtn, credErrors, err)
				return
			}
			//ioutil.WriteFile(configPath, b, 0755)
			encryptAndWriteToFile(configPath, b)
			fbRegisterBtn.SetDisabled(false)
			fbRegisterBtn.SetText("Register to get permanent token")
		}()

		fbRegisterBtn.SetDisabled(true)
		fbRegisterBtn.SetText("Please wait..")
		credErrors.SetHidden(true)
	})

	var facebookLayout = widgets.NewQGridLayout2()
	facebookLayout.AddWidget3(TokenURLLabel, 0, 0, 1, 1, 0)
	facebookLayout.AddWidget3(goToURLLineEdit, 0, 1, 1, 6, 0)
	facebookLayout.AddWidget3(fbPageURLLabel, 1, 0, 1, 1, 0)
	facebookLayout.AddWidget3(fbPageURLLineEdit, 1, 1, 1, 6, 0)
	facebookLayout.AddWidget3(certDirLabel, 2, 0, 1, 6, 0)
	facebookLayout.AddWidget3(certDirLineEdit, 2, 1, 1, 6, 0)
	facebookLayout.AddWidget3(fbRegisterBtn, 3, 0, 1, 7, 0)
	facebookLayout.AddWidget3(credErrors, 4, 0, 1, 7, 0)
	facebookGroup.SetLayout(facebookLayout)
	facebookGroup.SetMaximumHeight(350)
	return facebookGroup
}

func getInstaGroupWidget() widgets.QWidget_ITF {

	var (
		instaGroup       = widgets.NewQGroupBox2("Instagram Credentials", nil)
		userLabel        = widgets.NewQLabel2("Username:", nil, 0)
		passLabel        = widgets.NewQLabel2("Password:", nil, 0)
		usernameLineEdit = getCustomLineEdit("Username")
		passLineEdit     = getCustomLineEdit("Password")
		instaSubmitBtn   = getCustomPushButton("Submit")
		credErrors       = widgets.NewQLabel2("", nil, 0)
	)
	credErrors.SetStyleSheet("color:red")
	credErrors.SetHidden(true)

	instaSubmitBtn.ConnectClicked(func(c bool) {
		// save to a bpamedia file in home directory
		err := checkForEmptyInstaCred()
		if err != nil {
			setErroronLineEdit(instaSubmitBtn, credErrors, err)
			return
		}

		go func() {
			loadInstaProfile()
			b, err := json.Marshal(prof)
			if err != nil {
				setErroronLineEdit(instaSubmitBtn, credErrors, err)
				return
			}
			//ioutil.WriteFile(configPath, b, 0755)
			encryptAndWriteToFile(configPath, b)
			instaSubmitBtn.SetDisabled(false)
			instaSubmitBtn.SetText("Submit")
		}()

		instaSubmitBtn.SetDisabled(true)
		instaSubmitBtn.SetText("Please wait..")
		credErrors.SetHidden(true)
	})

	setLineEditColor(usernameLineEdit, func(txt string) {
		appInstaUsername = txt
	})

	setLineEditColor(passLineEdit, func(txt string) {
		appInstaPassword = txt
	})

	passLineEdit.SetEchoMode(widgets.QLineEdit__Password)

	var instaLayout = widgets.NewQGridLayout2()
	instaLayout.AddWidget3(userLabel, 0, 0, 1, 1, 0)
	instaLayout.AddWidget3(usernameLineEdit, 0, 1, 1, 6, 0)
	instaLayout.AddWidget3(passLabel, 1, 0, 1, 1, 0)
	instaLayout.AddWidget3(passLineEdit, 1, 1, 1, 6, 0)
	instaLayout.AddWidget3(instaSubmitBtn, 2, 0, 1, 7, 0)
	instaLayout.AddWidget3(credErrors, 3, 0, 1, 7, 0)
	instaGroup.SetLayout(instaLayout)
	instaGroup.SetMaximumHeight(200)
	return instaGroup
}

func getImportGroupWidget() widgets.QWidget_ITF {

	var (
		importGroup     = widgets.NewQGroupBox2("Import Credentials", nil)
		importSubmitBtn = getCustomPushButton("Import")
	)

	var importLayout = widgets.NewQGridLayout2()
	importLayout.AddWidget3(importSubmitBtn, 0, 0, 1, 1, 0)
	importGroup.SetLayout(importLayout)
	importGroup.SetMaximumHeight(150)

	importSubmitBtn.ConnectClicked(func(c bool) {
		var fileDialog = widgets.NewQFileDialog2(mainWindow, "Open crendentials JSON file", "", "")
		fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
		fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
		var mimeTypes = []string{"application/json"}
		fileDialog.SetMimeTypeFilters(mimeTypes)
		if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
			return
		}
		if len(fileDialog.SelectedFiles()) == 0 {
			return
		}

		jsonFile := fileDialog.SelectedFiles()[0]
		if jsonFile != "" {
			// copy file from one plae to another
			b, _ := ioutil.ReadFile(jsonFile)
			ioutil.WriteFile(configPath, b, 0755)
		}
	})
	return importGroup
}

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	centralLayout                                            *widgets.QGridLayout
	bpaTab1                                                  = "BPA Media (Upload)"
	bpaTab2                                                  = "BPA Media (Social Media Settings)"
	appInstaUsername, appInstaPassword, appFacebookPage      string
	appFacebookToken, appFacebookAppID, appFacebookAppSecret string
	currentFont, currentFontStyleSheet                       string
	uploadErrs                                               map[string]error // contains errors from instagram and/or facebook
)

type insta struct {
	Username, Password string
}

type fbS struct {
	PageID, PermToken string
}

type profile struct {
	Insta    insta  `json:"insta"`
	FB       fbS    `json:"fb"`
	ImageDir string `json:"imagedir"`
}

const cfgFileName = "bpamedia.json"

var (
	prof                = &profile{}
	mainWindow          *widgets.QMainWindow
	homeDir, configPath string
	errs                string
)

func setUpProfile() {
	// setup homeDir, prof variable
}

func main() {
	widgets.NewQApplication(len(os.Args), os.Args)
	var err error
	homeDir, err = homedir.Dir()
	if err == nil {
		configPath = filepath.Join(homeDir, cfgFileName)
	}
	certificateDir = homeDir // see settingsTab.go
	if configPath != "" {
		//credFile, err := ioutil.ReadFile(configPath)
		credFile, err := readFromFileAndDecrypt(configPath)
		if err != nil {
			errs = fmt.Sprintf("Decryption failed: %s", err.Error())
			//mainWindow.StatusBar().ShowMessage(core.QDir_ToNativeSeparators(fmt.Sprintf("Decryption failed: %s", err.Error())), 0)
			//return
		}
		err = json.Unmarshal(credFile, prof)
		if err != nil {
			errs = fmt.Sprintf("Could not read json file: %s", err.Error())
			//mainWindow.StatusBar().ShowMessage(core.QDir_ToNativeSeparators(fmt.Sprintf("Could not read json file: %s", err.Error())), 0)
			//return
		}
	}

	mainWindow = widgets.NewQMainWindow(nil, 0)
	mainWindow.SetWindowTitle(bpaTab1)

	centralWidget := widgets.NewQWidget(nil, 0)
	centralLayout = widgets.NewQGridLayout(centralWidget)

	addTabs()

	mainWindow.SetCentralWidget(centralWidget)
	mainWindow.ShowDefault()
	currentFont = mainWindow.FontInfo().Family()
	appFont := gui.NewQFont2(currentFont, 13, -1, false)
	mainWindow.SetFont(appFont)

	var menu = mainWindow.MenuBar().AddMenu2("&File")
	a := menu.AddAction("&Open")
	a.ConnectTriggered(func(checked bool) {
		if tabWidget.CurrentIndex() != 0 {
			return
		}
		addTabs()
		fileOpen(mainWindow)
	})
	a.SetShortcut(gui.NewQKeySequence2("Ctrl+D", gui.QKeySequence__NativeText))

	currentFontStyleSheet = mainWindow.StyleSheet()

	widgets.QApplication_Exec()
}

func addWidget(widget widgets.QWidget_ITF) {
	centralLayout.AddWidget(widget, 0, 0, core.Qt__AlignCenter)
}

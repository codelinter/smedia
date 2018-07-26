package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	fb "github.com/huandu/facebook"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/quick"
	"github.com/therecipe/qt/widgets"
	goinsta "gopkg.in/ahmdrz/goinsta.v2"
)

var (
	imageUploadPath string
	caption         string
)

func uploadToFacebook() error {
	sess := &fb.Session{}
	sess.SetAccessToken(prof.FB.PermToken)
	sess.Version = "v3.0"
	//pageid := ""
	photo, err := os.Open(imageUploadPath)
	if err != nil {
		return err
	}
	defer photo.Close()
	res, err := sess.Api("/"+prof.FB.PageID+"/photos", fb.POST, fb.Params{
		"message": caption,
		"source":  fb.Data(filepath.Base(imageUploadPath), photo),
	})
	if err != nil {
		return err
	}
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
func uploadToInstagram() error {
	insta := goinsta.New(prof.Insta.Username, prof.Insta.Password)

	if err := insta.Login(); err != nil {
		return err
	}
	var photo io.ReadCloser
	var err error
	photo, _ = getCroppedReader(imageUploadPath)
	if photo == nil {
		photo, err = os.Open(imageUploadPath)
		if err != nil {
			return err
		}

	}
	defer photo.Close()

	_, err = insta.UploadPhoto(photo, caption, 87, 0)
	if err != nil {
		return err
	}
	return nil
}

func firstTabWidgets(path string) widgets.QWidget_ITF {
	codd := `Image {
		fillMode: Image.PreserveAspectFit
		source: "file:///%s" 
		}`
	code := fmt.Sprintf(codd, path)
	imgWidget := createWidget("Image", code)
	var (
		imageUploadGroup = widgets.NewQGroupBox2("", nil)
		imageUploadLabel = widgets.NewQLabel2(path, nil, 0)

		uploadBtn       = getCustomPushButton("Upload")
		credErrors      = widgets.NewQLabel2("", nil, 0)
		captionLineEdit = getCustomLineEdit("Caption")
	)
	credErrors.SetStyleSheet("color:red")
	credErrors.SetHidden(true)

	setLineEditColor(captionLineEdit, func(txt string) {
		caption = txt
	})

	uploadBtn.ConnectClicked(func(c bool) {
		if errs != "" {
			uploadBtn.SetText("Error..")

			pushBtn.SetDisabled(false)
			mainWindow.StatusBar().ShowMessage(core.QDir_ToNativeSeparators(""), 1)
			return
		}

		go func() {
			var tasksDone string
			err := uploadToFacebook()
			if err != nil {
				setErroronLineEdit2(uploadBtn, credErrors, err)
				uploadBtn.SetText("Error..")

				pushBtn.SetDisabled(false)
				mainWindow.StatusBar().ShowMessage(core.QDir_ToNativeSeparators(""), 1)
				return

				tasksDone += fmt.Sprintf("(Facebook [ERROR] --> %s)", err.Error())
			} else {
				tasksDone += "(Facebook --> Success)"
			}
			err2 := uploadToInstagram()
			tasksDone += " / "
			if err2 != nil {
				tasksDone += fmt.Sprintf("(Instagram [ERROR] --> %s)", err2.Error())
			} else {
				tasksDone += "(Instagram --> Success)"
			}
			if err2 == nil && err == nil {
				uploadBtn.SetStyleSheet("color:green")
				uploadBtn.SetText(tasksDone)
			} else {
				uploadBtn.SetStyleSheet("color:orange")
				uploadBtn.SetText(tasksDone)
			}
			pushBtn.SetDisabled(false)
			mainWindow.StatusBar().ShowMessage(core.QDir_ToNativeSeparators(""), 1)
		}()
		uploadText := fmt.Sprintf("Uploading...%s , Caption Length: %d ", imageUploadPath, len(caption))
		mainWindow.StatusBar().ShowMessage(core.QDir_ToNativeSeparators(uploadText), 0)
		uploadBtn.SetText("Please Wait..")
		uploadBtn.SetDisabled(true)
		pushBtn.SetDisabled(true)
	})

	var imageUploadLayout = widgets.NewQGridLayout2()
	imageUploadLayout.AddWidget(imageUploadLabel, 0, 0, 0)
	imageUploadLayout.AddWidget3(captionLineEdit, 1, 0, 1, 2, 0)
	imageUploadLayout.AddWidget3(imgWidget, 2, 0, 1, 2, 0)
	imageUploadLayout.AddWidget3(uploadBtn, 3, 0, 1, 2, 0)
	imageUploadLayout.AddWidget3(credErrors, 4, 0, 1, 2, 0)
	imageUploadGroup.SetLayout(imageUploadLayout)

	return imageUploadGroup
	/* var layout = widgets.NewQGridLayout2()
	layout.AddWidget(settingsGroup, 1, 0, 0) */
}

func writeProfToFile() {
	b, _ := json.Marshal(prof)
	//ioutil.WriteFile(configPath, b, 0755)
	encryptAndWriteToFile(configPath, b)
}

func fileOpen(app widgets.QWidget_ITF) {
	var fileDialog = widgets.NewQFileDialog2(app, "Open JPEG Image...", prof.ImageDir, "")
	fileDialog.SetAcceptMode(widgets.QFileDialog__AcceptOpen)
	fileDialog.SetFileMode(widgets.QFileDialog__ExistingFile)
	var mimeTypes = []string{"image/jpeg"}
	fileDialog.SetMimeTypeFilters(mimeTypes)
	if fileDialog.Exec() != int(widgets.QDialog__Accepted) {
		return
	}
	if len(fileDialog.SelectedFiles()) == 0 {
		return
	}
	if prof.ImageDir == "" && fileDialog.SelectedFiles()[0] != "" {
		prof.ImageDir = filepath.Dir(fileDialog.SelectedFiles()[0])
		writeProfToFile()
	} else if prof.ImageDir != filepath.Dir(fileDialog.SelectedFiles()[0]) {
		prof.ImageDir = filepath.Dir(fileDialog.SelectedFiles()[0])
		writeProfToFile()
	}
	imageUploadPath = fileDialog.SelectedFiles()[0]
	loadImageWidget(imageUploadPath)
}

func reloadImageWidget(imgUploadPath string) {
	addTabs()
	loadImageWidget(imgUploadPath)
}

func loadImageWidget(imageUploadPath string) {
	groupedWidget := firstTabWidgets(imageUploadPath)
	firstTabLayout.AddWidget(groupedWidget, 0, 0)
}

var pushBtn *widgets.QPushButton

func buttonWidget() widgets.QWidget_ITF {
	pushBtn = getCustomPushButton("Add Image")
	pushBtn.ConnectClicked(func(checked bool) {
		addTabs()
		fileOpen(pushBtn)
	})
	return pushBtn
}

func createWidgetWithContext(name, code, ctxName string, ctx core.QObject_ITF) *quick.QQuickWidget {
	quickWidget := quick.NewQQuickWidget(nil)
	quickWidget.SetWindowTitle(name)

	quickWidget.RootContext().SetContextProperty(ctxName, ctx)

	quickWidget.SetResizeMode(quick.QQuickWidget__SizeRootObjectToView)

	path := filepath.Join(os.TempDir(), "tmp"+strings.Replace(name, " ", "", -1)+".qml")
	ioutil.WriteFile(path, []byte("import QtQuick 2.0\nimport QtQuick.Layouts 1.3\nimport QtQuick.Controls 1.4\n"+code), 0644)
	quickWidget.SetSource(core.QUrl_FromLocalFile(path))
	return quickWidget
}

func createWidget(name, code string) *quick.QQuickWidget {
	return createWidgetWithContext(name, code, "", nil)
}

var tabWidget *widgets.QTabWidget

func addWidgetToQStack(firstTabLayout *widgets.QVBoxLayout, wdg widgets.QWidget_ITF) {
	somewidget := widgets.NewQWidget(nil, 0)
	someWidgetLayout := widgets.NewQVBoxLayout2(somewidget)
	someWidgetLayout.AddWidget(wdg, 0, 0)
	firstTabLayout.AddWidget(somewidget, 0, 0)
}

func getFontStyle(typ int) string {
	if typ == 0 {
		return fmt.Sprintf("QLineEdit { color : red; font-family: %q; font: 17px }", currentFont)
	}
	return fmt.Sprintf("QLineEdit { color : blue; font-family: %q; font: 17px }", currentFont)
}

func getCustomLineEdit(ph string) *widgets.QLineEdit {
	le := widgets.NewQLineEdit(nil)
	le.SetPlaceholderText(ph)
	le.SetStyleSheet(getFontStyle(0))
	return le
}

func getCustomPushButton(name string) *widgets.QPushButton {
	btn := widgets.NewQPushButton2(name, nil)
	btn.SetMinimumHeight(35)
	return btn
}

func setLineEditColor(le *widgets.QLineEdit, f func(txt string)) *widgets.QLineEdit {
	le.ConnectTextChanged(func(txt string) {
		if len(txt) == 0 {
			le.SetStyleSheet(getFontStyle(0))
		} else {
			le.SetStyleSheet(getFontStyle(1))
		}
		f(txt)
	})
	return le
}

/*
// getPageID gets the page-id from the page-name
func getPageID() error {
	session := &fb.Session{}
	session.SetAccessToken(prof.FB.PermToken)
	session.Version = "v3.0"
	res, err := session.Get(appFacebookPage, fb.Params{
		"fields": "id",
	})
	if err != nil {
		return err
	}
	var id interface{}
	var ok bool
	if id, ok = res["id"]; !ok {
		return fmt.Errorf("Response doesn't PageID")
	}
	prof.FB.PageID = id.(string)
	return nil

}
*/
// getPermFBToken is called only from the settings tab
// its only supposed to be called when resetting the permanent token
func getPermFBToken() error {
	session := &fb.Session{}
	session.SetAccessToken(appFacebookToken)
	session.Version = "v3.0"
	if appFacebookAppID != "" && appFacebookAppSecret != "" && appFacebookToken != "" && appFacebookPage != "" {
		res, err := session.Get("/oauth/access_token", fb.Params{
			"grant_type":        "fb_exchange_token",
			"client_id":         appFacebookAppID,
			"client_secret":     appFacebookAppSecret,
			"fb_exchange_token": appFacebookToken,
		})
		if err != nil {
			return err
		}
		var toke interface{}
		var ok bool
		if toke, ok = res["access_token"]; !ok {
			return fmt.Errorf("Response doesn't have access token")
		}

		prof.FB.PermToken = toke.(string)
		return nil
	}
	return fmt.Errorf("Missing Facebook credentials")
}

func setErroronLineEdit(credentialsSubmitBtn *widgets.QPushButton, credErrors *widgets.QLabel, err error) {
	credentialsSubmitBtn.SetDisabled(false)
	credentialsSubmitBtn.SetText("Submit")
	credErrors.SetHidden(false)
	credErrors.SetText(err.Error())
}

func setErroronLineEdit2(credentialsSubmitBtn *widgets.QPushButton, credErrors *widgets.QLabel, err error) {
	credentialsSubmitBtn.SetDisabled(false)
	credentialsSubmitBtn.SetText("Upload")
	credErrors.SetHidden(false)
	credErrors.SetText(err.Error())
}

var firstTabLayout *widgets.QVBoxLayout

func uploadTab() widgets.QWidget_ITF {
	firstTabWidget := widgets.NewQWidget(nil, 0)
	firstTabWidget.SetWindowTitle("Widget")
	firstTabLayout = widgets.NewQVBoxLayout2(firstTabWidget)

	somewidget := widgets.NewQWidget(nil, 0)
	someWidgetLayout := widgets.NewQVBoxLayout2(somewidget)
	someWidgetLayout.AddWidget(buttonWidget(), 0, 0)
	firstTabLayout.AddWidget(somewidget, 0, 0)
	firstTabLayout.SetAlign(core.Qt__AlignTop)

	return firstTabWidget
}

func addTabs() {
	allTabs := make([]widgets.QWidget_ITF, 2)

	tabWidget = widgets.NewQTabWidget(nil)
	tabWidget.SetWindowTitle("Final Tabs")
	allTabs[0] = uploadTab()
	allTabs[1] = settingsTab()
	var label string
	for i, v := range allTabs {
		if i == 0 {
			label = "Upload"
		} else {
			label = "Settings"
		}
		tabWidget.AddTab(v, label)
	}
	tabWidget.SetFixedSize2(1024, 720)
	tabWidget.ConnectTabBarClicked(func(idx int) {
		switch idx {
		case 0:
			mainWindow.SetWindowTitle(bpaTab1)
		case 1:
			mainWindow.SetWindowTitle(bpaTab2)

		}
	})
	addWidget(tabWidget)
}

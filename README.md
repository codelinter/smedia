# What is smedia?
Upload images to instagram and facebook pages in one go. `smedia` is a desktop app hence it ***does not*** ship with FB app credentials. This app is created to help [bpaindia](http://www.bpaindia.org) to be more productive in the social media platform as they tend to post(upload) images quite often and doing so seperately slows them down.

# Licensing
I personally dont have any restrictive clause so you are free to do whatever you want. But since the underlying technology uses the `qt framework`, their license still holds true. Please visit `https://www.qt.io` for more information on this.

# Lets start

# Step 1. -> Get the binary

## Option 1. Build from source

0.0 Install golang
 
1.1 `go get -u github.com/therecipe/qt/cmd/...`

1.2 Install appropriate docker image 

    `docker pull therecipe/qt:windows_32_static` 
    
    (use linux tag for linux enviornment)

1.3 `go get github.com/mintyowl/smedia`

1.4 Change `appSymPassword` in creds.go file before compiling. Make sure to add creds.go in .gitignore file if you version controlling.

1.5 `qtdeploy -docker build linux`


# Option 2. Use compiled version instead
### If building by source is not an option, find appropriate binary from the `releases` section. Currently windows and linux binaries are supported.

# Step 2 -> Prepare server certificates and FB credentials
(This step is only required to be done once for the first time set up only. This one is a little bit cumbersome as well. Bare with me.)

2.1 Software(smedia) runs a localhost server over TLS on port `12345`, which requires server certificate and keys for `localhost` domain. Easiest way to get these files is by running this simple command

`go run $GOROOT/src/crypto/tls/generate_cert.go -host localhost`
This will create two files named cert.pem and key.pem. Rename both to `example.com.pem` and `example.com-key.pem` respectively

2.2 Get FB credentials (App ID and App Secret). 

2.2.1 Go to `https://developers.facebook.com`

2.2.2 Create an app. You will have a link like this 

    `https://developers.facebook.com/apps/<your-app-id>/dashboard/`

2.2.3 It may ask for type of app you are creating, choose www (or web)

2.2.4 Go to /fb-login/settings and provide `https://localhost:12345` as the Redirect URI option(this step is required)

2.2.5 Go to settings/basic and grab the `App ID` and `App Secret`

2.2.6 Provide your privacy declaration URI. 

# Step 3 -> Pack your credentials
(This step is continuation of Step 2 and is only required once)

3.1 Create a file name `fbcreds` and populate it with content as shown below

    clientid = <your-app-id-from-step-2.2.5>
    clientsecret = <your-app-secret-from-step-2.2.5>

3.2 Put `fbcreds`, `example.com.pem` and `example.com-key.pem` (all three of them) in the same folder (anywhere you link)

# Step 4 -> Start smedia (if not running already)
4.1 Go the settings tab and fill in instagram credentials

4.2 In `Cert Directory` provide the path to the folder that contains files mentioned in step 3

4.3 FB Page link is the link to your publicly available page
    `https://www.facebook.com/<your-public-page>/`

4.4 Hit `Register to get permanent token`. This creates and a link on `Token URL`. Copy that link and hop into your browser to visit that link, then login. Accept the terms and come back to the software(smedia). If everything went well you will not see any errors in red color in smedia.

# THATS IT, YOU'RE DONE ! ! !

Now go to the upload tab and start uploading images to both instagram and FB simultaneously.



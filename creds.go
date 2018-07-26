package main

const (
	// appSymPassword is a symmetric password for NaCl, used to encrypt and decrypt above creds among others
	// resultant encrypted data will then gets stored in the correct json file
	// appSymPassword is not used by the client. Its only meant for the software
	appSymPassword string = "myPassword"
)

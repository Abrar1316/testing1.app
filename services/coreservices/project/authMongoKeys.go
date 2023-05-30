package project

import (
	"fmt"
	"log"
	"net/http"

	auth "github.com/Snawoot/go-http-digest-auth-client"
)

func AuthMongokeys(orgID, publicKey, privateKey string) (bool, error) {

	// Send the authentication request to the MongoDB Atlas API
	apiURL := fmt.Sprintf("https://cloud.mongodb.com/api/atlas/v1.0/orgs/%s/", orgID)

	client := &http.Client{
		Transport: auth.NewDigestTransport(publicKey, privateKey, http.DefaultTransport),
	}

	response, err := client.Get(apiURL)
	if err != nil {
		log.Fatalln(err)
		return false, fmt.Errorf("Internal server error ")
	}

	defer response.Body.Close()

	// Check the response status code
	if response.StatusCode == http.StatusOK {
		log.Println("Authentication successful!")
		return true,nil

	} else {
		log.Printf("Authentication failed")
		return false,fmt.Errorf("authentication failed")
	}
}

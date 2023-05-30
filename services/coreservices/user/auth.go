package user

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var (
	secretKey = []byte("billingApp_SecretKey")
)

func GenerateJwt(userid int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = userid

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", errors.New("unable to create JWT access token")
	}
	return tokenString, nil
}

func VerifyAccess(endpoint http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		user := params["uid"]

		userId, _ := strconv.Atoi(user)

		cookie, err := r.Cookie("session_id")
		if err != nil {
			fmt.Printf("Error in VerifyAccess: unable to extract User Access : %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Generate token based on incoming request UserId
		generateToken, err := GenerateJwt(userId)
		if err != nil {
			fmt.Printf("Error in VerifyAccess: unable to verify user : %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Extracting the token value
		incomingToken := cookie.Value

		// Checks to verify that no other user have access to resources
		if generateToken != incomingToken {
			fmt.Println("Error in VerifyAccess: User is not Authorized")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		token, err := jwt.Parse(incomingToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("server error in parsing the token")
			}
			return secretKey, nil
		})

		if err != nil {
			fmt.Printf("Error in VerifyAccess: Error in verifying token : %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if token.Valid {
			endpoint(w, r)
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}

	})
}

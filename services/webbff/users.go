package webbff

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/coreservices/user"
)

func loginPageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")
		// TODO: check form inputs
		fmt.Println(email, password)

		user := &user.UserInfo{
			Email:    email,
			Password: password,
		}

		userData, err := user.Login()
		if err != nil {
			t := template.Must(template.ParseFiles(loginTamplates...))
			t.Execute(w, struct{ Error string }{err.Error()})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   userData.Token,
			Path:    "/",
			Expires: time.Now().Add(30 * time.Minute),
		})

		http.Redirect(w, r, fmt.Sprintf("/%d/days/30", userData.Key), http.StatusFound)
		return

	} else {
		t := template.Must(template.ParseFiles(loginTamplates...))
		t.Execute(w, struct{}{})

	}
}

func signupPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		user := &user.UserInfo{
			Name:     name,
			Email:    email,
			Password: password,
		}

		err := user.SignupService()
		if err != nil {

			t := template.Must(template.ParseFiles(signinTamplates...))
			t.Execute(w, struct{ Error string }{err.Error()})

			return
		}

		userData, err := user.Login()
		if err != nil {
			t := template.Must(template.ParseFiles(loginTamplates...))
			t.Execute(w, struct{ Error string }{err.Error()})
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "session_id",
			Value:   userData.Token,
			Path:    "/",
			Expires: time.Now().Add(30 * time.Minute),
		})

		// Note: Use a 3xx status code to redirect the client
		// https://go.dev/src/net/http/status.go
		http.Redirect(w, r, fmt.Sprintf("/%d/days/30", userData.Key), http.StatusFound)

	} else {
		t := template.Must(template.ParseFiles(signinTamplates...))
		t.Execute(w, struct{}{})
	}
}

func logoutPageHandler(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // Set the MaxAge to a negative value to expire the cookie immediately
	}
	http.SetCookie(w, &cookie)

	// redirect to home page
	http.Redirect(w, r, "/", http.StatusFound)
}

// This is a default handler, if path is not defined, then it will redirect to login page
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

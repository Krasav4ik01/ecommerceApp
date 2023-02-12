package controllers

import (
	"errors"
	"github.com/Krasav4ik01/ecommerceApp/config"
	"github.com/Krasav4ik01/ecommerceApp/entities"
	"github.com/Krasav4ik01/ecommerceApp/libraries"
	"github.com/Krasav4ik01/ecommerceApp/models"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type UserInput struct {
	Username string
	Pass     string
}

var userModel = models.NewUserModel()
var validation = libraries.NewValidation()

// главная страница.
// Если мы не залогинились, то откроется стр логина

func Index(w http.ResponseWriter, r *http.Request) {

	session, _ := config.Store.Get(r, config.SESSION_ID)

	if len(session.Values) == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {

		if session.Values["loggedIn"] != true {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		} else {

			data := map[string]interface{}{
				"name": session.Values["name"],
			}

			temp, _ := template.ParseFiles("templates/index.html")
			temp.Execute(w, data)
		}
	}
}

//страница логина

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		temp, err := template.ParseFiles("templates/login.html")
		if err != nil {
			panic(err)
		}
		temp.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		r.ParseForm()
		UserInput := &UserInput{
			Username: r.Form.Get("username"),
			Pass:     r.Form.Get("pass"),
		}
		var user entities.User
		userModel.Where(&user, "username", UserInput.Username)
		//editing
		var message error
		if user.Username == "" {
			message = errors.New("Invalid Username or Password!")
		} else {
			// password
			errPassword := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(UserInput.Pass))
			if errPassword != nil {
				message = errors.New("Invalid Username or Password!")
			}
		}

		if message != nil {

			data := map[string]interface{}{
				"error": message,
			}

			temp, _ := template.ParseFiles("templates/login.html")
			temp.Execute(w, data)
		} else {
			// set session
			session, _ := config.Store.Get(r, config.SESSION_ID)

			session.Values["loggedIn"] = true
			session.Values["username"] = user.Username
			session.Values["pass"] = user.Pass
			session.Values["name"] = user.Name

			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}

}

// функция logout, чтобы выйти из акк

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := config.Store.Get(r, config.SESSION_ID)
	// delete session
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		temp, _ := template.ParseFiles("templates/register.html")
		temp.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		// процесс регистрации

		r.ParseForm()

		user := entities.User{
			Name:     r.Form.Get("name"),
			Email:    r.Form.Get("email"),
			Username: r.Form.Get("username"),
			Pass:     r.Form.Get("pass"),
			Confpass: r.Form.Get("confpass"),
		}

		errorMessages := validation.Struct(user)

		if errorMessages != nil {

			data := map[string]interface{}{
				"validation": errorMessages,
				"user":       user,
			}

			temp, _ := template.ParseFiles("templates/register.html")
			temp.Execute(w, data)
		} else {

			// защифровка пароля
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Pass), bcrypt.DefaultCost)
			user.Pass = string(hashPassword)

			// инсерт в БД
			userModel.Create(user)
			//Успешная регистрация
			data := map[string]interface{}{
				"event": "Success register",
			}
			temp, _ := template.ParseFiles("templates/register.html")
			temp.Execute(w, data)
		}
	}

}
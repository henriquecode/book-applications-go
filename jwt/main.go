package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"

	"github.com/cristalhq/jwt"
)

type user struct {
	id       int
	Name     string
	Age      int
	Email    string
	password string
	Book     book
}

type book struct {
	id    int
	Title string
}

type userClaims struct {
	IsAdmin bool   `json:"is_admin"`
	Email   string `json:"email"`
	ID      int    `json:"id"`
}

var (
	users []user
	books []book
)

const secretKey string = "key-secret-123"

func init() {
	makeDataBooks()
	makeDataUsers()
	setRoutes()
}

func main() {
	http.ListenAndServe(":9000", nil)
}

func makeDataBooks() {

	books = []book{
		book{1, "Livro do Leandro"},
		book{2, "Livro da Andressa"},
	}
}

func makeDataUsers() {

	users = []user{
		user{1, "Leandro", 28, "leandro@gmail.com", "123456", books[0]},
		user{2, "Andressa", 30, "andressa@gmail.com", "123", books[1]},
	}
}

func setRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		view(w, "home", "")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			view(w, "login", "")
		} else {
			r.ParseForm()
			email := r.Form.Get("user")
			password := r.Form.Get("password")

			fmt.Println(email, password)
			user, ok := getUserByEmailPassword(email, password)

			if ok == false {
				view(w, "login", struct {
					Sucess  bool
					Message string
				}{
					false,
					"Usuário não encontrado",
				})
			} else {
				token := generateToken(user)

				view(w, "token", struct{ Token string }{
					token,
				})
			}
		}
	})

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		tokenReceived := params.Get("token")

		verifier, errVerifier := jwt.NewVerifierHS(jwt.HS256, []byte(secretKey))

		if errVerifier != nil {
			panic(errVerifier)
		}

		newToken, err := jwt.ParseString(tokenReceived, verifier)

		if err != nil {
			panic(err)
		}

		var formatUser userClaims
		json.Unmarshal(newToken.RawClaims(), &formatUser)

		dataUser, ok := getUserByID(formatUser.ID)

		if ok == false {
			view(w, "info", struct {
				Sucess  bool
				Message string
			}{
				false,
				"Usuário não encontrado",
			})
		} else {
			view(w, "info", struct {
				Success bool
				Data    user
			}{
				false,
				dataUser,
			})
		}
	})
}

func generateToken(u user) string {
	key := []byte(secretKey)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)

	if err != nil {
		panic(err)
	}

	builder := jwt.NewBuilder(signer)

	claims := &userClaims{
		IsAdmin: true,
		Email:   u.Email,
		ID:      u.id,
	}

	token, err := builder.Build(claims)

	if err != nil {
		panic(err)
	}

	return token.String()
}

func getUserByID(id int) (user, bool) {
	for _, u := range users {
		if u.id == id {
			return u, true
		}
	}

	return user{}, false
}

func getUserByEmailPassword(email, password string) (user, bool) {
	for _, u := range users {
		if u.Email == email && u.password == password {
			return u, true
		}
	}

	return user{}, false
}

func view(w http.ResponseWriter, view string, data interface{}) {
	directory, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	tmpl := template.Must(template.ParseFiles(directory + "/views/" + view + ".tpl"))

	tmpl.Execute(w, data)
}

package web

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/tonyalaribe/converze-api/config"
	"github.com/tonyalaribe/entrenet/messages"

	"github.com/tonyalaribe/converze-api/models"
)

//LoginResponse sent to the cllient, carrying the token, upon login
type LoginResponse struct {
	User    models.User
	Message string
	Token   string
}

//GetMeDetails returns details about the current logged user, based on the content of his JWT token
func GetMeDetails(w http.ResponseWriter, r *http.Request) {
	user, err := Userget(r)
	if err != nil {
		log.Println(err)
	}

	userDetails, err := user.Get(config.Get(), user.Email)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userDetails)
}

//NewUserAccount is a handler for new user account creation
func NewUserAccount(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}
	//log.Println(user)

	data := struct {
		Message string
	}{
		Message: "Account Created Successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)

	err = user.Add(config.Get()) //placed at this point so that image handling does not delay returning response
	if err != nil {
		log.Println(err)
	}
}

//UserLogin exchanges an email and password for an authentication token, which is then attached to every request to the database
func UserLogin(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	log.Println(user)
	log.Printf("%+v", user)

	user2, err := user.Get(config.Get(), user.Email)
	if err != nil {
		log.Println(err)
	}

	err = bcrypt.CompareHashAndPassword(user2.Password, []byte(user.P))
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(messages.ErrWrongPassword)
	}
	user.P = ""

	response, err := GenerateJWT(user2)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(messages.ErrInternalServer)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

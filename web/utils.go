package web

import (
	"net/http"
	"time"

	"github.com/tonyalaribe/converze-api/config"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
	"github.com/tonyalaribe/converze-api/messages"
	"github.com/tonyalaribe/converze-api/models"
)

//Userget reads the json web token(JWT) content from context and marshals it ito a user struct,
func Userget(r *http.Request) (models.User, error) {
	//id := context.Get(r, "UserID")
	//u := context.Get(r, "User")

	u := r.Context().Value("User")

	var user models.User
	err := mapstructure.Decode(u, &user)

	if err != nil {
		return user, err
	}
	return user, nil

}

//Setup could be run at a predefined route, and be used to make important initializations, like default interests, admin users, etc.
func Setup(w http.ResponseWriter, r *http.Request) {
	messages.WriteError(w, messages.Success)
}

//GenerateJWT urn user details into a hasked token that can be used to recognize the user in the future.
func GenerateJWT(user models.User) (resp LoginResponse, err error) {
	claims := jwt.MapClaims{}

	// set our claims
	claims["User"] = user
	claims["Name"] = user.Name

	// set the expire time

	claims["exp"] = time.Now().Add(time.Hour * 24 * 30 * 12).Unix() //24 hours inn a day, in 30 days * 12 months = 1 year in milliseconds

	// create a signer for rsa 256
	t := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	pub, err := jwt.ParseRSAPrivateKeyFromPEM(config.Get().Encryption.Private)
	if err != nil {
		return
	}
	tokenString, err := t.SignedString(pub)

	if err != nil {
		return
	}

	resp = LoginResponse{
		User:    user,
		Message: "Token succesfully generated",
		Token:   tokenString,
	}

	return

}

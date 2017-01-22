package models

import (
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"golang.org/x/crypto/bcrypt"

	uuid "github.com/satori/go.uuid"
	"github.com/tonyalaribe/converze-api/config"
)

//Listing struct holds all submitted form data for listings
type User struct {
	Name            string
	Email           string
	Image           string
	Password        []byte `json:"-"`
	P               string `json:"Password" bson:"-"`
	DateCreated     time.Time
	Coordinates     string
	CompletedSignup bool
	Location        struct {
		Type        string
		Coordinates []float64
	}
	BioDetails struct {
		SecondaryInstitution string
		TertiaryInstitution  string
		CourseOfStudy        string
		Profession           string
		Workplace            string
		StateOfOrigin        string
		CountryOfOrigin      string
		StateOfResidence     string
		CountryOfResidence   string
	}
	Interests []string
}

//Return a single users details from d, based on his email
func (r User) Get(conf *config.Conf, email string) (User, error) {
	user := User{}
	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	collection := conf.Database.C(config.USERSCOLLECTION).With(mgoSession)

	err := collection.Find(bson.M{
		"email": email,
	}).One(&user)

	if err != nil {
		log.Println(err)
		return user, err
	}
	return user, nil
}

//Add a user to the database
func (r User) Add(conf *config.Conf) error {
	phash, err := bcrypt.GenerateFromPassword([]byte(r.P), conf.PasswordEncryptionCost)
	if err != nil {
		log.Println(err)
		return err
	}
	r.Password = phash //Store the hash of the password into the password field
	r.P = ""           //Delete the text based password

	if r.Image != "" {
		var imageURL string
		imagepath := "profileimage/" + uuid.NewV1().String()
		imageURL, err = UploadBase64Image(conf.S3Bucket, r.Image, imagepath, 250)
		if err != nil {
			log.Println(err)
			return err
		}
		r.Image = imageURL
	}

	r.DateCreated = time.Now()

	mgoSession := conf.Database.Session.Copy()
	defer mgoSession.Close()

	collection := conf.Database.C(config.USERSCOLLECTION).With(mgoSession)

	// Index
	index := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = collection.EnsureIndex(index)
	if err != nil {
		log.Println(err)
		return err
	}

	err = collection.Insert(r)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

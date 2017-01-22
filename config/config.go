package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	mgo "gopkg.in/mgo.v2"
)

//Conf nbfmjh
type Conf struct {
	MongoDB                string
	MongoServer            string
	Database               *mgo.Database
	S3Bucket               *s3.Bucket
	PasswordEncryptionCost int
	Encryption             struct {
		Private []byte
		Public  []byte
	}
}

var (
	config     Conf
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

const (
	ADVERT  = "advert"
	LISTING = "listing"

	FACEBOOK = "facebook"

	USERSCOLLECTION = "Users"

	TAGSCOLLECTION = "Tags"
)

func Init() {
	MONGOSERVER := os.Getenv("MONGO_URL")
	MONGODB := os.Getenv("MONGODB")
	if MONGOSERVER == "" {
		log.Println("No mongo server address set, resulting to default address")
		MONGOSERVER = "127.0.0.1:27017"
		MONGODB = "converze"
		//MONGODB = "yellowListings"
		//MONGODB = "y"
		//mongodb://localhost
	}

	session, err := mgo.Dial(MONGOSERVER)
	if err != nil {
		log.Println(err)
	}
	log.Println(session)

	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err)
	}
	client := s3.New(auth, aws.USWest2)
	bucketname := os.Getenv("bucket_name")
	if bucketname == "" {
		bucketname = "test-past3"
	}
	bucket := client.Bucket(bucketname)

	config = Conf{
		MongoDB:                MONGODB,
		MongoServer:            MONGOSERVER,
		Database:               session.DB(MONGODB),
		S3Bucket:               bucket,
		PasswordEncryptionCost: 10,
	}
	log.Println(basepath)
	config.Encryption.Public, err = ioutil.ReadFile("./config/encryption_keys/public.pem")
	if err != nil {
		log.Println("Error reading public key")
		log.Println(err)
		return
	}

	config.Encryption.Private, err = ioutil.ReadFile("./config/encryption_keys/private.pem")
	if err != nil {
		log.Println("Error reading private key")
		log.Println(err)
		return
	}

	log.Printf("mongoserver %s", MONGOSERVER)
}

func Get() *Conf {
	return &config
}

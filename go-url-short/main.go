package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

const mongo_db = "urlshortdb"
const mongo_collection = "urlshort"

// if running in docker, use the service name
// const mongo_url = "mongodb://mongo:27017"
const mongo_default_url = "mongodb://localhost:27017/urlshortdb"
const mongo_default_username = "admin"
const mongo_default_password = "password"

var client_url = "http://localhost:5173"

var client *mongo.Client
var collection *mongo.Collection

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getClient() (*mongo.Client, error) {
	client_url = getEnv("CLIENT_URL", client_url)
	mongo_url := getEnv("MONGO_URL", mongo_default_url)
	mongo_username := getEnv("MONGO_USERNAME", mongo_default_username)
	mongo_password := getEnv("MONGO_PASSWORD", mongo_default_password)

	fmt.Println("Connecting to MongoDB at", mongo_url)
	fmt.Println("Username:", mongo_username)

	clientOptions := options.Client().ApplyURI(mongo_url)

	// check if username and password are provided
	if clientOptions.Auth != nil {
		clientOptions.Auth.Username = mongo_username
		clientOptions.Auth.Password = mongo_password
	}

	options.Client().SetMaxConnIdleTime(10 * time.Second)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type URLDoc struct {
	ShortURL  string    `bson:"shorturl"`
	LongURL   string    `bson:"longurl"`
	CreatedAt time.Time `bson:"createdAt"`
}

func shorten(w http.ResponseWriter, r *http.Request) {
	// get the URL from the request
	longurl := r.FormValue("url")
	if longurl == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	// check if the URL is valid
	if _, err := url.Parse(longurl); err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	ctx, _ := context.WithTimeout(r.Context(), 10*time.Second)

	// generate a new short URL
	var shorturl string
	md5hash := md5.Sum([]byte(longurl))
	start := 0
	for {
		var hashcode []byte
		hashcode = md5hash[start : start+2]
		// convert the bytes to decimal
		padding := make([]byte, 2-len(hashcode))
		hashcode = append(hashcode, padding...)
		key := binary.BigEndian.Uint16(hashcode)
		// add a random integer to the key
		key += uint16(time.Now().UnixNano())
		// use base64 encoding to generate a short URL
		shorturl = base64.URLEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(key), 10)))
		// remove padding characters '='
		shorturl = strings.TrimRight(shorturl, "=")
		// check collision
		count, err := collection.CountDocuments(ctx, bson.M{"short": shorturl})
		if err != nil {
			http.Error(w, "Error checking for collision", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			break
		}
		start += 2
	}

	// insert the new short URL into the database
	newDoc := URLDoc{
		ShortURL:  shorturl,
		LongURL:   longurl,
		CreatedAt: time.Now(),
	}
	_, err := collection.InsertOne(ctx, newDoc)
	if err != nil {
		http.Error(w, "Error inserting URL", http.StatusInternalServerError)
		return
	}

	// return the document as JSON
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, newDoc)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		http.Error(w, "Short URL is required", http.StatusBadRequest)
		return
	}

	ctx, _ := context.WithTimeout(r.Context(), 10*time.Second)
	var result URLDoc
	err := collection.FindOne(ctx, bson.M{"shorturl": code}).Decode(&result)
	if err != nil {
		notfoundURL := client_url + "/notfound"
		http.Redirect(w, r, notfoundURL, http.StatusFound)
		return
	}

	http.Redirect(w, r, result.LongURL, http.StatusFound)
}

func init() {
	var err error
	client, err = getClient()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		os.Exit(1)
	}

	// Check the connection
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB:", err)
		os.Exit(1)
	}
	fmt.Println("Connected to MongoDB")

	// get the database
	db := client.Database(mongo_db)

	// get the collection
	collection = db.Collection(mongo_collection)

	fmt.Println("Connected to MongoDB and collection", mongo_collection)
}

func main() {
	defer client.Disconnect(context.Background())

	// Add endpoint
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
	}))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("URL Shortener"))
	})
	r.Post("/shorten", shorten)
	r.Get("/{code}", redirect)

	http.ListenAndServe(":3000", r)
}

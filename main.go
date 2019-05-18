package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	//"log"
	"net/http"
	//"log"
	//"log"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2/bson"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

//mongo DB
const (
	hosts      = "188.166.239.20:27017"
	database   = "regblogs_mongo"
	username   = "pomtcom"
	password   = "P@ssw0rd"
	collection = "regblogs"
)

// The person Type (more like an object)
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type Blogs struct {
	Publish_no 			int 	`json:"publish_no"`
	Topic_name 			string 	`json:"topic_name"`
	Thumbnail			string	`json:"thumbnail"`
	Short_description	string	`json:"short_description"`
	Html_code			string	`json:"html_code"`
}

type HealthCheck struct {
	Response  string `json:"response"`
}

var people []Person


// Application HealthCheck
func GetAppHealthCheck(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	if (*r).Method == "GET" {

	}else if (*r).Method == "POST" {
		fmt.Println("NOT ALLOW for this method")
	}

	fmt.Println(("GET method for healthcheck endpoint is working"))
	var healthCheck HealthCheck
	healthCheck.Response = "OK"
	json.NewEncoder(w).Encode(healthCheck)

}


// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	fmt.Println("/people is executing")

	if (*r).Method == "GET" {

	}else if (*r).Method == "POST" {
		fmt.Println("NOT ALLOW for this method")
	}

	fmt.Println(("THIS IS GET METHOD"))
	fmt.Println("GetPeople web-service")
	json.NewEncoder(w).Encode(people)

}

// create a new item
func GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if (*r).Method == "POST" {
		fmt.Println("NOT ALLOW for this method")
	}else if (*r).Method == "GET" {

	}

	//blogs := connectAndQueryBlog()
	blogs := FirestoreConnect()
	json.NewEncoder(w).Encode(blogs)

}



// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

//func enableCors(w *http.ResponseWriter) {
//	(*w).Header().Set("Access-Control-Allow-Origin", "*")
//}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func connectAndQueryBlog() []Blogs {
	info := &mgo.DialInfo{
		Addrs:    []string{hosts},
		Timeout:  60 * time.Second,
		Database: database,
		Username: username,
		Password: password,
	}

	session, err1 := mgo.DialWithInfo(info)
	if err1 != nil {
		panic(err1)
	}

	//get collection (regblogs)
	col := session.DB(database).C(collection)

	count, err2 := col.Count()
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(fmt.Sprintf("Messages count: %d", count))

	var blogs []Blogs
	//blogs := []Blogs
	//blogs := Blogs{}
	//err := col.Find(bson.M{"publish_no": 1}).One(&blogs)
	err := col.Find(bson.M{}).All(&blogs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("blogs:", blogs)

	return blogs
}

// Application HealthCheck
func Firestoretest(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	if (*r).Method == "GET" {

	}else if (*r).Method == "POST" {
		fmt.Println("NOT ALLOW for this method")
	}

	fmt.Println(("GET method for firestore test is working"))

	blogs := FirestoreConnect()
	//fmt.Println(blogs)
	json.NewEncoder(w).Encode(blogs)

}

func FirestoreConnect() []Blogs {
	fmt.Println("Firestore connect is executing")
	sa := option.WithCredentialsFile("firestore-pomt.json")
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	q := client.Collection("regblogs")
	docs, err := q.Documents(context.Background()).GetAll()
	if err != nil {
		// TODO: Handle error.
	}

	var blogs []Blogs
	for _, doc := range docs {
		//fmt.Println(doc)

		var blog Blogs
		if err := doc.DataTo(&blog); err != nil {
			// TODO: Handle error.
		}
		//fmt.Println(blog)
		blogs = append(blogs, blog)

	}
	fmt.Println("Blogs query is completed")
	return blogs
}

// main function to boot up everything
func main() {
	fmt.Println("Prepare web-services latest")
	//log.Fatal("TEST12345")
	router := mux.NewRouter()
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/api/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	router.HandleFunc("/api/getallblogs", GetAllBlogs).Methods("GET")
	router.HandleFunc("/App-HealthCheck", GetAppHealthCheck).Methods("GET")
	router.HandleFunc("/api/firestoretest", Firestoretest).Methods("GET")

	fmt.Println("Web-services are starting")
	log.Fatal(http.ListenAndServe(":8087", router))


	//connectMongoDB()


}
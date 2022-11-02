package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
	Artists   string `json:"artists"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
	Relation  string `json:"relation"`
}
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relations    string   `json:"relations"`
	Relationinfo Relationss
}

type Relationss struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

var artists []Artist
var rel []Relationss
var response_Object *Response

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/groups/", handleArtists)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil { //to initialise our server on :8080
		log.Fatal("HTTP status 500 - Internal server error: %s", err)
	}
}

func retData(url string) ([]byte, error) {

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return body, nil
}

func handleArtists(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[8:]
	p, _ := strconv.Atoi(id)
	t, err := template.ParseFiles("artists.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, &artists[p-1])
	if err != nil {
		panic(err)
	}

}

func Home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "404 Not Found", 404)
	} else {
		url := "https://groupietrackers.herokuapp.com/api"
		data, err := retData(url)
		if err != nil {
			fmt.Println("Error retrieving api")
			log.Println(err)
			return
		}
		json.Unmarshal(data, &response_Object)
		artData, err := retData(response_Object.Artists)
		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(artData, &artists)

		relationData, err := retData(response_Object.Relation)

		if err != nil {
			log.Println(err)
		} else {
			relationData = relationData[9 : len(relationData)-2]
			json.Unmarshal(relationData, &rel)

			for i := range artists {
				artists[i].Relationinfo = rel[i]
			}

			t, errParse := template.ParseFiles("index.html")
			if errParse != nil {
				log.Println(err)
			} else {
				t.Execute(w, &artists)
			}
		}
	}
}
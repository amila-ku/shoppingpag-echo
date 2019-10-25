package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	// blank import is for adding swagger docs
	_ "github.com/amila-ku/shoppingpal/api/docs"
	"github.com/amila-ku/shoppingpal/pkg/entity"
	store "github.com/amila-ku/shoppingpal/pkg/store"
	"github.com/gorilla/mux"
	swagger "github.com/swaggo/http-swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ItemList hods the list of items
var ItemList = entity.NewItems()

// Home page handler
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

// health endpoint handler
func healthEndpoint(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}


// Add Item godoc
// @Summary Add an Item
// @Description add an item
// @ID add-item
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} entity.Item
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} entity.APIError "We need ID!!"
// @Failure 404 {object} entity.APIError "Can not find ID"
// @Failure 500 {object} entity.APIError "We had a problem"
// @Router /items/ [post]
func createNewItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: CreateNewItem")

	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	//fmt.Fprintf(w, "%+v", string(reqBody))
	var itm entity.Item
	json.Unmarshal(reqBody, &itm)

	// update our global item array to include our new item
	ItemList = append(ItemList, itm)

	// save to db
	db, err := store.NewTable("itemtable")

	if err != nil {
		log.Fatal("Failed to create table", err)
	}
	err = db.CreateItem(itm)
	if err != nil {
		log.Fatal("Unable to insert item", err)
	}

	fmt.Println(ItemList)

	prettyJSON(w, itm)

}
// List Single Item godoc
// @Summary List Single Item
// @Description get Item
// @Accept  json
// @Produce  json
// @Param id query string false "item search by id"
// @Success 200 {array} entity.Item
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} entity.APIError "We need ID!!"
// @Failure 404 {object} entity.APIError "Can not find ID"
// @Router /items [get]
func returnSingleItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	//Check items slice for matching item
	for _, item := range ItemList {

		if item.ID == key {
			prettyJSON(w, item)
		}
	}
}
// ListALLItems godoc
// @Summary List Items
// @Description get Items
// @Accept  json
// @Produce  json
// @Success 200 {array} entity.Item
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} entity.APIError "We need ID!!"
// @Failure 404 {object} entity.APIError "Can not find ID"
// @Router /items [get]
func returnAllItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllItems")

	//json.NewEncoder(w).Encode(ItemList)

	// Print Json with indents, the pretty way:
	prettyJSON(w, ItemList)

}
// DeleteItems godoc
// @Summary Delete Items
// @Description Delete Items
// @Accept  json
// @Produce  json
// @Param id query string false "item delete by id"
// @Success 200 {array} entity.Item
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} entity.APIError "We need ID!!"
// @Failure 404 {object} entity.APIError "Can not find ID"
// @Router /items/id [del]
func deleteItem(w http.ResponseWriter, r *http.Request) {
	// parse the path parameters
	vars := mux.Vars(r)
	// extract the `id` of the item
	id := vars["id"]

	//loop through all our items
	for index, item := range ItemList {
		// delete if item id matches
		if item.ID == id {
			ItemList = append(ItemList[:index], ItemList[index+1:]...)
		}
	}

}

func prettyJSON(w http.ResponseWriter, list interface{}) {
	pretty, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(pretty)
}

//HandleRequests defines all the route mappings
func HandleRequests() {
	// Remnants of default http handler for comparison
	// http.HandleFunc("/", returnAllItems)
	// log.Fatal(http.ListenAndServe(":10000", nil))

	// Creates a new router
	myRouter := mux.NewRouter().StrictSlash(true)

	// Application Operations related mappings
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/health", healthEndpoint)
	myRouter.PathPrefix("/metrics").Handler(promhttp.Handler())

	// App functionality mappings
	myRouter.PathPrefix("/swagger/").Handler(swagger.Handler(swagger.URL("http://localhost:10000/swagger/doc.json")))
	myRouter.HandleFunc("/items", returnAllItems).Methods("GET")
	myRouter.HandleFunc("/item/{id}", returnSingleItem).Methods("GET")
	myRouter.HandleFunc("/item/{id}", deleteItem).Methods("DELETE")
	myRouter.HandleFunc("/item", createNewItem).Methods("POST")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
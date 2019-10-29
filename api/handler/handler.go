package handler

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"

	//echo framework
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// blank import is for adding swagger docs
	_ "github.com/amila-ku/shoppingpal-echo/api/docs"
	"github.com/amila-ku/shoppingpal-echo/pkg/entity"
	store "github.com/amila-ku/shoppingpal-echo/pkg/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// ItemList hods the list of items
var ItemList = entity.NewItems()

// Home page handler
func homePage(c echo.Context) error {
	//c.Logger.Info("Endpoint Hit: homePage")
	return c.String(http.StatusOK, "Welcome to the HomePage!")
}

// health endpoint handler
func healthEndpoint(c echo.Context) error {
	//c.Logger.Info("Endpoint Hit: health")
	return c.String(http.StatusOK, "Up and Running")
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
func createNewItem(c echo.Context) error {
	fmt.Println("Endpoint Hit: CreateNewItem")

	// get the body of our POST request
	// return the string response containing the request body
	defer c.Request().Body.Close()

	var itm entity.Item
	
	err := json.NewDecoder(c.Request().Body).Decode(&itm)

	if err != nil {
		log.Printf("Failed processing addDog request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	//json.Unmarshal(reqBody, &itm)

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
	return c.JSON(http.StatusOK, itm)

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
func returnSingleItem(c echo.Context) error {
	//vars := mux.Vars(r)
	key := c.Param("id")

	//Check items slice for matching item
	for _, item := range ItemList {

		if item.ID == key {
			//prettyJSON(w, item)
			return c.JSON(http.StatusOK, item)
		}
	}
	return c.JSON(http.StatusOK, key)
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
func deleteItem(c echo.Context) error {
	// parse the path parameters
	// vars := mux.Vars(r)
	// extract the `id` of the item
	// id := vars["id"]

	id := c.Param("id")

	//loop through all our items
	for index, item := range ItemList {
		// delete if item id matches
		if item.ID == id {
			ItemList = append(ItemList[:index], ItemList[index+1:]...)
		}
	}

	return c.NoContent(http.StatusNoContent)

}

//HandleRequests defines all the route mappings
func HandleRequests() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Application Operations related mappings
	e.GET("/", homePage)
	e.GET("/health", healthEndpoint)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// myRouter.PathPrefix("/metrics").Handler()

	// OpenAPI3 docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// // App functionality mappings
	e.GET("/item/:id", returnSingleItem)
	e.DELETE("/item/:id", deleteItem)
	e.POST("/item", createNewItem)

	// myRouter.HandleFunc("/items", returnAllItems).Methods("GET")
	// myRouter.HandleFunc("/item/{id}", returnSingleItem).Methods("GET")
	// myRouter.HandleFunc("/item/{id}", deleteItem).Methods("DELETE")
	// myRouter.HandleFunc("/item", createNewItem).Methods("POST")

	e.Logger.Fatal(e.Start(":10000"))
}

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Service struct {
	UID                    string `json:"uid"`
	RunDate                string `json:"runDate"`
	Origin                 string `json:"origin"`
	Destination            string `json:"destination"`
	OriginWorkingTime      string `json:"originWorkingTime"`
	DestinationWorkingTime string `json:"destinationWorkingTime"`
	OriginPlatform         string `json:"originPlatform"`
	DestinationPlatform    string `json:"destinationPlatform"`
}

type ServiceByStation struct {
	Location struct {
		Name    string `json:"name"`
		Crs     string `json:"crs"`
		Tiploc  string `json:"tiploc"`
		Country string `json:"country"`
		System  string `json:"system"`
	} `json:"location"`
	Services []struct {
		LocationDetail struct {
			RealtimeActivated   bool   `json:"realtimeActivated"`
			Tiploc              string `json:"tiploc"`
			Crs                 string `json:"crs"`
			Description         string `json:"description"`
			GbttBookedArrival   string `json:"gbttBookedArrival"`
			GbttBookedDeparture string `json:"gbttBookedDeparture"`
			Origin              []struct {
				Tiploc      string `json:"tiploc"`
				Description string `json:"description"`
				WorkingTime string `json:"workingTime"`
				PublicTime  string `json:"publicTime"`
			} `json:"origin"`
			Destination []struct {
				Tiploc      string `json:"tiploc"`
				Description string `json:"description"`
				WorkingTime string `json:"workingTime"`
				PublicTime  string `json:"publicTime"`
			} `json:"destination"`
			IsCall                  bool   `json:"isCall"`
			IsPublicCall            bool   `json:"isPublicCall"`
			RealtimeArrival         string `json:"realtimeArrival"`
			RealtimeArrivalActual   bool   `json:"realtimeArrivalActual"`
			RealtimeDeparture       string `json:"realtimeDeparture"`
			RealtimeDepartureActual bool   `json:"realtimeDepartureActual"`
			Platform                string `json:"platform"`
			PlatformConfirmed       bool   `json:"platformConfirmed"`
			PlatformChanged         bool   `json:"platformChanged"`
			DisplayAs               string `json:"displayAs"`
		} `json:"locationDetail"`
		ServiceUID      string `json:"serviceUid"`
		RunDate         string `json:"runDate"`
		TrainIdentity   string `json:"trainIdentity"`
		RunningIdentity string `json:"runningIdentity"`
		AtocCode        string `json:"atocCode"`
		AtocName        string `json:"atocName"`
		ServiceType     string `json:"serviceType"`
		IsPassenger     bool   `json:"isPassenger"`
	} `json:"services"`
}

type ServiceByRoute struct {
	Location struct {
		Name    string `json:"name"`
		Crs     string `json:"crs"`
		Tiploc  string `json:"tiploc"`
		Country string `json:"country"`
		System  string `json:"system"`
	} `json:"location"`
	Services []struct {
		LocationDetail struct {
			RealtimeActivated   bool   `json:"realtimeActivated"`
			Tiploc              string `json:"tiploc"`
			Crs                 string `json:"crs"`
			Description         string `json:"description"`
			GbttBookedArrival   string `json:"gbttBookedArrival"`
			GbttBookedDeparture string `json:"gbttBookedDeparture"`
			Origin              []struct {
				Tiploc      string `json:"tiploc"`
				Description string `json:"description"`
				WorkingTime string `json:"workingTime"`
				PublicTime  string `json:"publicTime"`
			} `json:"origin"`
			Destination []struct {
				Tiploc      string `json:"tiploc"`
				Description string `json:"description"`
				WorkingTime string `json:"workingTime"`
				PublicTime  string `json:"publicTime"`
			} `json:"destination"`
			IsCall                  bool   `json:"isCall"`
			IsPublicCall            bool   `json:"isPublicCall"`
			RealtimeArrival         string `json:"realtimeArrival"`
			RealtimeArrivalActual   bool   `json:"realtimeArrivalActual"`
			RealtimeDeparture       string `json:"realtimeDeparture"`
			RealtimeDepartureActual bool   `json:"realtimeDepartureActual"`
			Platform                string `json:"platform"`
			PlatformConfirmed       bool   `json:"platformConfirmed"`
			PlatformChanged         bool   `json:"platformChanged"`
			DisplayAs               string `json:"displayAs"`
		} `json:"locationDetail"`
		ServiceUID      string `json:"serviceUid"`
		RunDate         string `json:"runDate"`
		TrainIdentity   string `json:"trainIdentity"`
		RunningIdentity string `json:"runningIdentity"`
		AtocCode        string `json:"atocCode"`
		AtocName        string `json:"atocName"`
		ServiceType     string `json:"serviceType"`
		IsPassenger     bool   `json:"isPassenger"`
	} `json:"services"`
}

func main() {
	getURI()
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/services/:station", getServiceByDestination)
	router.GET("/services/:station/to/:toStation", getServiceByRoute)
	router.Run("localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getServiceByDestination(c *gin.Context) {
	baseURL := os.Getenv("REALTIME_TRAINS_URI")
	username := os.Getenv("TRAINS_USERNAME")
	password := os.Getenv("PASSWORD")
	station := c.Param("station")

	client := &http.Client{}

	req, err1 := http.NewRequest("GET", baseURL+"json/search/"+station, nil)
	if err1 != nil {
		fmt.Println(err1)

		panic(err1)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	fmt.Println(req.Header.Get("Authorization"))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(resp.StatusCode, err.Error())
	}

	var currServices ServiceByStation
	error1 := json.Unmarshal([]byte(body), &currServices)
	if error1 != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	c.IndentedJSON(http.StatusOK, currServices)

}

func getServiceByRoute(c *gin.Context) {
	baseURL := os.Getenv("REALTIME_TRAINS_URI")
	username := os.Getenv("TRAINS_USERNAME")
	password := os.Getenv("PASSWORD")
	station := c.Param("station")
	toStation := c.Param("toStation")

	client := &http.Client{}

	req, err1 := http.NewRequest("GET", baseURL+"json/search/"+station+"/to/"+toStation, nil)
	if err1 != nil {
		fmt.Println(err1)

		panic(err1)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	fmt.Println(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.IndentedJSON(resp.StatusCode, err.Error())
	}

	var currServices ServiceByRoute
	error1 := json.Unmarshal([]byte(body), &currServices)
	if error1 != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	c.IndentedJSON(http.StatusOK, currServices)
}

func getURI() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

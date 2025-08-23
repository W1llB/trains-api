package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ncruces/go-strftime"
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

type ServiceDetail struct {
	ServiceUID           string `json:"serviceUid"`
	RunDate              string `json:"runDate"`
	ServiceType          string `json:"serviceType"`
	IsPassenger          bool   `json:"isPassenger"`
	TrainIdentity        string `json:"trainIdentity"`
	PowerType            string `json:"powerType"`
	TrainClass           string `json:"trainClass"`
	AtocCode             string `json:"atocCode"`
	AtocName             string `json:"atocName"`
	PerformanceMonitored bool   `json:"performanceMonitored"`
	Origin               []struct {
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
	Locations []struct {
		RealtimeActivated   bool   `json:"realtimeActivated"`
		Tiploc              string `json:"tiploc"`
		Crs                 string `json:"crs"`
		Description         string `json:"description"`
		GbttBookedDeparture string `json:"gbttBookedDeparture,omitempty"`
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
		RealtimeDeparture       string `json:"realtimeDeparture,omitempty"`
		RealtimeDepartureActual bool   `json:"realtimeDepartureActual,omitempty"`
		Platform                string `json:"platform"`
		PlatformConfirmed       bool   `json:"platformConfirmed"`
		PlatformChanged         bool   `json:"platformChanged"`
		Line                    string `json:"line,omitempty"`
		LineConfirmed           bool   `json:"lineConfirmed,omitempty"`
		DisplayAs               string `json:"displayAs"`
		GbttBookedArrival       string `json:"gbttBookedArrival,omitempty"`
		RealtimeArrival         string `json:"realtimeArrival,omitempty"`
		RealtimeArrivalActual   bool   `json:"realtimeArrivalActual,omitempty"`
		ServiceLocation         string `json:"serviceLocation,omitempty"`
		Path                    string `json:"path,omitempty"`
		PathConfirmed           bool   `json:"pathConfirmed,omitempty"`
	} `json:"locations"`
	RealtimeActivated bool   `json:"realtimeActivated"`
	RunningIdentity   string `json:"runningIdentity"`
}

func main() {
	getURI()
	router := gin.Default()
	router.Use(ErrorHandler())
	router.Use(CORSMiddleware())
	router.GET("/services/:station", getServicesByDestination)
	router.GET("/services/:station/to/:toStation", getServicesByRoute)
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

// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next() // Step1: Process the request first.

        // Step2: Check if any errors were added to the context
        if len(c.Errors) > 0 {
            // Step3: Use the last error
            err := c.Errors.Last().Err

            // Step4: Respond with a generic error message
            c.JSON(http.StatusInternalServerError, map[string]any{
                "success": false,
                "message": err.Error(),
            })
        }

        // Any other steps if no errors are found
    }
}

func getURI() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getServicesByDestination(c *gin.Context) {
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

func getServicesByRoute(c *gin.Context) {
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
	}
	servicesList := currServices.Services
	detailedServicesList := []ServiceDetail{}
	for i := 0; i < len(servicesList); i++ {
		detailedServicesList = append(detailedServicesList, getServiceDetail(servicesList[i].ServiceUID))
	}
	fmt.Println(detailedServicesList)
	c.IndentedJSON(http.StatusOK, detailedServicesList)
}

func getServiceDetail(serviceUid string) ServiceDetail {

	baseURL := os.Getenv("REALTIME_TRAINS_URI")
	username := os.Getenv("TRAINS_USERNAME")
	password := os.Getenv("PASSWORD")
	// now, timeError := time.Parse("01/02/2006", time.Now().String())

	client := &http.Client{}
	req, err1 := http.NewRequest("GET", baseURL+"json/service/"+serviceUid+"/"+strftime.Format("%Y/%m/%d", time.Now()), nil)
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
		fmt.Println("Error reading body:", err)
		// c.IndentedJSON(resp.StatusCode, err.Error())
	}

	var detailedService ServiceDetail
	error1 := json.Unmarshal([]byte(body), &detailedService)
	if error1 != nil {
		fmt.Println("Error decoding JSON:", error1)
	}
	return detailedService

}

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type InefficientServersResponse struct {
	Status  string    `json:"status"`
	Servers []Servers `json:"servers"`
}

type Servers struct {
	Ip       string `json:"ip"`
	HostName string `json:"host_name"`
	Active   bool   `json:"active"`
}

func getAllInefficiantServers(threshold int) ([]string, error) {
	//return []string{"mta-prod-5", "mta-prod-3"}, nil
	client := &http.Client{Timeout: 10 * time.Second}
	endpoint := "https://optimizing-servers/get"
	values := map[string]int{"threshold": threshold}
	jsonReq, err := json.Marshal(values)

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatalf("Error Occurred. %+v", err)
		return nil, err
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
		return nil, err
	}

	defer response.Body.Close()

	// read json http response
	jsonDataFromHttp, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Error in reading the response body. %+v", err)
		return nil, err
	}

	var inefficientServersResponse InefficientServersResponse

	err = json.Unmarshal([]byte(jsonDataFromHttp), &inefficientServersResponse) // here!

	if err != nil {
		log.Fatalf("Error in unmarshalling the response body. %+v", err)
		return nil, err
	}

	if inefficientServersResponse.Status == "failed" {
		log.Fatal("Status of the response is failed.")
		return nil, errors.New("failed status")
	}

	HostNameCounttMap := make(map[string]int)

	for _, emp := range inefficientServersResponse.Servers {
		if emp.Active == false {
			HostNameCounttMap[emp.HostName]++
		}
	}

	finalHostName := []string{}

	for hostName, count := range HostNameCounttMap {
		if count >= threshold {
			finalHostName = append(finalHostName, hostName)
		}
	}

	return finalHostName, nil
}

func main() {

	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	value, ok := os.LookupEnv("THRESHOLD")

	if !ok {
		log.Fatal("No environment variable found")
		return
	}

	threshold, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal("Error in converting string to int")
	}

	fmt.Println("The threshold is", threshold)
	servers, err := getAllInefficiantServers(threshold)
	if err != nil {
		log.Fatalf("Error occuring feting from the API %v", err)
		return
	}

	fmt.Println("The inefficient servers of threshols", threshold, " are: ")
	for _, val := range servers {
		fmt.Println(val)
	}
}

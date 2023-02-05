package main

import (
	"log"
	"testing"
)

func TestAllEfficientServers(t *testing.T) {

	expectedTestOutput := []string{"mta-prod-1", "mta-prod-2"}
	actualTestOutput, err := getAllInefficiantServers(1)
	if err != nil {
		log.Fatalf("Error occuring feting from the API with threshold %d", 1)
	} else {
		for k, v := range actualTestOutput {
			if v != expectedTestOutput[k] {
				t.Errorf("Output %s not equal to expected %s", v, expectedTestOutput[k])
			}
		}
	}

}

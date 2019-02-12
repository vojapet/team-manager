package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"os"
	"flag"
	"log"
)

type Team struct {
	Name        string
	Description string
	MemberCount int
}

func readTeamApi(url string, login string, passwd string) []Team {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	req.SetBasicAuth(login, passwd)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	switch res.StatusCode {
	case http.StatusOK:
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		teams := make([]Team, 0)
		json.Unmarshal(body, &teams)
		log.Print("Request processed successfully.")
		return teams

	case http.StatusUnauthorized:
		log.Printf("Request processed unsuccessfully [%d] - unauthorized", res.StatusCode)
		os.Exit(1)
	default:
		log.Printf("Request processed unsuccessfully [%d]", res.StatusCode)
		os.Exit(1)
	}
	return nil
}

func printTeams(theTeams []Team, aTeam string) {
	for _, team := range(theTeams) {
		if aTeam != "" && aTeam != team.Name {
			continue
		}
		fmt.Printf("%s,%d\n", team.Name, team.MemberCount)
	}
}

func main() {
	urlPtr := flag.String("url", "http://localhost:8000/api/team", "team api url")
	loginPtr := flag.String("login", "", "restapi login name")
	passwdPtr := flag.String("password", "", "restapi password")
	teamPtr := flag.String("team", "", "specify team - no team means all the teams")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Client expect the JSON retrieved from API is array of 'Team' object. Team is struct with 'name', 'description' and ''membercount' fields.\n\n")
	}
	
	flag.Parse()
	printTeams(readTeamApi(*urlPtr, *loginPtr, *passwdPtr), *teamPtr)
}

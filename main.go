package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ratings []struct {
	Id   int `json:"Id"`
	Star int `json:"Star"`
}

type reviewes struct {
	Id       int    `json:"Id"`
	Star     int    `json:"Star"`
	Reviewer string `json:"Reviewer"`
	Review   string `json:"Review"`
	Color    string `json:"color"`
}

func main() {
	reviewer1 := reviewes{
		Id:       1,
		Star:     0,
		Reviewer: "Reviewer1",
		Review:   "An extremely entertaining play by Shakespeare. The slapstick humour is refreshing!",
		Color:    "",
	}
	reviewer2 := reviewes{
		Id:       2,
		Star:     0,
		Reviewer: "Reviewer2",
		Review:   "Absolutely fun and entertaining. The play lacks thematic depth when compared to other plays by Shakespeare.",
		Color:    "",
	}
	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get("http://localhost:8000/rating")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(data))

		var rat ratings
		json.Unmarshal(data, &rat)

		for _, v := range rat {
			if v.Id == 1 {
				reviewer1.Star = v.Star
			} else if v.Id == 2 {
				reviewer2.Star = v.Star
			}
		}

		reviewer := []reviewes{reviewer1, reviewer2}
		bs, err := json.Marshal(reviewer)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(bs)
	})
	http.ListenAndServe(":80", nil)

}

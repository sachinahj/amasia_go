package yelp

import "fmt"

type CategoryConfig struct {
	Alias            string   `json:"alias"`
	Title            string   `json:"title"`
	CountryWhitelist []string `json:"country_whitelist"`
	CountryBlacklist []string `json:"country_blacklist"`
	Parents          []string `json:"parents"`
}

type Category struct {
	Alias      string `json:"alias"`
	Title      string `json:"title"`
	CreatedAt  string
	ModifiedAt string
}

func (c *Category) Update() {

	fmt.Println("update this shit", c)

}

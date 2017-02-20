package yelp

type CategoryConfig struct {
	Alias            string   `json:"alias"`
	Title            string   `json:"title"`
	CountryWhitelist []string `json:"country_whitelist"`
	CountryBlacklist []string `json:"country_blacklist"`
	Parents          []string `json:"parents"`
}

type Category struct {
	Id         string
	Alias      string
	Title      string
	CreatedAt  string
	ModifiedAt string
}

type CategoriesConfig []CategoryConfig
type Categories []Category

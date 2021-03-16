package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/oauth2"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type Repo struct {
	Language  string
	URL       string
	Name      string
	StarCount int
}

type Org struct {
	URL      string
	ReposURL string
	Name     string
	Repos    []Repo
}

var personalAccessToken string
var q string
var lang string
var loc string
var page int = 0
var perPage int = 10
var count int = 0

func init() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("specify language and location")
	}

	lang = args[0]
	loc = args[1]

	if len(args) == 3 {
		perPage, _ = strconv.Atoi(args[2])
	}

	q = fmt.Sprintf("language:%s+type:org+location:%s", lang, loc)

	f, err := os.Open(".env")
	if err != nil {
		fmt.Println("no .env file found")
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := strings.Split(scanner.Text(), "=")
		if len(l) == 2 {
			fmt.Println(l[0])
			fmt.Println(l[1])
			personalAccessToken = l[1]
			return
		}
		fmt.Println("check your .env")
		return
	}
}

func main() {
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)

	display(client)

}

func fetchMore(client *github.Client) (*github.UsersSearchResult, error) {
	fmt.Printf("Fetching page %d\n", page+1)
	opts := &github.SearchOptions{ListOptions: github.ListOptions{Page: page + 1, PerPage: perPage}}
	sr, _, err := client.Search.Users(context.Background(), q, opts)
	return sr, err
}

func display(client *github.Client) {
	for {
		sr, err := fetchMore(client)
		if err != nil {
			log.Fatalf("error while searching: %v", err)
		}

		orgs := []Org{}
		for _, org := range sr.Users {
			count++

			o := Org{URL: org.GetHTMLURL(), ReposURL: org.GetReposURL(), Name: org.GetLogin()}

			rq := fmt.Sprintf("org:%s", org.GetLogin())
			ropts := &github.SearchOptions{Sort: "stars", ListOptions: github.ListOptions{PerPage: perPage}}
			projects, _, err := client.Search.Repositories(context.Background(), rq, ropts)
			if err != nil {
				fmt.Println("error: ", err.Error())
				continue
			}

			rs := []Repo{}
			for _, p := range projects.Repositories {
				l := strings.TrimSpace(p.GetLanguage())
				if l != "" && strings.ToLower(l) == strings.ToLower(strings.TrimSpace(lang)) {
					r := Repo{Language: l, URL: p.GetHTMLURL(), StarCount: p.GetStargazersCount(), Name: p.GetName()}
					rs = append(rs, r)
				}
			}

			if len(rs) > 0 {
				o.Repos = rs
				orgs = append(orgs, o)
			}
		}

		data := [][]string{}

		for _, or := range orgs {
			for _, r := range or.Repos {
				data = append(data, []string{or.Name, r.Name, fmt.Sprint(r.StarCount), r.Language, r.URL})
			}
		}

		sort.Slice(data, func(i, j int) bool {
			inum, err := strconv.Atoi(data[i][2])
			if err != nil {
				return false
			}
			jnum, err := strconv.Atoi(data[j][2])
			if err != nil {
				return false
			}

			return jnum < inum
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Organization", "Repo-Name", "Star", "Language", "URL"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output

		fmt.Printf("count is %d and total is %d\n", count, sr.GetTotal())
		if count == sr.GetTotal() {
			fmt.Println("scanned all repos")
			return
		}

		var input string
		for {
			fmt.Print("Do you want to fetch more companies? (y/n) ")
			fmt.Scanln(&input)
			input = strings.TrimSpace(input)
			if input == "y" || input == "n" {
				break
			}
		}
		if input == "y" {
			page++
		}
		if input == "n" {
			fmt.Println("Exit")
			break
		}
	}

}

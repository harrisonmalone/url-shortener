package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type App struct {
	db *[]Link
}

type Link struct {
	id        int
	url       string
	createdAt time.Time
}

func links() *[]Link {
	return &[]Link{
		{
			id:        100000,
			url:       "https://google.com",
			createdAt: time.Now(),
		},
		{
			id:        100001,
			url:       "https://apple.com",
			createdAt: time.Now().Local().Add(100),
		},
	}
}

func (a App) findByHash(hash string) (Link, error) {
	for _, link := range *a.db {
		idToHash := strconv.FormatInt(int64(link.id), 36)
		if idToHash == hash {
			return link, nil
		}
	}
	return Link{}, fmt.Errorf("No link found")
}

func (a *App) create() {
	rows := *a.db
	length := len(rows)
	lastLink := rows[length-1]
	link := Link{
		id:        lastLink.id + 1,
		url:       "https://example.com",
		createdAt: time.Now(),
	}
	*a.db = append(*a.db, link)
}

func (a *App) show() {
	for _, link := range *a.db {
		hash := strconv.FormatInt(int64(link.id), 36)
		row := fmt.Sprintf("id: %d, original url: %s, short url: https://short.com/%s", link.id, link.url, hash)
		fmt.Println(row)
	}
}

func (a App) menu() {
	fmt.Println("1. Create shortened link")
	fmt.Println("2. Show all links")
	fmt.Println("3. Redirect to website with link")
	fmt.Println("4. Exit")
}

func getHash() string {
	fmt.Println("Enter short url")
	var userInput string
	fmt.Scanln(&userInput)
	stringSlice := strings.Split(userInput, "/")
	return stringSlice[len(stringSlice)-1]
}

func (a App) selection() int {
	var userInput int
	fmt.Scanln(&userInput)
	return userInput
}

func main() {
	app := App{
		db: links(),
	}
	for true {
		app.menu()
		userInput := app.selection()
		if userInput == 1 {
			app.create()
			fmt.Println("Url created!")
		} else if userInput == 2 {
			app.show()
		} else if userInput == 3 {
			link, err := app.findByHash(getHash())
			if err != nil {
				fmt.Println("No url found!")
			} else {
				exec.Command("open", link.url).Start()
			}
		} else {
			os.Exit(0)
		}
	}
}

// 1. user enters into textbox link
// 2. goes to POST backend endpoint
// 3. backend endpoint creates new link made up of primary key, url and createdAt, would be SQL insert
// 4. primary key is auto imcremented by 1
// 5. id returned from DB is converted to base 36, append this base 36 string to your domain
// => example.com/255s
// 6. user able to copy this url
// 7. when user actually enters this url into browser request sent to GET link backend endpoint
// 8. extract base 36 from url
// 9. lookup link in DB with base 36, SELECT * FROM LINKS WHERE hash = param
// 10. one link should be returned from DB
// 11. return redirect from endpoint that takes user to url of link returned from DB
// 12. all urls are uniquely identified by id and base_36_hash

// notes
// reason why we use base 36 instead of just id is to ensure url is as short as possible
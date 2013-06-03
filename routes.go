package main

import (
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

/*
	/{team}/{board}
*/

var urlValidator = regexp.MustCompile("^/([a-zA-Z0-9\\-]+)/([a-zA-Z0-9\\-]+)$")

var views *template.Template

func formatTime(t *time.Time, layout string) string {
	return t.Format(layout)
}

func takeTop(s Records, top int) Records {
	if len(s) > top {
		return s[:top]
	} else {
		return s
	}
}

func skipTop(s Records, top int) Records {
	if len(s) > top {
		return s[top:]
	} else {
		return nil
	}
}

func makeBoardHandler(fn func(http.ResponseWriter, *http.Request, *Board)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//println("URL: " + r.URL.Path)
		m := urlValidator.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		board, err := LoadBoard(m[1], m[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fn(w, r, board)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, board *Board) {
	if r.Method == "POST" {
		who := strings.TrimSpace(r.FormValue("who"))
		if who == "" {
			who = "Unknown"
		}

		email := strings.TrimSpace(r.FormValue("email"))

		when := time.Now()

		score, err := strconv.ParseInt(r.FormValue("score"), 10, 0)
		if err != nil {
			score = 0
		}

		// Store the last entered information in a cookie so we can populate
		// the form with it in the future.
		cookie := &http.Cookie{Name: "memory", Value: "who"}
		http.SetCookie(w, cookie)

		recordUpdated := false
		if email != "" {
			// If we have an email, update any existing record for this person.
			for i := range board.Records {
				if strings.EqualFold(board.Records[i].Email, email) {
					board.Records[i].Email = email
					board.Records[i].When = when
					board.Records[i].Score = score
					recordUpdated = true
					break
				}
			}
		} else {
			// If we have a nickname only, update any existing record based on the nick.		
			for i := range board.Records {
				if strings.EqualFold(board.Records[i].Who, who) {
					board.Records[i].Who = who
					board.Records[i].When = when
					board.Records[i].Score = score
					recordUpdated = true
					break
				}
			}
		}

		if !recordUpdated {
			board.Records = append(board.Records, &Record{Who: who, When: when, Score: score})
		}

		// Sort the board's scores.
		sort.Sort(board.Records)

		// Only keep the top 10 scores.
		max := 10
		if len(board.Records) > max {
			board.Records = board.Records[:max]
		}

		if err = board.Save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/"+board.Team+"/"+board.Name, http.StatusFound)
		return
	} else {
		if err := views.ExecuteTemplate(w, "board.gohtml", board); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func init() {
	funcMap := template.FuncMap{
		"formatTime": formatTime,
		"takeTop":    takeTop,
		"skipTop":    skipTop,
	}

	views = template.Must(template.New("board.gohtml").Funcs(funcMap).ParseFiles("./views/board.gohtml"))

	http.HandleFunc("/", makeBoardHandler(viewHandler))

	http.Handle("/static/", http.FileServer(http.Dir("./")))
}

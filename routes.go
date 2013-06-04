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

// Use a regex to parse the team and board name from the URL.
var urlValidator = regexp.MustCompile("^/([a-zA-Z0-9\\-]+)/([a-zA-Z0-9\\-]+)(/([a-zA-Z0-9\\-]+))?/?$")

var views *template.Template

// Helper function for formatting timestamps in our templates.
func formatTime(t *time.Time, layout string) string {
	return t.Format(layout)
}

// Helper function for grabbing the first X records in our template.
func takeTop(s Records, top int) Records {
	if len(s) > top {
		return s[:top]
	} else {
		return s
	}
}

// Helper function for grabbing all /except/ the first X records in our template.
func skipTop(s Records, top int) Records {
	if len(s) > top {
		return s[top:]
	} else {
		return nil
	}
}

// Wrap a handler in logic that parses the team and board tokens from the URL
// and loads the Board instance for our handler to work with easily.
func makeBoardHandler(fn func(http.ResponseWriter, *http.Request, *Board, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//println("URL: " + r.URL.Path)
		m := urlValidator.FindStringSubmatch(r.URL.Path)
		if m == nil {
			//http.NotFound(w, r)
			homepageHandler(w, r)
			return
		}

		board, err := LoadBoard(m[1], m[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fn(w, r, board, m[4])
	}
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if err := views.ExecuteTemplate(w, "home.gohtml", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func descHandler(w http.ResponseWriter, r *http.Request, board *Board) {
	if r.Method == "POST" {
		newDesc := strings.TrimSpace(r.FormValue("desc"))

		board.Desc = newDesc
		if err := board.Save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, board *Board, action string) {
	if strings.EqualFold(action, "desc") {
		descHandler(w, r, board)
		return
	}

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
		cookie := &http.Cookie{Name: "who", Value: who}
		http.SetCookie(w, cookie)
		cookie = &http.Cookie{Name: "email", Value: email}
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
					board.ActivityCount++
					recordUpdated = true
					break
				}
			}
		}

		if !recordUpdated {
			// If we didn't find an existing record to update then add a new one.
			board.Records = append(board.Records, &Record{Who: who, When: when, Email: email, Score: score})
			board.ActivityCount++
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

		var vm = &BoardViewModel{Board: board, TeamBoards: GetTeamBoardNames(board.Team)}
		cookie, _ := r.Cookie("who")
		if cookie != nil {
			vm.Who = cookie.Value
		}
		cookie, _ = r.Cookie("email")
		if cookie != nil {
			vm.Email = cookie.Value
		}

		if err := views.ExecuteTemplate(w, "board.gohtml", vm); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func init() {

	// Register helper functions to make them available in the templates.
	funcMap := template.FuncMap{
		"formatTime": formatTime,
		"takeTop":    takeTop,
		"skipTop":    skipTop,
		"equals": func(a, b string) bool {
			return strings.EqualFold(a, b)
		},
	}

	// Compile our template.
	views = template.Must(template.New("board.gohtml").Funcs(funcMap).ParseFiles("./views/board.gohtml", "./views/home.gohtml"))

	// Register our HTTP handlers.
	http.HandleFunc("/", makeBoardHandler(viewHandler))

	// ... including a static file server.
	http.Handle("/static/", http.FileServer(http.Dir("./")))
}

package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// View model for the scoreboard page.
type BoardViewModel struct {
	Board      *Board
	TeamBoards []string
	Who        string
	Email      string
}

// Gets a list of all the given team's boards.
func GetTeamBoardNames(team string) []string {
	// TODO: Handle the possible error here.
	files, _ := ioutil.ReadDir("./data/" + team)

	boardNames := make([]string, 0, len(files))

	for i := range files {
		if !files[i].IsDir() {
			filename := files[i].Name()
			boardNames = append(boardNames, filename[0:len(filename)-len(filepath.Ext(filename))])
		}
	}

	return boardNames
}

type Board struct {
	Team          string  `json:"team"`
	Name          string  `json:"name"`
	Desc          string  `json:"desc"`
	Records       Records `json:"records"`
	ActivityCount int64   `json:"activityCount"`
}

// Create a slice type for our Records so that we can implement
// the sort logic for it.
type Records []*Record

// Add the functions required for sorting our Records type.
func (s Records) Len() int           { return len(s) }
func (s Records) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Records) Less(i, j int) bool { return s[i].Score > s[j].Score }

type Record struct {
	Who   string    `json:"who"`
	Email string    `json:"email"`
	When  time.Time `json:"when"`
	Score int64     `json:"score"`
}

func (r *Record) GetAvatarUrl() string {
	h := md5.New()
	io.WriteString(h, strings.ToLower(r.Email))
	return "http://www.gravatar.com/avatar/" + fmt.Sprintf("%x", h.Sum(nil)) + "?r=R&s=60&d=mm"
}

// Gets the file path where the given team and board should be saved.
func getSavePath(team string, board string) string {
	dir := "./data/" + team
	// TODO: Handle the potential error here.
	os.MkdirAll(dir, 0777)
	return dir + "/" + board + ".json"
}

// Gets the file path to save the serialized Board instance to.
func (b *Board) getSavePath() string {
	return getSavePath(b.Team, b.Name)
}

// Saves the Board by serializing it as JSON to disk.
func (b *Board) Save() error {
	boardJson, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(b.getSavePath(), boardJson, 0600)
}

// Gets the Board for the given team and board name by deserializing
// it from JSON on disk. If the JSON file doesn't exist then we will
// return a new Board instance.
func LoadBoard(team string, name string) (*Board, error) {
	boardJson, err := ioutil.ReadFile(getSavePath(team, name))
	var b Board
	if err != nil {
		if os.IsNotExist(err) {
			// The requested board doesn't exist.
			// Create it.
			b = Board{Team: team, Name: name, Records: make([]*Record, 0, 1)}
		} else {
			return nil, err
		}
	} else {
		err = json.Unmarshal(boardJson, &b)
		if err != nil {
			return nil, err
		}
	}

	return &b, nil
}

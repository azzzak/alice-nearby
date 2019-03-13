package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"

	_ "github.com/azzzak/alice-nearby/statik"
	flags "github.com/jessevdk/go-flags"
	"github.com/rakyll/statik/fs"
)

// Opts with cli commands and flags
type Opts struct {
	Webhook string `long:"webhook" env:"WEBHOOK" required:"true" description:"url to webhook"`
	Port    string `long:"port" short:"p" env:"PORT" required:"false" default:"2345" description:"port to web frontend"`
}

type question struct {
	Message  int    `json:"message_id"`
	Session  string `json:"session_id"`
	Skill    string `json:"skill_id"`
	User     string `json:"user_id"`
	Question string `json:"question"`
}

type packet struct {
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

var opts Opts

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		// fmt.Println(err)
		os.Exit(1)
	}

	// fmt.Println(opts)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/app", listen)
	http.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	// http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(fmt.Sprintf(":%s", opts.Port), nil)
}

func makeRequest(p question) *Request {
	r := Request{}

	r.Meta.ClientID = "ru.yandex.searchplugin/7.16 (none none; android 4.4.2)"
	r.Meta.Locale = "ru-RU"
	r.Meta.Timezone = "UTC"

	r.Request.Command = p.Question
	r.Request.Utterance = p.Question
	r.Request.Type = CommandTypeVoice

	re := regexp.MustCompile(`[\s,.]`)
	tmp := re.Split(p.Question, -1)
	var tokens = []string{}
	for _, v := range tmp {
		if v != "" && v != "-" {
			tokens = append(tokens, v)
		}
	}
	r.Request.NLU.Tokens = tokens
	r.Request.NLU.Entities = []struct {
		Tokens struct {
			Start int `json:"start"`
			End   int `json:"end"`
		} `json:"tokens"`
		Type  string           `json:"type"`
		Value *json.RawMessage `json:"value"`
	}{}

	r.Session.MessageID = p.Message

	if p.Session != "" {
		r.Session.SessionID = p.Session
	} else {
		r.Session.SessionID = generateID()
	}

	if p.Skill != "" {
		r.Session.SkillID = p.Skill
	} else {
		r.Session.SkillID = generateID()
	}

	if p.User != "" {
		r.Session.UserID = p.User
	} else {
		r.Session.UserID = generateID()
	}

	r.Version = "1.0"
	return &r
}

func generateID() string {
	str := "0123456789abcdef"
	var buf bytes.Buffer
	for i := 0; i < 32; i++ {
		if i == 8 || i == 16 || i == 24 {
			buf.Write([]byte("-"))
		}
		rnd := rand.Intn(len(str))
		buf.Write([]byte{str[rnd]})
	}
	return buf.String()
}

func useWebhook(u *Request) *Response {
	var client = &http.Client{
		Timeout: time.Second * 3,
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(u)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Post(opts.Webhook, "application/json; charset=utf-8", b)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	var body Response
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Printf("%s -> %s\n", u.Request.Command, body.Response.Text)

	return &body
}

func listen(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var p question
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

	// fmt.Printf("%+v\n", p)

	request := makeRequest(p)
	resp := useWebhook(request)

	if resp == nil {
		return
	}

	profile := packet{
		request,
		resp,
	}

	js, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

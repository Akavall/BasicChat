package main 

import (
	"text/template"
	"net/http"
	"log"
	"fmt"
	"io"
	"io/ioutil"
	"encoding/json"
	"strings"
)

type ChatWindow struct {
	// Should be protected by Mutex
	messages []string
	user_name []string 
}

func handler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("chat.html")
	if err != nil {
		log.Fatalln(err)
	}
	err = t.Execute(w, nil) 
	if err != nil {
		log.Fatalln(err)
	}

}

// Post request
func (cw *ChatWindow) send(w http.ResponseWriter, r *http.Request) {

	log.Println("Starting")

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	
	err = r.Body.Close()
	if err != nil {
		panic(err)
	}

	message_info := map[string]string {}

	err = json.Unmarshal(body, &message_info)
	if err != nil {
		panic(err)
	}

	message := message_info["message"]
	user_name := message_info["user_name"]

	log.Printf("Sending message: %s: %s", user_name, message)

	full_entry := fmt.Sprintf("%s: %s", user_name, message)
	(*cw).messages = append((*cw).messages, full_entry)
}

// Get request
func (cw *ChatWindow) receive_all(w http.ResponseWriter, r *http.Request) {
	var messages_str string
	if len((*cw).messages) >= 10 {
		messages_str = strings.Join((*cw).messages[len((*cw).messages) - 10 :], "<br>")
	} else {
		messages_str = strings.Join((*cw).messages, "<br>")
	}

	fmt.Fprint(w, messages_str)
}


func main() {

	chat_window := &ChatWindow {}

	http.HandleFunc("/chat", handler)
	http.HandleFunc("/send", chat_window.send)
	http.HandleFunc("/receive_all", chat_window.receive_all)

	ip_address := "0.0.0.0:8090"
	fmt.Println(ip_address)
	
	error := http.ListenAndServe(ip_address, nil)
	if error != nil {
		log.Fatalln(error)
	}
}

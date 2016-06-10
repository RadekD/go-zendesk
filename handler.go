package zendesk

import (
	"encoding/json"
	v "github.com/asaskevich/govalidator"
	"log"
	"net/http"
	"strings"
)

//Strategy defines what to do after user sends message
type Strategy int

const (
	//AlwaysNewTicket creates new ticket every time user sends message
	AlwaysNewTicket Strategy = iota
	//AddToLastTicket adds message to last not resolved ticket
	AddToLastTicket
)

//ContactHandler is a helper for sending messages
type ContactHandler struct {
	*Auth
	BeforeCreateTicket func(r *http.Request, ticket *NewTicket) *NewTicket
	Strategy           Strategy
}

type response struct {
	Success bool
	Errors  map[string]string `json:"Errors,omitempty"`
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	rJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("Error while marshaling response", err)
		return
	}
	w.Write(rJSON)
	return
}

func (z ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		panic("Wrong method")
	}
	w.Header().Add("Content-Type", "application/json;charset=utf-8")

	resp := response{Success: false, Errors: make(map[string]string)}

	var err error
	errors := make(map[string]string)

	email := r.FormValue("Email")
	subject := r.FormValue("Subject")
	content := r.FormValue("Message")

	if email == "" {
		resp.Errors["Email"] = "empty"
	}
	if email != "" && !v.IsEmail(email) {
		resp.Errors["Email"] = "invalid"
	}
	if subject == "" {
		resp.Errors["Subject"] = "empty"
	}
	if content == "" {
		resp.Errors["Content"] = "empty"
	}
	if len(resp.Errors) > 0 {
		sendJSON(w, resp)
		return
	}

	users, _ := z.SearchUser(email)
	var user User
	if len(users) == 0 {
		user, err = z.CreateUser(NewUser{
			Name:  strings.Split(email, "@")[0],
			Email: email,
			Roles: "end-user",
		})

		if err != nil {
			resp.Errors["User"] = err.Error()
			sendJSON(w, resp)
			return
		}
	} else {
		user = users[0]
	}

	if z.Strategy == AddToLastTicket {
		requests, _ := z.ListOpenRequests(user)
		if len(requests) > 0 {
			_, err = z.AddCommentToTicket(requests[0].ID, &Comment{Body: content, AuthorID: user.ID})
			if err != nil {
				resp.Errors["Internal"] = err.Error()
				sendJSON(w, resp)
				return
			}

			resp.Success = true
			sendJSON(w, resp)
			return
		}
	}

	newTicket := &NewTicket{
		Subject:     subject,
		RequesterID: UserID(user.ID),
		Status:      "new",
		Comment: Comment{
			Body: content,
		},
	}
	if z.BeforeCreateTicket != nil {
		newTicket = z.BeforeCreateTicket(r, newTicket)
	}
	_, err = z.CreateTicket(newTicket)
	if err != nil {
		errors["Internal"] = err.Error()
		sendJSON(w, resp)
		return
	}
	resp.Success = true
	sendJSON(w, resp)
	return
}

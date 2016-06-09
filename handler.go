package zendesk

import (
	v "github.com/asaskevich/govalidator"
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

//Handler is a helper for sending messages
type Handler struct {
	*Auth

	GetFunc     http.HandlerFunc
	ErrorFunc   func(w http.ResponseWriter, r *http.Request, errors map[string]string)
	SuccessFunc func(w http.ResponseWriter, r *http.Request, ticket Ticket)

	CustomFieldsFunc func(r *http.Request) CustomFields
	TagsFunc         func(r *http.Request) []string
	ExternalIDFunc   func(r *http.Request) *uint64

	Strategy Strategy
}

func (z Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		z.GetFunc(w, r)
		return
	}
	var err error
	errors := make(map[string]string)

	email := r.FormValue("Email")
	subject := r.FormValue("Subject")
	content := r.FormValue("Message")

	if email == "" {
		errors["Email"] = "empty"
	}
	if !v.IsEmail(email) {
		errors["Email"] = "invalid"
	}
	if subject == "" {
		errors["Subject"] = "empty"
	}
	if content == "" {
		errors["Content"] = "empty"
	}
	if len(errors) > 0 {
		z.ErrorFunc(w, r, errors)
		return
	}

	users, _ := z.SearchUser(email)
	var user User
	if len(users) == 0 {
		var externalID *uint64
		if z.ExternalIDFunc != nil {
			externalID = z.ExternalIDFunc(r)
		}

		user, err = z.CreateUser(NewUser{
			Name:       strings.Split(email, "@")[0],
			Email:      email,
			ExternalID: externalID,
			Verified:   true,
			Roles:      "end-user",
		})

		if err != nil {
			errors["User"] = err.Error()
			z.ErrorFunc(w, r, errors)
			return
		}
	} else {
		user = users[0]
	}

	if z.Strategy == AddToLastTicket {
		requests, _ = z.ListOpenRequests(user)
		if len(requests) > 0 {
			ticket, err = z.AddCommentToTicket(requests[0].ID, &Comment{Body: content, AuthorID: user.ID})
			if err != nil {
				errors["Internal"] = err.Error()
				z.ErrorFunc(w, r, errors)
				return
			}
			z.SuccessFunc(w, r, *ticket)
			return
		}
	}

	var customFields CustomFields
	if z.CustomFieldsFunc != nil {
		customFields = z.CustomFieldsFunc(r)
	}

	var tags []string
	if z.TagsFunc != nil {
		tags = z.TagsFunc(r)
	}

	newTicket := &NewTicket{
		Subject:     subject,
		RequesterID: UserID(user.ID),
		Status:      "new",
		Comment: Comment{
			Body: content,
		},
		Tags:         tags,
		CustomFields: customFields,
	}
	ticket, err := z.CreateTicket(newTicket)

	if err != nil {
		errors["Internal"] = err.Error()
		z.ErrorFunc(w, r, errors)
		return
	}
	z.SuccessFunc(w, r, *ticket)
}

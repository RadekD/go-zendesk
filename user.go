package zendesk

import (
	"encoding/json"
	"fmt"
	"time"
)

type UserID uint64

type User struct {
	ID          UserID  `json:"id,omitempty"`
	URL         string  `json:"url,omitempty"`
	Name        string  `json:"name,omitempty"`
	ExternalID  *uint64 `json:"external_id,omitempty"`
	Alias       string  `json:"alias,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
	Active      bool    `json:"active,omitempty"`
	Verified    bool    `json:"verified,omitempty"`
	LastLoginAt string  `json:"last_login_at,omitempty"`
	Email       string  `json:"email,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Details     string  `json:"details,omitempty"`
	Suspended   bool    `json:"suspended,omitempty"`
}

type Users struct {
	Users []User `json:"users"`
}

//SearchUser a search
func (auth Auth) SearchUser(query string) ([]User, error) {
	data, err := api(auth, "GET", "/users/search.json?query="+query, "")
	if err != nil {
		return nil, err
	}

	users := &Users{}
	err = json.Unmarshal(data, users)

	return users.Users, err
}

//NewUser represents new user
type NewUser struct {
	Email      string  `json:"email,omitempty"`
	Name       string  `json:"name,omitempty"`
	Roles      string  `json:"roles,omitempty"`
	ExternalID *uint64 `json:"external_id,omitempty"`
	Verified   bool    `json:"verified,omitempty"`
}

//CreateUser creates user
func (auth Auth) CreateUser(user NewUser) (User, error) {
	newUser, err := json.Marshal(map[string]interface{}{"user": user})
	if err != nil {
		return User{}, err
	}

	data, err := api(auth, "POST", "/users.json", string(newUser))
	if err != nil {
		return User{}, err
	}
	createdUser := &struct {
		User User `json:"user"`
	}{}

	err = json.Unmarshal(data, createdUser)
	return createdUser.User, err
}

type Request struct {
	ID              uint64       `json:"id"`
	AssigneeID      uint64       `json:"assignee_id"`
	CanBeSolvedByMe bool         `json:"can_be_solved_by_me"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
	CustomFields    CustomFields `json:"custom_fields,omitempty"`
	RequesterID     int64        `json:"requester_id,omitempty"`
	Description     string       `json:"description,omitempty"`
	DueAt           *time.Time   `json:"due_at,omitempty"`
	OrganizationID  int64        `json:"organization_id,omitempty"`
	Priority        string       `json:"priority,omitempty"`
	Status          string       `json:"status,omitempty"`
	Subject         string       `json:"subject,omitempty"`
	Type            string       `json:"type,omitempty"`
	URL             string       `json:"url"`
	Via             struct {
		Channel string `json:"channel,omitempty"`
	} `json:"via"`
	CollaboratorIds []*uint64 `json:"collaborator_ids"`
}
type RequestsSearch struct {
	Count    uint       `json:"count"`
	NextPage *string    `json:"next_page"`
	PrevPage *string    `json:"previous_page"`
	Requests []*Request `json:"requests"`
}

func (auth Auth) ListOpenRequests(user User) ([]*Request, error) {
	data, err := api(auth, "GET", fmt.Sprintf("/users/%d/requests.json?sort_order=desc", user.ID), "")
	if err != nil {
		return nil, err
	}
	sr := RequestsSearch{}
	err = json.Unmarshal(data, &sr)
	return sr.Requests, err
}

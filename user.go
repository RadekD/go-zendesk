package zendesk

import "encoding/json"
import "fmt"
import "time"

type UserId uint64
type ExternalId *uint64
type Email string

type User struct {
	Id                  UserId        `json:"id,omitempty"`
	Url                 string        `json:"url,omitempty"`
	Name                string        `json:"name,omitempty"`
	ExternalId          ExternalId    `json:"external_id,omitempty"`
	Alias               string        `json:"alias,omitempty"`
	CreatedAt           string        `json:"created_at,omitempty"`
	UpdatedAt           string        `json:"updated_at,omitempty"`
	Active              bool          `json:"active,omitempty"`
	Verified            bool          `json:"verified,omitempty"`
	Shared              bool          `json:"shared,omitempty"`
	SharedAgent         bool          `json:"shared_agent,omitempty"`
	Locale              string        `json:"locale,omitempty"`
	LocaleId            uint32        `json:"locale_id,omitempty"`
	TimeZone            string        `json:"time_zone,omitempty"`
	LastLoginAt         string        `json:"last_login_at,omitempty"`
	Email               Email         `json:"email,omitempty"`
	Phone               string        `json:"phone,omitempty"`
	Signature           string        `json:"signature,omitempty"`
	Details             string        `json:"details,omitempty"`
	Notes               string        `json:"notes,omitempty"`
	OrganizationId      uint32        `json:"organization_id,omitempty"`
	Role                string        `json:"role,omitempty"`
	CustomerRoleId      uint32        `json:"custom_role_id,omitempty"`
	Moderator           bool          `json:"moderator,omitempty"`
	TicketRestriction   string        `json:"ticket_restriction,omitempty"`
	OnlyPrivateComments bool          `json:"only_private_comments,omitempty"`
	Tags                []string      `json:"tags,omitempty"`
	RestrictedAgent     bool          `json:"restricted_agent,omitempty"`
	Suspended           bool          `json:"suspended,omitempty"`
	Photo               []*UsersPhoto `json:"photo,omitempty"`
	UserFields          []*UserField  `json:"user_fields,omitempty"`
}
type UsersPhoto struct {
	Id          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ContentUrl  string `json:"content_url,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Size        int    `json:"size,omitempty"`
}
type UserField struct {
	UserDecimal  float32 `json:"user_decimal,omitempty"`
	UserDropdown string  `json:"user_dropdown,omitempty"`
	UserDate     string  `json:"user_date,omitempty"`
}

type Users struct {
	Users []User `json:"users"`
}

func (auth Auth) Search(query string) ([]User, error) {
	data, err := api(auth, "GET", "/users/search.json?query="+query, "")
	if err != nil {
		return nil, err
	}

	users := &Users{}
	json.Unmarshal(data, users)

	return users.Users, nil
}

func (auth Auth) UpdateUser(user *User, update map[string]interface{}) error {

	to_update, err := json.Marshal(map[string]interface{}{"user": update})
	if err != nil {
		return err
	}

	_, err = api(auth, "PUT", fmt.Sprintf("/users/%d.json", user.Id), string(to_update))
	return err

}

type NewUser struct {
	Email      Email      `json:"email,omitempty"`
	Name       string     `json:"name,omitempty"`
	Roles      string     `json:"roles,omitempty"`
	ExternalId ExternalId `json:"external_id,omitempty"`
	Verified   bool       `json:"verified,omitempty"`
}

func (auth Auth) CreateUser(user NewUser) (User, error) {
	new_user, err := json.Marshal(map[string]interface{}{"user": user})
	if err != nil {
		return User{}, err
	}

	data, err := api(auth, "POST", "/users.json", string(new_user))

	created_user := &struct {
		User User `json:"user"`
	}{}
	json.Unmarshal(data, created_user)

	return created_user.User, nil
}

type Request struct {
	Id              uint64       `json:"id"`
	AssigneeId      uint64       `json:"assignee_id"`
	CanBeSolvedByMe bool         `json:"can_be_solved_by_me"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
	CustomFields    CustomFields `json:"custom_fields,omitempty"`
	RequesterId     int64        `json:"requester_id,omitempty"`
	Description     string       `json:"description,omitempty"`
	DueAt           *time.Time   `json:"due_at,omitempty"`
	OrganizationId  int64        `json:"organization_id,omitempty"`
	Priority        Priority     `json:"priority,omitempty"`
	Status          Status       `json:"status,omitempty"`
	Subject         string       `json:"subject,omitempty"`
	Type            string       `json:"type,omitempty"`
	Url             string       `json:"url"`
	Via             Via          `json:"via"`
	CollaboratorIds []*uint64    `json:"collaborator_ids"`
}
type RequestsSearch struct {
	Count    uint       `json:"count"`
	NextPage *string    `json:"next_page"`
	PrevPage *string    `json:"previous_page"`
	Requests []*Request `json:"requests"`
}

func (auth Auth) ListOpenRequests(user User) []*Request {
	data, _ := api(auth, "GET", fmt.Sprintf("/users/%d/requests.json", user.Id), "")

	sr := RequestsSearch{}
	json.Unmarshal(data, &sr)
	return sr.Requests
}

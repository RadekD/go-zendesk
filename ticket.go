package zendesk

import (
	"encoding/json"
	"fmt"
	"time"
)

type CustomField struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}
type CustomFields []CustomField

//Ticket
type Ticket struct {
	ID              uint64     `json:"id,omitempty"`
	URL             string     `json:"url,omitempty"`
	ExternalID      string     `json:"external_id,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
	Type            string     `json:"type,omitempty"`
	Subject         string     `json:"subject,omitempty"`
	RawSubject      string     `json:"raw_subject,omitempty"`
	Description     string     `json:"description,omitempty"`
	Priority        string     `json:"priority,omitempty"`
	Status          string     `json:"status,omitempty"`
	Recipient       string     `json:"recipient,omitempty"`
	RequesterID     int64      `json:"requester_id,omitempty"`
	SubmitterID     int64      `json:"submitter_id,omitempty"`
	AssigneeID      int64      `json:"assignee_id,omitempty"`
	OrganizationID  int64      `json:"organization_id,omitempty"`
	GroupID         int64      `json:"group_id,omitempty"`
	CollaboratorIds []int64    `json:"collaborator_ids,omitempty"`
	ForumTopicID    int64      `json:"forum_topic_id,omitempty"`
	ProblemID       int64      `json:"problem_id,omitempty"`
	HasIncidents    bool       `json:"has_incidents,omitempty"`
	DueAt           *time.Time `json:"due_at,omitempty"`
	Tags            []string   `json:"tags,omitempty"`
	Via             struct {
		Channel string `json:"channel,omitempty"`
	} `json:"via,omitempty"`
	CustomFields       CustomFields `json:"custom_fields,omitempty"`
	SatisfactionRating struct {
		ID      int    `json:"id"`
		Score   string `json:"score"`
		Comment string `json:"comment"`
	} `json:"satisfaction_rating,omitempty"`
	SharingAgreementIds []int64 `json:"sharing_agreement_ids,omitempty"`
}

//NewTicket represents new tickets
type NewTicket struct {
	Subject             string       `json:"subject,omitempty"`
	Comment             Comment      `json:"comment,omitempty"`
	RequesterID         UserID       `json:"requester_id,omitempty"`
	SubmitterID         UserID       `json:"submitter_id,omitempty"`
	AssigneeID          UserID       `json:"assignee_id,omitempty"`
	GroupID             uint64       `json:"group_id,omitempty"`
	CollaboratorIDs     []UserID     `json:"collaborator_ids,omitempty"`
	Type                string       `json:"type,omitempty"`
	Priority            string       `json:"priority,omitempty"`
	Status              string       `json:"status,omitempty"`
	Tags                []string     `json:"tags,omitempty"`
	ExternalID          *uint64      `json:"external_id,omitempty"`
	ForumTopicID        *uint64      `json:"forum_topic_id,omitempty"`
	ProblemID           *uint64      `json:"problem_id,omitempty"`
	DueAt               *time.Time   `json:"due_at,omitempty"`
	CustomFields        CustomFields `json:"custom_fields,omitempty"`
	ViaFollowupSourceID *uint64      `json:"via_followup_source_id,omitempty"`
}

//Comment represents single message in ticket
type Comment struct {
	ID           uint64        `json:"id,omitempty"`
	Type         string        `json:"type,omitempty"`
	Body         string        `json:"body,omitempty"`
	Public       bool          `json:"public,omitempty"`
	CreatedAt    *time.Time    `json:"created_at,omitempty"`
	AuthorID     UserID        `json:"author_id,omitempty"`
	Attachements []Attachement `json:"attachements,omitempty"`
	Via          *struct {
		Channel string `json:"channel,omitempty"`
	} `json:"via,omitempty"`
}

//Attachement represents attachement in ticket
type Attachement struct {
	ID          uint64   `json:"id,omitempty"`
	Name        string   `json:"name"`
	ContentURL  string   `json:"content_url"`
	ContentType string   `json:"content_type"`
	Size        uint64   `json:"size"`
	Thumbnails  []string `json:"thumbnails"`
}

func returnTicketOrError(data []byte, err error) (*Ticket, error) {
	if err != nil {
		return nil, err
	}
	var toJSON struct {
		Ticket *Ticket `json:"ticket"`
	}
	err = json.Unmarshal(data, &toJSON)
	if err != nil {
		return nil, err
	}
	return toJSON.Ticket, nil
}

//CreateTicket creates new ticket
func (auth Auth) CreateTicket(ticket *NewTicket) (*Ticket, error) {
	ticketJSON, err := json.Marshal(map[string]*NewTicket{"ticket": ticket})
	if err != nil {
		return nil, err
	}
	data, err := api(auth, "POST", "/tickets.json", string(ticketJSON))
	return returnTicketOrError(data, err)
}

//AddCommentToTicket adds comment to existing ticket
func (auth Auth) AddCommentToTicket(ticketID uint64, comment *Comment) (*Ticket, error) {
	toUpdate := map[string]interface{}{
		"ticket": map[string]interface{}{
			"comment": comment,
		},
	}
	commentJSOM, err := json.Marshal(toUpdate)
	if err != nil {
		return nil, err
	}
	data, err := api(auth, "PUT", fmt.Sprintf("/tickets/%d.json", ticketID), string(commentJSOM))
	return returnTicketOrError(data, err)
}

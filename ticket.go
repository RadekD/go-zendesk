package zendesk

import "time"
import "encoding/json"
import "fmt"

type Via struct {
	Channel string `json:"channel,omitempty"`
}

type Tag string
type Tags []Tag

type CustomField struct {
	Id    int    `json:"id"`
	Value string `json:"value"`
}
type CustomFields []CustomField

type SatisfactionRating struct {
	Id      int    `json:"id"`
	Score   string `json:"score"`
	Comment string `json:"comment"`
}

type Ticket struct {
	Id                  uint64             `json:"id,omitempty"`
	Url                 string             `json:"url,omitempty"`
	ExternalId          string             `json:"external_id,omitempty"`
	CreatedAt           time.Time          `json:"created_at,omitempty"`
	UpdatedAt           time.Time          `json:"updated_at,omitempty"`
	Type                string             `json:"type,omitempty"`
	Subject             string             `json:"subject,omitempty"`
	RawSubject          string             `json:"raw_subject,omitempty"`
	Description         string             `json:"description,omitempty"`
	Priority            Priority           `json:"priority,omitempty"`
	Status              Status             `json:"status,omitempty"`
	Recipient           string             `json:"recipient,omitempty"`
	RequesterId         int64              `json:"requester_id,omitempty"`
	SubmitterId         int64              `json:"submitter_id,omitempty"`
	AssigneeId          int64              `json:"assignee_id,omitempty"`
	OrganizationId      int64              `json:"organization_id,omitempty"`
	GroupId             int64              `json:"group_id,omitempty"`
	CollaboratorIds     []int64            `json:"collaborator_ids,omitempty"`
	ForumTopicId        int64              `json:"forum_topic_id,omitempty"`
	ProblemId           int64              `json:"problem_id,omitempty"`
	HasIncidents        bool               `json:"has_incidents,omitempty"`
	DueAt               *time.Time         `json:"due_at,omitempty"`
	Tags                Tags               `json:"tags,omitempty"`
	Via                 Via                `json:"via,omitempty"`
	CustomFields        CustomFields       `json:"custom_fields,omitempty"`
	SatisfactionRating  SatisfactionRating `json:"satisfaction_rating,omitempty"`
	SharingAgreementIds []int64            `json:"sharing_agreement_ids,omitempty"`
}

type NewTicket struct {
	Subject             string       `json:"subject,omitempty"`
	Comment             Comment      `json:"comment,omitempty"`
	RequesterId         UserId       `json:"requester_id,omitempty"`
	SubmitterId         UserId       `json:"submitter_id,omitempty"`
	AssigneeId          UserId       `json:"assignee_id,omitempty"`
	GroupId             uint64       `json:"group_id,omitempty"`
	CollaboratorIds     []UserId     `json:"collaborator_ids,omitempty"`
	Type                string       `json:"type,omitempty"`
	Priority            Priority     `json:"priority,omitempty"`
	Status              Status       `json:"status,omitempty"`
	Tags                Tags         `json:"tags,omitempty"`
	ExternalId          *uint64      `json:"external_id,omitempty"`
	ForumTopicId        *uint64      `json:"forum_topic_id,omitempty"`
	ProblemId           *uint64      `json:"problem_id,omitempty"`
	DueAt               *time.Time   `json:"due_at,omitempty"`
	CustomFields        CustomFields `json:"custom_fields,omitempty"`
	ViaFollowupSourceId *uint64      `json:"via_followup_source_id,omitempty"`
}

type Comment struct {
	Id           uint64        `json:"id,omitempty"`
	Type         string        `json:"type,omitempty"`
	Body         string        `json:"body,omitempty"`
	Public       bool          `json:"public,omitempty"`
	CreatedAt    *time.Time    `json:"created_at,omitempty"`
	AuthorId     uint64        `json:"author_id,omitempty"`
	Attachements []Attachement `json:"attachements,omitempty"`
	Via          *Via          `json:"via,omitempty"`
}
type Attachement struct {
	Id          uint64   `json:"id,omitempty"`
	Name        string   `json:"name"`
	ContentUrl  string   `json:"content_url"`
	ContentType string   `json:"content_type"`
	Size        uint64   `json:"size"`
	Thumbnails  []string `json:"thumbnails"`
}

func (auth Auth) CreateTicket(ticket *NewTicket) (*Ticket, error) {

	ticket_json, err := json.Marshal(map[string]*NewTicket{"ticket": ticket})
	if err != nil {
		return nil, err
	}
	//return nil, nil
	data, err := api(auth, "POST", "/tickets.json", string(ticket_json))

	var to_json struct {
		Ticket *Ticket `json:"ticket"`
	}
	json.Unmarshal(data, &to_json)
	return to_json.Ticket, nil
}

func (auth Auth) AddCommentToTicket(ticket_id uint64, comment *Comment) (*Ticket, error) {
	to_update := map[string]interface{}{
		"ticket": map[string]interface{}{
			"comment": comment,
		},
	}
	comment_json, _ := json.Marshal(to_update)
	api(auth, "PUT", fmt.Sprintf("/tickets/%d.json", ticket_id), string(comment_json))
	return nil, nil
}

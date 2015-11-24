package zendesk

import "fmt"

func main() {
	auth := Auth{"LOGIN", "KEY", "ADDR", true}
	e, _ := auth.Search("email")
	var user User
	if len(e) == 0 {
		user, _ = auth.CreateUser(NewUser{
			Name:     "name",
			Email:    Email("email"),
			Verified: true,
			Roles:    "end-user",
		})
		fmt.Printf("\n%s\n", user)
	} else {
		user = e[0]
	}

	/*n := &NewTicket{
		Subject:     "test",
		RequesterId: UserId(user.Id),
		Status:      New,
		Comment: Comment{
			Body: "test",
		},
	}
	auth.CreateTicket(n)*/

	requests := auth.ListOpenRequests(user)
	if len(requests) > 0 {
		auth.AddCommentToTicket(requests[0].Id, &Comment{Body: "test44", AuthorId: uint64(user.Id)})
	}
}

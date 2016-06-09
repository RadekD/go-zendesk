// +build ignore

package main

import zendesk "github.com/RadekD/go-zendesk"
import "github.com/gorilla/mux"
import "fmt"
import "net/http"

func main() {
	r := mux.NewRouter()

	auth := zendesk.NewAuth("radek@dejnek.pl", "GySziP9xbK0OG0JPgxzDl6Zy8FO74mzSMy383Vcc", "dejnek.zendesk.com", true)

	handler := zendesk.Handler{
		Auth: auth,

		GetFunc:     handleForm,
		ErrorFunc:   handleError,
		SuccessFunc: handleSuccess,

		CustomFieldsFunc: func(r *http.Request) zendesk.CustomFields {
			return zendesk.CustomFields{
				0: zendesk.CustomField{
					ID:    27226329,
					Value: "eeeeeeeeee",
				},
			}
		},
		TagsFunc: func(r *http.Request) []string {
			return []string{"super", "tagi", "bulwo"}
		},

		Strategy: zendesk.AddToLastTicket,
	}
	r.Handle("/contact", handler)

	fmt.Println(http.ListenAndServe("127.0.0.1:8088", r))

}

func handleError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	fmt.Fprintf(w, "%s", errors)
}
func handleSuccess(w http.ResponseWriter, r *http.Request, ticket zendesk.Ticket) {
	fmt.Fprintf(w, "udało się %s", ticket.URL)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html;charset=utf-8")
	fmt.Fprintf(w, `<form class="ui inverted form contact" action="/contact" method="POST">
            <div class="ui success message">Wiadomość została wysłana</div>
            <div class="field">
                <label>Twój adres e-mail</label>
                <div class="field">
                    <input type="text" name="Email">
                </div>
            </div>
            <div class="field">
                <label>Temat</label>
                <div class="field">
                    <input type="text" name="Subject">
                </div>
            </div>
            <div class="field">
                <label>Wiadomość</label>
                <div class="field">
                    <textarea name="Message"></textarea>
                </div>
            </div>
            <button type="submit" class="ui positive button" tabindex="0"><i class="ui icon mail"></i>Wyślij wiadomość</button>
        </form>`)
}

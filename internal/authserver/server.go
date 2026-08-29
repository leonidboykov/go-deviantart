package authserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2/authhandler"
)

func AuthHandler(callbackURL string) authhandler.AuthorizationHandler {
	return func(authCodeURL string) (code string, state string, err error) {
		fmt.Println("Visit the URL for the auth dialog:", authCodeURL)

		uri, err := url.Parse(callbackURL)
		if err != nil {
			return "", "", fmt.Errorf("parse callback url: %w", err)
		}

		codeChan := make(chan string)
		stateChan := make(chan string)

		mux := http.NewServeMux()
		mux.HandleFunc(uri.Path, func(w http.ResponseWriter, r *http.Request) {
			// TODO: Cleanup code.
			log.Println("Code:", r.FormValue("code"), "State:", r.FormValue("state"))

			// Auto close a new tab.
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<script>window.close();</script>`))

			codeChan <- r.FormValue("code")
			stateChan <- r.FormValue("state")
		})

		go func() {
			if err := http.ListenAndServe(uri.Host, mux); err != nil {
				log.Fatalln(err)
			}
		}()

		// TODO: Stop http server.
		return <-codeChan, <-stateChan, err
	}
}

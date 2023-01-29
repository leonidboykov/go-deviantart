package authserver

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type AuthServer struct {
	conf  *oauth2.Config
	state string
}

func NewAuthServer(conf *oauth2.Config) *AuthServer {
	return &AuthServer{
		conf:  conf,
		state: RandString(32),
	}
}

func (a *AuthServer) ListenToken() (string, error) {
	url, err := url.Parse(a.conf.RedirectURL)
	if err != nil {
		return "", fmt.Errorf("unable to parse redirectURL: %w", err)
	}

	codeChan := make(chan string)

	mux := http.NewServeMux()
	mux.HandleFunc(url.Path, func(w http.ResponseWriter, r *http.Request) {
		// TODO: Cleanup code.
		log.Println("Code:", r.FormValue("code"), "State:", r.FormValue("state"))

		// Auto close a new tab.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<script>window.close();</script>`))

		codeChan <- r.FormValue("code")
	})

	go func() {
		if err := http.ListenAndServe(url.Host, mux); err != nil {
			//return "", err
			log.Fatalln(err)
		}
	}()

	return <-codeChan, nil
}

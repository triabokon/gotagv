package server

import (
	"fmt"
	"net/http"

	"github.com/triabokon/gotagv/auth"
)

func (s *Server) SetRoutes() {
	s.router.HandleFunc("/", auth.HandleAuth(s.HelloHandler))
	s.router.HandleFunc("/signup", s.SignUp)
	s.router.HandleFunc("/signin", s.SignIn)
}

func (s *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	w.Write([]byte(fmt.Sprintf("Hello, %s!", username)))
}

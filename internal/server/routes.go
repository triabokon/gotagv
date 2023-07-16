package server

import (
	"fmt"
	"net/http"

	"github.com/triabokon/gotagv/internal/auth"
)

func (s *Server) SetRoutes() {
	// todo: make simple health check
	s.router.HandleFunc("/", s.auth.HandleAuth(s.HelloHandler))

	s.router.HandleFunc("/signup", s.SignUp)
	s.router.HandleFunc("/signin", s.SignIn)

	s.router.HandleFunc("/videos/add", s.auth.HandleAuth(s.CreateVideo))
	s.router.HandleFunc("/videos", s.auth.HandleAuth(s.ListVideos))
	s.router.HandleFunc(fmt.Sprintf("/videos/delete/{%s}", entityIDKey), s.auth.HandleAuth(s.DeleteVideo))

	s.router.HandleFunc("/annotations/add", s.auth.HandleAuth(s.CreateAnnotation))
	s.router.HandleFunc("/annotations", s.auth.HandleAuth(s.ListAnnotations))
	s.router.HandleFunc(
		fmt.Sprintf("/annotations/update/{%s}", entityIDKey),
		s.auth.HandleAuth(s.UpdateAnnotation),
	)
	s.router.HandleFunc(
		fmt.Sprintf("/annotations/delete/{%s}", entityIDKey),
		s.auth.HandleAuth(s.DeleteAnnotation),
	)
}

func (s *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		s.ErrorResponse(w, fmt.Errorf("can't access user id"), http.StatusUnauthorized)
		return
	}
	s.SuccessResponse(w, fmt.Sprintf("Hello, %s!", userID))
}

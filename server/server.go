package main

import (
	"fmt"
	"github.com/cppforlife/checkman_jenkins_proxy/server/storer"
	"io"
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr   string
	storer storer.Storer
	server *http.Server
	logger *log.Logger
}

func NewServer(addr string, storer storer.Storer, logger *log.Logger) *Server {
	return &Server{
		addr:   addr,
		storer: storer,
		logger: logger,
	}
}

func (s *Server) ListenAndServe() error {
	s.server = &http.Server{
		Addr:           s.addr,
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.logger.Printf("server.listen-and-serve.run addr=%s\n", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) ServeHTTP(respWriter http.ResponseWriter, req *http.Request) {
	key := req.URL.Path[1:] // strip '/'

	switch req.Method {
	case "PUT":
		defer req.Body.Close()

		err := s.storer.Put(key, req.Body)
		if err != nil {
			s.logger.Printf("server.serve-http.put.fail url=%v err=%v\n", req.URL, err)
			s.respondWithError(500, respWriter)
		}

	case "GET":
		content, err := s.storer.Get(key)
		if err != nil {
			s.logger.Printf("server.serve-http.get.fail url=%v err=%v\n", req.URL, err)
			s.respondWithError(500, respWriter)
			return
		}

		if content == nil {
			s.logger.Printf("server.serve-http.get.no-content url=%v\n", req.URL)
			s.respondWithError(404, respWriter)
			return
		}

		defer content.Close()

		written, err := io.Copy(respWriter, content)
		if err != nil {
			s.logger.Printf("server.serve-http.get.copy.fail url=%v written=%d err=%v\n", req.URL, written, err)
			// cannot send error back since writing HTTP response failed
		}

	case "DELETE":
		err := s.storer.Delete(key)
		if err != nil {
			s.logger.Printf("server.serve-http.delete.fail url=%v err=%v\n", req.URL, err)
			s.respondWithError(500, respWriter)
		}

	default:
		s.logger.Printf("server.serve-http.unknown.fail url=%v method=%s\n", req.URL, req.Method)
		s.respondWithError(501, respWriter) // 501 Not Implemented
	}
}

func (s *Server) respondWithError(code int, respWriter http.ResponseWriter) {
	s.logger.Printf("server.respond-with-error code=%d\n", code)
	body := fmt.Sprintf("%d %s", code, http.StatusText(code))
	http.Error(respWriter, body, code)
}

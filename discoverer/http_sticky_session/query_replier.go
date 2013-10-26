package http_sticky_session

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type QueryReplier struct {
	uuid   string
	logger *log.Logger
}

func NewQueryReplier(uuid string, logger *log.Logger) *QueryReplier {
	return &QueryReplier{
		uuid:   uuid,
		logger: logger,
	}
}

// net/http.Handler
func (qr *QueryReplier) ServeHTTP(
	respWriter http.ResponseWriter, req *http.Request) {

	// Only add session cookie if one is given
	for _, cookie := range req.Cookies() {
		if cookie.Name == SessionKey {
			http.SetCookie(respWriter, &http.Cookie{
				Name:  SessionKey,
				Value: "any-value",
				Path:  "/",
			})
			break
		}
	}

	resp := map[string]string{"uuid": qr.uuid}
	bytes, err := json.Marshal(resp)
	if err != nil {
		qr.logger.Printf("http-sticky-session.query-replier.serve-http.fail err=%v\n", err)
		qr.respondWithError(500, respWriter)
	}

	respWriter.Write(bytes)
}

func (qr *QueryReplier) respondWithError(code int, respWriter http.ResponseWriter) {
	qr.logger.Printf("http-sticky-session.query-replier.respond-with-error code=%d\n", code)
	body := fmt.Sprintf("%d %s", code, http.StatusText(code))
	http.Error(respWriter, body, code)
}

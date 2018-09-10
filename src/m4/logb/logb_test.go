package logb

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogb(t *testing.T) {

	err := SetLogPath("./logs/")

	if err != nil {
		t.Error(err)
	}

	r, err := http.NewRequest("GET", "https://www.logb.com/", nil)
	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(http.HandlerFunc(testHandler))
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Error("Error Code: ", w.Code)
	}

}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

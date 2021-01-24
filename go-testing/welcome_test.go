package go_testing

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWelcome_name(t *testing.T) {
	resp, err := http.Get("http://localhost:3999/welcome?name=Frank")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		t.Fatal(err)
	}
	if v := doc.Find("h1.header-name span.name").Text(); v != "Frank" {
		t.Fatalf("expected markup to contain 'Frank', got '%s'", v)
	}
}

func TestWelcome_name_JSON(t *testing.T) {
	resp, err := http.Get("http://localhost:3999/welcome.json?name=Frank")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var dst struct{ Salutation string }
	if err := json.NewDecoder(resp.Body).Decode(&dst); err != nil {
		t.Fatal(err)
	}
	if dst.Salutation != "Hello Frank!" {
		t.Fatalf("expected 'Hello Frank!', got '%s'", dst.Salutation)
	}
}

func setup() *httptest.Server {
	return httptest.NewServer(nil)
}

func teardown(s *httptest.Server) {
	s.Close()
}

func TestWelcome_name_setup_teardown(t *testing.T) {
	srv := setup()

	url := fmt.Sprintf("%s/welcome.json?name=Frank", srv.URL)
	resp, err := http.Get(url)
	// verify errors & run assertions as usual
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)

	teardown(srv)
}

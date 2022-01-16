package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/kenshin579/tutorials-go/go-oauth2/mingrammer/auth"
	"golang.org/x/oauth2"
)

const (
	Templates = "/Users/user/GolandProjects/tutorials-go/go-oauth2/mingrammer/templates"
)

var store = sessions.NewCookieStore([]byte("secret"))

func main() {
	http.HandleFunc("/", RenderMainView)
	http.HandleFunc("/auth", RenderAuthView)
	http.HandleFunc("/auth/callback/google", Authenticate)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, _ := template.ParseFiles(name)
	tmpl.Execute(w, data)
}

// 메인 뷰 핸들러
func RenderMainView(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, Templates+"/main.html", nil)
}

// 랜덤 state 값을 가진 구글 로그인 링크를 렌더링 해주는 뷰 핸들러
// 랜덤 state는 유저를 식별하는 용도로 사용된다
func RenderAuthView(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options = &sessions.Options{
		Path:   "/auth",
		MaxAge: 300,
	}
	state := auth.RandToken()
	fmt.Printf("state: %+v\n", state)
	session.Values["state"] = state
	session.Save(r, w)
	url := auth.GetLoginURL(state)
	fmt.Println("url: ", url)
	RenderTemplate(w, Templates+"/auth.html", url)
}

// Google OAuth 인증 콜백 핸들러
func Authenticate(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	state := session.Values["state"]

	//todo: 왜 삭제를 하나?
	//delete(session.Values, "state")
	//session.Save(r, w)

	if state != r.FormValue("state") {
		http.Error(w, "Invalid session state", http.StatusUnauthorized)
		return
	}

	token, err := auth.OAuthConf.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client := auth.OAuthConf.Client(oauth2.NoContext, token)
	userInfoResp, err := client.Get(auth.UserInfoAPIEndpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer userInfoResp.Body.Close()
	userInfo, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var authUser auth.User

	json.Unmarshal(userInfo, &authUser)

	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400,
	}
	session.Values["user"] = authUser.Email
	session.Values["username"] = authUser.Name
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

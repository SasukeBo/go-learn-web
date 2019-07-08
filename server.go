package main

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	_ "github.com/SasukeBo/learn-web/db"
	"github.com/astaxie/beego/session"
	"golang.org/x/net/websocket"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func report(ws *websocket.Conn) {
	for {
		var receive string

		if err := websocket.Message.Receive(ws, &receive); err != nil {
			fmt.Println("error:", err)
			break
		}

		fmt.Println(receive)
		websocket.Message.Send(ws, "server received")
	}
}

func getUserSlice(w http.ResponseWriter, _ *http.Request) {
	f1 := Friend{Fname: "liziqi"}
	f2 := Friend{Fname: "jackiezq"}
	t := template.New("fieldname example")
	t, _ = t.Parse(`
		hello {{.UserName}}!
		{{range .Emails}}
			an email {{.}}
		{{end}}
		{{with .Friends}}
		{{range .}}
			my friend name is {{.Fname}}
		{{end}}
		{{end}}
	`)
	p := Person{
		UserName: "sasuke",
		Emails:   []string{"809754210@qq.com", "wangbo@giabbs.cn"},
		Friends:  []*Friend{&f1, &f2},
	}

	t.Execute(w, p)
}

func getTemplate(w http.ResponseWriter, _ *http.Request) {
	u := User{UserName: "sasuke", UserAge: 25}
	t := template.New("fieldname example")
	t, _ = t.Parse("username {{.UserName}}, age {{.UserAge}}!")
	t.Execute(w, u)
}

func logRequest(r *http.Request) {
	fmt.Println("method:", r.Method, " URL:", r.URL)
}

func requestJSON(w http.ResponseWriter, _ *http.Request) {
	var u UserSlice
	u.Users = append(u.Users, User{UserName: "wangbo", UserAge: 25})
	u.Users = append(u.Users, User{UserName: "liziqi", UserAge: 24})
	b, _ := json.Marshal(u)
	w.Header().Set("Content-Type", "text/json")

	fmt.Fprintf(w, string(b))
}

func requestXML(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	v := &Servers{Version: "1"}
	v.Svs = append(v.Svs, server{"Shanghai_VPN", "127.0.0.1"})
	v.Svs = append(v.Svs, server{"Beijing_VPN", "127.0.0.2"})
	output, _ := xml.MarshalIndent(v, "", "    ")
	w.Header().Set("Content-Type", "application/xml")

	fmt.Fprintf(w, xml.Header+string(output))
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	cookie := http.Cookie{Name: "username", Value: "Sasuke Bo", Expires: expiration}
	http.SetCookie(w, &cookie)
	fmt.Fprintf(w, "Hello Sasuke!")
}

func gethelloName(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	for _, cookie := range r.Cookies() {
		fmt.Fprint(w, cookie.Name+": "+cookie.Value+"\n")
	}
}

func xss(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	fmt.Println("name", name)
	w.Header().Add("X-XSS-Protection", "0")
	fmt.Fprintf(w, name)
}

func login(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	r.ParseForm()
	defer sess.SessionRelease(w)
	username := sess.Get("username")
	fmt.Println("session username:", username)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("html/login.html")
		t.Execute(w, nil)
	} else {
		sess.Set("username", r.Form["username"])
	}
}

func count(w http.ResponseWriter, r *http.Request) {
	sess, _ := globalSessions.SessionStart(w, r)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
	} else {
		sess.Set("countnum", (ct.(int) + 1))
	}
	t, _ := template.ParseFiles("html/count.html")
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, sess.Get("countnum"))
}

/*
func login(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method == "GET" {
		currentTime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(currentTime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("html/login.html")
		t.Execute(w, token)
	} else {
		r.ParseForm()
		token := r.Form.Get("token")
		if token != "" {
			// TODO 验证token的合法性
			fmt.Println("验证token的合法性")
		} else {
			// TODO 不存在token报错
			fmt.Println("不存在token报错")
		}

		fmt.Println("username:", template.HTMLEscapeString(r.Form.Get("username")))
		fmt.Println("password:", template.HTMLEscapeString(r.Form.Get("password")))
		template.HTMLEscape(w, []byte(r.Form.Get("username")))
	}
}
*/

func upload(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("html/upload.html")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	port := "9090"

	// logic.TestDB()

	http.HandleFunc("/", sayhelloName)
	http.HandleFunc("/xss", xss)
	http.HandleFunc("/count", count)
	http.HandleFunc("/gethello", gethelloName)
	http.HandleFunc("/login", login)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/xml", requestXML)
	http.HandleFunc("/json", requestJSON)
	http.HandleFunc("/template", getTemplate)
	http.HandleFunc("/userslice", getUserSlice)
	http.Handle("/websocket", websocket.Handler(report))

	fmt.Println("OK! Server start listening on", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var globalSessions *session.Manager

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "GOSESSIONID",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "./tmp",
	}
	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}

// Servers doc false
type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"`
	Svs     []server `xml:"server"`
}

type server struct {
	ServerName string `xml:"serverName"`
	ServerIP   string `xml:"ServerIP"`
}

// User doc false
type User struct {
	UserName string `json:"username"`
	UserAge  int    `json:"age"`
}

// UserSlice doc false
type UserSlice struct {
	Users []User `json:"users"`
}

// Friend doc false
type Friend struct {
	Fname string
}

// Person doc false
type Person struct {
	UserName string
	Emails   []string
	Friends  []*Friend
}

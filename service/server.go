/*
* Name: Go-Http-Server
* Usage:
*   1. Regist a new user (Duplication check)
*   2. Login the user
* Use Frames: negroni, mux
*/

package service

import (
    "fmt"
    "encoding/json"
	"net/http"
    "io/ioutil"
    "crypto/sha256"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type User struct {
    Name string `json:"name"`
    Pass [32]byte `json:"pass"`
}

// User List to store all user into json file
type UserList struct {
    Users []User
}

// Use negroni to create new server
func NewServer() *negroni.Negroni {

    formatter := render.New(render.Options{
        IndentJSON: true,
    })

    n := negroni.Classic()
    mx := mux.NewRouter()

    initRoutes(mx, formatter)

    n.UseHandler(mx)
    return n
}

func initRoutes(mx *mux.Router, formatter *render.Render) {
    // Create regist and login handler
    mx.HandleFunc("/regist", registHandler(formatter)).Methods("POST")
    mx.HandleFunc("/login", loginHandler(formatter)).Methods("POST")
}

// Regist handler, check username and password validation and push new user
func registHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        req.ParseForm()
        name := req.Form.Get("name")
        if len(name) == 0 {
            formatter.JSON(w, http.StatusBadRequest, struct{ Error string }{"Bad Request from client!"})
            return
        }
        pass := req.Form.Get("pass")
        if len(pass) == 0 {
            formatter.JSON(w, http.StatusBadRequest, struct{ Error string }{"Bad Request from client!"})
        } else if checkDuplicate(name) {
            // Duplication check
            formatter.JSON(w, http.StatusNotAcceptable, struct{ Error string }{"This user name has already registed!"})
        } else {
            formatter.JSON(w, http.StatusOK, struct{ RegistInfo string }{fmt.Sprintf("Regist user %s successfully!", name)})
            newUser(name, pass)
        }
    }
}

// Login handler, do username and password check
func loginHandler(formatter *render.Render) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        req.ParseForm()
        name := req.Form.Get("name")
        if len(name) == 0 {
            formatter.JSON(w, http.StatusBadRequest, struct{ Error string }{"Bad Request from client!"})
            return
        }
        pass := req.Form.Get("pass")
        if len(pass) == 0 {
            formatter.JSON(w, http.StatusBadRequest, struct{ Error string }{"Bad Request from client!"})
            return
        }
        // Name and password check
        if checkLogin(name, pass) {
            formatter.JSON(w, http.StatusOK, struct{ LoginInfo string }{fmt.Sprintf("Login user %s successfully!", name)})
        } else {
            formatter.JSON(w, http.StatusNotAcceptable, struct{ Error string }{"Invalid user or wrong password!"})
        }

    }
}
// Create a new user
func newUser(name_, pass_ string) {
    // Get origin userInfo content
    data, err := ioutil.ReadFile("./userInfos/users.json")
    if err != nil {
        fmt.Println("New User fail...")
        return
    }
    var allUsers UserList
    json.Unmarshal(data, &allUsers)
    // Encode with json format
    allUsers.Users = append(allUsers.Users, User{Name:name_, Pass:sha256.Sum256([]byte(pass_))})
    res, err1 := json.Marshal(allUsers)
    if err1 != nil {
        fmt.Println("Marshal error")
        return
    }
    ioutil.WriteFile("./userInfos/users.json", res, 0666)
}

// Access fake db to check user duplication
func checkDuplicate(name string) bool {
    data, _ := ioutil.ReadFile("./userInfos/users.json")
    // Check duplication
    var allUsers UserList
    json.Unmarshal(data, &allUsers)
    for i:=0; i<len(allUsers.Users); i++ {
        if name == allUsers.Users[i].Name {
            return true
        }
    }
    return false
}

// Access fake db to check username and password 
func checkLogin(name, pass string) bool {
    data, _ := ioutil.ReadFile("./userInfos/users.json")
    // Create decoder
    var allUsers UserList
    json.Unmarshal(data, &allUsers)
    for i:=0; i<len(allUsers.Users); i++ {
        if name == allUsers.Users[i].Name && sha256.Sum256([]byte(pass)) == allUsers.Users[i].Pass {
            return true
        }
    }
    return false
}
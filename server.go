package main

import (
	"fmt"
	"sort"
	"os"
	"flag"
	"strings"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"github.com/gorilla/mux"
)

//------------------------------------------------------------------------------
//User struct

// User struct represents users stored in the system
type User struct {
	Email     string   `json:"email"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Password  string   `json:"password"`
}

// Return marshled User struct without password
func (user User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Email string
		Firstname string
		Lastname string
	}{
		Email: user.Email,
		Firstname: user.Firstname,
		Lastname: user.Lastname,
	})
}


// Primitive authetication function - login and password are compared with
// those stored in the given struct.
func (user *User) Authenticate(aLogin string, aPassword string) bool {
	return aLogin == user.Email && aPassword == user.Password
}

// Updates all the fields of user struct except e-mail.
// This is how the password can be changed.
func (user *User) Update(anUser *User) {
	if nil != anUser {
		user.Firstname = anUser.Firstname
		user.Lastname = anUser.Lastname
		user.Password = anUser.Password
	}
}

//------------------------------------------------------------------------------
//Users struct

// Users are list of the users known by app.
// In real life application this will not be in memory variable - the database is
// state of art solution.
type Users struct {
	UserList map[string]*User
	UserListMutex sync.Mutex
}

// Creates and initialize empty user 'db'
func InitUsers() *Users {
	newUsers := &Users{}
	newUsers.UserList = make(map[string]*User)
	return newUsers
}

// Returns the User struct with given e-mail.
func (users *Users) GetUserByEmail(anEmail string) *User {
	users.UserListMutex.Lock()
	defer users.UserListMutex.Unlock()

	user, _:= users.UserList[anEmail]
	// Comment
	return user
}

// Call update on User object that has the same login/email
// as given anUser object.
// Returns true if object was found and update was called false otherwise.
func (users *Users) UpdateUserInfo(anUser *User) bool {
	if nil == anUser {
		return false
	}

	if nil == users.GetUserByEmail(anUser.Email) {
		//user doesn't exists - nothing to update
		return false
	}

	users.UserListMutex.Lock()
	defer users.UserListMutex.Unlock()
	users.UserList[anUser.Email].Update(anUser)
	return true
	// Comment
}
// Comment

// Insert anUser to user 'db'
// Returns true if user was successfully stored false otherwise
func (users *Users) InsertUser(anUser *User) bool {
	if nil == anUser {
		return false
	}

	if nil != users.GetUserByEmail(anUser.Email) {
		// user already exists - cannot be inserted
		return false
	}

	users.UserListMutex.Lock()
	defer users.UserListMutex.Unlock()
	users.UserList[anUser.Email] = anUser
	return true
}

//------------------------------------------------------------------------------
//Team struct

// Team representation
type Team struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`

	Members      map[string]*User `json:"-"`
	MembersMutex sync.Mutex       `json:"-"`
}

// Return marshaled team with added MemberCount field
func (team Team) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string
		Description string
		MemberCount int
	}{
		Name: team.Name,
		Description: team.Description,
		MemberCount: team.GetMembersCount(),
	})
}

// Creates instance of Team and initialize it
func InitTeam() *Team {
	newTeam := &Team{}
	newTeam.Members = make(map[string]*User)
	return newTeam
}

// Return count of members in given team
func (team *Team) GetMembersCount() int {
	team.MembersMutex.Lock()
	defer team.MembersMutex.Unlock()

	return len(team.Members)
}

// Add anUser to the given team members.
// Return true if successfully added
func (team *Team) Subscribe(anUser *User) bool {
	if nil == anUser {
		return true
	}

	team.MembersMutex.Lock()
	defer team.MembersMutex.Unlock()

	if _, exist := team.Members[anUser.Email]; exist {
		//subscription already done
		return false
	}
	//new subscription done
	team.Members[anUser.Email] = anUser
	return true
}

// Remove anUser from given team members.
// Return true if successfully removed.
func (team *Team) Unsubscribe(anUser *User) bool {
	if nil == anUser {
		return true
	}

	team.MembersMutex.Lock()
	defer team.MembersMutex.Unlock()

	if _, exist := team.Members[anUser.Email]; exist {
		//subscription will be deleted
		delete(team.Members, anUser.Email)
		return true
	}

	//no subscription found
	return false
}

//------------------------------------------------------------------------------
//Teams struct

// Teams 'db'
type Teams struct {
	TeamList map[string]*Team
	TeamListMutex sync.Mutex
}

// Create instance of Teams and initialize it
func InitTeams() *Teams {
	newTeams := &Teams{}
	newTeams.TeamList = make(map[string]*Team)
	return newTeams
}

// Insert aTeam to Teams 'db'
// Returns true if successfully added false otherwise
func (teams *Teams) Insert(aTeam *Team) bool {
	teams.TeamListMutex.Lock()
	defer teams.TeamListMutex.Unlock()

	if _, exists := teams.TeamList[aTeam.Name]; exists {
		return false
	} else {
		teams.TeamList[aTeam.Name] = aTeam
		return true
	}
}

// TeamList is map that
// Returns slice of the teams.
func (teams *Teams) GetTeamSlice() []*Team {
	teams.TeamListMutex.Lock()
	defer teams.TeamListMutex.Unlock()

	var result []*Team

	for _, team := range(teams.TeamList) {
		result = append(result, team)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// Returns the Team struct with name that equals aTeamName or nil if not found
func (teams *Teams) GetTeamByName(aTeamName string) *Team {
	teams.TeamListMutex.Lock()
	defer teams.TeamListMutex.Unlock()

	team, _ := teams.TeamList[aTeamName]
	return team
}


//------------------------------------------------------------------------------
//http auth

// HandlerFunc chain
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}

	return h
}

// Retrieve credentials from http request
// Return true, login, passwd if the creadentials found false, "", "" otherwise
func getCredentials(r *http.Request) (bool, string, string) {
	s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 {
		return false, "", ""
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return false, "", ""
	}

	pair := strings.SplitN(string(b), ":", 2)
	if len(pair) != 2 {
		return false, "", ""
	}

	//returns login and passwd
	return true, pair[0], pair[1]
}

// Authenticaion agains 'user db'
func basicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		if ok, login, passwd := getCredentials(r); ok {
			var user *User

			if user = users.GetUserByEmail(login); user == nil {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}

			if !user.Authenticate(login, passwd) {
				http.Error(w, "Not authorized", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}


		h.ServeHTTP(w, r)
	}
}
// Comment

//------------------------------------------------------------------------------
//http handlers
func GetUser(w http.ResponseWriter, r *http.Request) {
	_, login, _ := getCredentials(r)
	// Comment

	if user := users.GetUserByEmail(login); user != nil {
		json.NewEncoder(w).Encode(user)
		log.Printf("User [%s] get his info.", login)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u User
	decoder.Decode(&u)

	if u.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		// Comment
		log.Print("Email missing - user not created.")
		return
	}

	if u.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Password missing - user not created.")
		return
	}

	if users.InsertUser(&u) {
		w.WriteHeader(http.StatusCreated)
		log.Printf("User [%s] created.", u.Email)
	} else {
		w.WriteHeader(http.StatusNotModified)
		log.Printf("User [%s] not created.", u.Email)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u User
	decoder.Decode(&u)
	if users.UpdateUserInfo(&u) {
		w.WriteHeader(http.StatusAccepted)
		log.Print("User [%s] modified.", u.Email)
	} else {
		w.WriteHeader(http.StatusNotModified)
		log.Print("User [%s] not modified.", u.Email)
	}
}

func GetTeams(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(teams.GetTeamSlice())
	log.Print("TeamList returned.")
}

func CreateTeam(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	team := InitTeam()
	decoder.Decode(team)

	if team.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Team name is missing - team not created.")
		return
	}

	if teams.Insert(team) {
		w.WriteHeader(http.StatusCreated)
		log.Printf("Team [%s] successfully created.", team.Name)
	} else {
		w.WriteHeader(http.StatusNotModified)
		log.Printf("Team [%s] not created (already exists).", team.Name)
		// Comment
	}
}

func GetTeam(w http.ResponseWriter, r *http.Request) {
	team_name := mux.Vars(r)["team_name"]
	if team := teams.GetTeamByName(team_name); team != nil {
		json.NewEncoder(w).Encode(team)
		log.Print("Team [%s] returned.", team_name)
	} else {
		w.WriteHeader(http.StatusNoContent)
		log.Print("Team [%s] not found.", team_name)
	}
}

func SubscribeToTeam(w http.ResponseWriter, r *http.Request) {
	_, login, _ := getCredentials(r)

	if user := users.GetUserByEmail(login); user == nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		team_name := mux.Vars(r)["team_name"]
		if team := teams.GetTeamByName(team_name); team != nil {
			if team.Subscribe(user) {
				w.WriteHeader(http.StatusAccepted)
				log.Printf("User [%s] added to team [%s].", login, team_name)
			} else {
				w.WriteHeader(http.StatusNotModified)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func UnsubscribeToTeam(w http.ResponseWriter, r *http.Request) {
	_, login, _ := getCredentials(r)

	if user := users.GetUserByEmail(login); user == nil {
		w.WriteHeader(http.StatusForbidden)
	} else {
		team_name := mux.Vars(r)["team_name"]
		if team := teams.GetTeamByName(team_name); team != nil {
			if team.Unsubscribe(user) {
				w.WriteHeader(http.StatusAccepted)
				log.Printf("User [%s] removed from team [%s].", login, team_name)
			} else {
				w.WriteHeader(http.StatusNotModified)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

var users *Users
var teams *Teams

func main() {
	bindPtr := flag.String("bind", ":8000", "server host and port")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Http server with implemented two RestApi endpoints. The server uses `basic http auth` - username ~ user.email, password ~ user.password.\n\n")

		fmt.Fprintf(flag.CommandLine.Output(), "/api/user - PUT ~ creates user\n")
		fmt.Fprintf(flag.CommandLine.Output(), "/api/user - GET ~ retrieve user info (http auth)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "/api/user - POST ~ updates user info (http auth)\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "/api/team - GET ~ retrieve list of all teams (http auth)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "/api/team - PUT ~ creates new team (http auth)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "/api/team/{team_name} - GET ~ retrieve team info (http auth)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "/api/team/{team_name}/subscribe - POST ~ user will be added to members of team (http auth)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "/api/team/{team_name}/unsubscribe - POST ~ user will be removed from members of team (http auth)\n\n")

		fmt.Fprintf(flag.CommandLine.Output(), "Data description:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "User ~ {\"email\": \"\", \"firstname\": \"\", \"secondname\": \"\", \"password\": \"\"}.\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Team ~ {\"name\": \"\", \"description\": \"\"}.\n\n")

	}


	flag.Parse()

	users = InitUsers()
	teams = InitTeams()

	router := mux.NewRouter()
	router.HandleFunc("/api/user", CreateUser).Methods("PUT")
	router.HandleFunc("/api/user", use(GetUser, basicAuth)).Methods("GET")
	router.HandleFunc("/api/user", use(UpdateUser, basicAuth)).Methods("POST")
	// Comment

	router.HandleFunc("/api/team", use(GetTeams, basicAuth)).Methods("GET")
	router.HandleFunc("/api/team", use(CreateTeam, basicAuth)).Methods("PUT")
	router.HandleFunc("/api/team/{team_name:.*}", use(GetTeam, basicAuth)).Methods("GET")
	router.HandleFunc("/api/team/{team_name:.*}/subscribe", use(SubscribeToTeam, basicAuth)).Methods("POST")
	router.HandleFunc("/api/team/{team_name:.*}/unsubscribe", use(UnsubscribeToTeam, basicAuth)).Methods("POST")

	log.Printf("Server running [%s]...", *bindPtr)
	log.Fatal(http.ListenAndServe(*bindPtr, router))
}

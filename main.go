package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// the language used is Golang
const (
	dbUsername = "root"
	dbPassword = ""
	dbHostname = "localhost"
	dbPort     = "3306"
	dbName     = "go"
)

// BELOW I CALLED THE FUNCTIONS THAT I WAS TO USE AND DEFINED A PORT SO AS TO RUN THE PAGE WHICH WILL BE PORT 8080

func main()  {

	http.HandleFunc("/", Initiate)
	http.HandleFunc("/LoginInit", LoginInit)
	http.HandleFunc("/signup", SignupHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/home", rendeHome)
	http.HandleFunc("/logout", LogOut)
	// routing css files
	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))

	fmt.Println("Server ready at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
	
}
//I INITIATE THE SIGN UP PAGE HERE
func Initiate(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)


}

//I HANDLE SIGN UP PROCESSES HERE. CONNECTIONS ARE MADE TO THE DATABASE AND DATA IS INSERTED


func SignupHandler(w http.ResponseWriter, r *http.Request){
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHostname,dbPort,dbName)
	db, err := sql.Open("mysql",dns)

	if err != nil {
		log.Fatalf("impossible to create the connection: %s", err)
	}
	defer db.Close()

	if r.Method != http.MethodPost{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == http.MethodPost{

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		_,ins := db.Exec("INSERT INTO `users` (name, email, password) VALUES(?, ?, ?)", name, email, password)
		if ins != nil{
			log.Printf("Impossible to insert: %s", ins)
			return
		}
		//ONCE DATA IS INSERTED SUCCESSFULLY, THE USER IS REDIRECTED TO THE HOME PAGE
		http.Redirect(w,r,"/home", http.StatusSeeOther)
	} 
}

// HERE LOG IN PAGE IS INITIATED IF THE USER IS LOGGING IN

func LoginInit(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("template/login.html")
	if err != nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)

}

//LOG IN PROCESSES ARE HANDLED HERE IN A SIMILAR WAY TO THE SIGN UP

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHostname, dbPort, dbName)
    db, err := sql.Open("mysql", dns)

    if err != nil {
        log.Printf("Unable to connect to db: %s", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    if r.Method != http.MethodPost {
        http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
        return
    }

    email := r.FormValue("email")
    password := r.FormValue("password")

    // Use QueryRow instead of Exec and scan the result into variables
	// WE SELECT FROM THE DATABASE THIS TIME, NOT INSERT. THIS IS THE MAIN DIFFERENCE FROM THE SIGN UP PAGE
    var resultEmail, resultPassword string
    err = db.QueryRow("SELECT email, password FROM `users` WHERE email = ? AND password = ?", email, password).Scan(&resultEmail, &resultPassword)

    if err != nil {
		errorMessage := "Invalid credentials. Please check your email and password."
        renderLoginPageWithMessage(w, errorMessage)
        return
    }else{
		successMessage:="Logged in successfully!"
		//SUCCESS MESSAGE IS DISPLAYED AND USER IS REDIRECTED TO THE HOME PAGE
		renderHomePageWithMessage(w, successMessage)
		return
	}

    // Check if the retrieved email and password match the input
    if resultEmail != email || resultPassword != password {
        log.Printf("Invalid credentials for email: %s", email)
		//IF EMAIL AND PASSWORD DO NOT EXIST, THE USER IS REDIRECTED TO THE SIGN UP PAGE
        http.Redirect(w, r, "/signup", http.StatusSeeOther)
        return
    }

   
}
// THIS FUNCTION DISPLAYS THE ERROR MESSAGE
func renderLoginPageWithMessage(w http.ResponseWriter, errorMessage string) {
    tmpl, err := template.ParseFiles("template/login.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    data := struct {
        ErrorMessage string
    }{
        ErrorMessage: errorMessage,
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
// THIS FUNCTION DISPLAYS THE SUCCESS MESSAGE
func renderHomePageWithMessage(w http.ResponseWriter, successMessage string) {
    tmpl, err := template.ParseFiles("template/home.html")
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    data := struct {
		SuccessMessage string
    }{
		SuccessMessage: successMessage,
    }

    err = tmpl.Execute(w, data)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)

        return
    }
	
}
//THIS FUNCTION INITIATES HOME PAGE

func rendeHome(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("template/home.html")
	if err != nil{
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil{
		log.Fatal(err)
	}

}
//THIS FUNCTION DOES LOG OUT PROCESSES
func LogOut(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("template/login.html")
	if err != nil{
		log.Fatal(err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil{
		log.Fatal(err)
	}
}

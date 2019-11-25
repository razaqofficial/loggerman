package HomeController

import (
	"fmt"
	"github.com/leekchan/accounting"
	"github.com/thedevsaddam/govalidator"
	"html/template"
	"io"
	"log"
	"loggerman/database/model/userModel"
	"math/big"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"sync"
	"time"
)

var wG sync.WaitGroup

type PageTitle struct {
	Title string
}

func Index(w http.ResponseWriter, r * http.Request) {
	tmpl,_ :=  template.ParseFiles("templates/layout/default.html","templates/index.html")
	pageData := PageTitle{Title:"Welcome to Loggerman"}
	tmpl.Execute(w,pageData)
}

func About(w http.ResponseWriter, r * http.Request) {
	tmpl,_ :=  template.ParseFiles("templates/layout/default.html","templates/about.html")
	tmpl.Execute(w,nil)
}
type Logger struct {
	Name string
	Age string
	Description string
	Email string
	Strength int
	Image string
	Price string
}

func formatPrice(l []Logger)  []Logger  {
	for index,_ := range l {
		ac := accounting.Accounting{Symbol:"", Precision: 2}
		parseToFloat,_ := strconv.ParseFloat(l[index].Price,10)
		price := ac.FormatMoneyBigFloat(big.NewFloat(parseToFloat))
		l[index].Price = price
	}
	return l
}

func Loggers(w http.ResponseWriter, r * http.Request) {
	db := userModel.Connection()
	var logs []Logger
	type loggerAgg struct {
		Title string
		Loggers []Logger
	}
	db.Table("users").Select("name, age, description, email, strength, image, price").Scan(&logs)

	p := loggerAgg{Title: "Hire A Logger",Loggers:formatPrice(logs)}
	tmpl,_ :=  template.ParseFiles("templates/layout/default.html","templates/loggers.html")
	tmpl.Execute(w,p)
}

func AddLogger(w http.ResponseWriter, r * http.Request)  {
	if r.Method == "GET" {
		tmpl,_ :=  template.ParseFiles("templates/layout/default.html","templates/add-logger.html")
		tmpl.Execute(w,nil)
		return
	}

	type UserData struct {
		Name string
		Age int
		Strength int
		Email string
	}
	type PageData struct {
		Error map[string][]string
	}
	rules := govalidator.MapData{
		"name" : []string{"required","max:40"},
		"age" : []string{"required","max:2","numeric"},
		"strength" : []string{"required","numeric"},
		"email" : []string{"required","email"},
	}
	message := govalidator.MapData{
		"name":[]string{"required: Your Fullname is required","max: Maximum chars is 40"},
		"age":[]string{"required: Your Age is required","max: Maximum chars is 2","numeric: Age must be numeric"},
		"strength":[]string{"required: Your Strength is required","numeric: Strength must be numeric"},
		"email":[]string{"required: Your Email Address is required","email: Your email is not a valid email address"},
	}
	opts := govalidator.Options{
		Data:            nil,
		Request:         r,
		RequiredDefault: false,
		Rules:           rules,
		Messages:        message,
		TagIdentifier:   "",
		FormSize:        0,
	}
	v := govalidator.New(opts)
	e := v.Validate()

	if len(e) > 0 {
		fmt.Println("Error naa")
		http.Redirect(w, r, "/add/logger",301)
	}
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fileName := time.Now().Format("20160102150405") + ".jpg"
	fmt.Println("Original Filename is:", handler.Filename)
	fmt.Println("File size is:",handler.Size)
	fmt.Println("File type is:", handler.Header)
	f, err := os.OpenFile("./static/images/" + fileName, os.O_WRONLY|os.O_CREATE,0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	strength,_ := strconv.Atoi(r.FormValue("strength"))
	logger := Logger{
		r.FormValue("name"),
		r.FormValue("age"),
		r.FormValue("description"),
		r.FormValue("email"),
		strength,
		fileName,
		r.FormValue("price"),
	}
	db := userModel.Connection()
	if err := db.Table("users").Create(logger); err.Error != nil {
		fmt.Println("error occurred", err)
		http.Redirect(w, r,"/add/logger",301)
	}
	http.Redirect(w, r,"/loggers",301)
}

func Forum(w http.ResponseWriter, r * http.Request) {
	tmpl,_ :=  template.ParseFiles("templates/layout/default.html","templates/forum.html")
	tmpl.Execute(w,nil)
}

func Contact(w http.ResponseWriter, r * http.Request) {
	if r.Method == "GET" {
		tmpl,_ :=  template.ParseFiles("templates/layout/default.html","templates/contact.html")
		tmpl.Execute(w,nil)
		return
	}
	var (
		from       = "razaqofficial@gmail.com"
		msg = []byte("To: razaqofficial@gmail.com\r\n" +
			"From: " + r.FormValue("email") + "\r\n" +
			"Subject: Contact us message!\r\n" +
			"\r\n" +
			"You have a new message from " + r.FormValue("name")  + "\r\n" +
			r.FormValue("message") +".\r\n")
		recipients = []string{"razaqofficial@gmail.com"}
	)

	wG.Add(1)
	go sendMail(from, msg,recipients)
	wG.Wait()

	http.Redirect(w, r, "/",http.StatusFound)
}


func sendMail(from string, msg []byte, recipients []string) {
	defer wG.Done()
	// Set up authentication information.
	auth := smtp.PlainAuth("", "05550cf3b93299", "31f7efc89925b9", "smtp.mailtrap.io")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail("smtp.mailtrap.io:2525", auth, "",recipients, msg)
	if err != nil {
		log.Fatal(err)
	}
}
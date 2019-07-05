package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"

	"github.com/froprintoai/modernWeb/data"
	"github.com/froprintoai/modernWeb/loglog"
	"github.com/julienschmidt/httprouter"
)

//struct for storing temporary data when signing up
type UserSignUp struct {
	first, last      string
	email            string
	password         string
	month, day, year string
	uuid             string
}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if c, err := r.Cookie("_cookie"); err != nil {
		loginSignup(w, r)
	} else {
		fmt.Fprintln(w, "You have right Cookie", c)
	}
}

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user_data := UserSignUp{}
	if err := r.ParseForm(); err != nil {
		//Need to be modified to show an custim error page
		loglog.LogWTF("Parsing Error:", err)
	}
	//Put form values into struct UserSignUp
	user_data.first = r.PostFormValue("FirstName")
	user_data.last = r.PostFormValue("LastName")
	user_data.email = r.PostFormValue("Email")
	user_data.password = r.PostFormValue("Password")
	user_data.month = r.PostFormValue("Month")
	user_data.day = r.PostFormValue("Day")
	user_data.year = r.PostFormValue("Year")

	if msg, f := checkValidInput(&user_data); f == true {
		//generate UUID
		var err error
		user_data.uuid, err = createUUID()
		if err != nil {
			//Need to be modified to show an custim error page
			loglog.LogWTF("Cannot Create UUID : ", err)
		}
		//encrypt password
		encrypted_pass, salt := encryptPassword(user_data.password)
		//send an email(auth file activation)
		//create Activation Code
		activation_code, err := createACode()
		if err != nil {
			//Need to be modified to show an custim error page
			loglog.LogWTF("Cannot create activation code", err)
		}
		//create temporary File
		err = createTempFile(user_data, activation_code, encrypted_pass, salt)
		if err != nil {
			//Need to be modified to show an custim error page
			loglog.LogWTF("cannot create file", err)
		}
		//create mail
		url := createURL(user_data.email, activation_code)
		err = sendMail(url, user_data.email)
		if err != nil {
			//Need to be modified to show an custim error page
			loglog.LogWTF("there is an error when sending an mail", err)
		}
		//show the mail-sent page
		fmt.Fprintf(w, "We sent an email to %s. In order to finish registration, please access to an URL attached to the email in 30 minutes.", user_data.email)
	} else {
		fmt.Fprintln(w, user_data, msg)
	}

}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		loglog.LogWTF("Cannot parse the login form", err)
	}
	email := r.PostFormValue("Email")
	user, err := data.User_by_email(email)
	if err != nil {
		http.Redirect(w, r, "/", 301)
	}
	encrypted := fmt.Sprintf("%x", sha256.Sum256([]byte(r.PostFormValue("Password")+user.Salt)))
	if encrypted == user.Password {
		setCookie(user.Uuid, w)
		http.Redirect(w, r, "/", 301)
	}
}

func confirm(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//deal with accesses to URL like localhost:8080/confirm/(hashed email)/(activation code)
	hashed_email_challenged := ps.ByName("hashed_email")
	activation_code_challenged := ps.ByName("activation_code")

	filename := "temporary/" + hashed_email_challenged + ".txt"
	f, err := os.Open(filename)
	if err != nil {
		loglog.LogWarning("Challenge failed", err)
		fmt.Fprintln(w, "Your access is invalid:1")
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan() //scan the activation code on the  first line
	activation_code := scanner.Text()
	if activation_code == activation_code_challenged {
		//success
		user := data.User{}
		//Register user information to DB
		scanner.Scan()
		user.Uuid = scanner.Text()
		scanner.Scan()
		user.First = scanner.Text()
		scanner.Scan()
		user.Last = scanner.Text()
		scanner.Scan()
		user.Email = scanner.Text()
		scanner.Scan()
		user.Birthday = scanner.Text()
		scanner.Scan()
		user.Password = scanner.Text()
		scanner.Scan()
		user.Salt = scanner.Text()

		f.Close()
		os.Remove(filename)
		err := user.Register()
		if err != nil {
			loglog.LogWTF("Registration problem", err)
		}
		//ここでクッキーをセットするとまずい。なぜなら適用ドメインが/confirm/hashed_email/になり、
		//ここからのアクセスじゃないと、クッキーを送ってくれないから
		//よってCookie structのdomainとpathをいじることで、どこからのクッキーを設定する。
		setCookie(user.Uuid, w)
		//lead the user to the home page
		//successRegistration(w, r)
		http.Redirect(w, r, "/", 301)
	} else {
		f.Close()
		//failure
		loglog.LogWarning("Activation code doesn't match", err)
		fmt.Fprintln(w, "Your access is invalid:2")
	}
}

func loginSignup(w http.ResponseWriter, r *http.Request) {
	var b bytes.Buffer
	err := signup_template.Execute(&b, nil)
	if err != nil {
		fmt.Fprint(w, "An error occured.")
		return
	}
	b.WriteTo(w)
}

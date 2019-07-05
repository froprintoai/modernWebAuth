package main

import (
	"bufio"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/SlyMarbo/gmail"
	"github.com/froprintoai/modernWeb/data"
)

//check user data is valid when signing in
func checkValidInput(user_data *UserSignUp) (msg string, f bool) {
	msg = ""
	f = false
	//check 1.every input has value except for and birthday
	if user_data.first == "" || user_data.last == "" || user_data.password == "" || user_data.email == "" {
		msg = "There is a problem in name, password, or email."
		return
	}
	//check 1a. birthday
	day, err1 := strconv.Atoi(user_data.day)
	year, err2 := strconv.Atoi(user_data.year)
	if err1 != nil || err2 != nil {
		msg = "day or year is not the number"
		return
	} else if day < 1 || day > 31 || year < 0 { //It's desirable to check if the date really exists
		msg = "day and year are numbers but it's not valid for birthday"
		return
	}
	//2.password has 8 characters or more
	if len(user_data.password) < 8 {
		msg = "password should be 8 or more"
		return
	}
	//3.there is not the same email address in DB
	_, err := data.User_by_email(user_data.email)
	if err == nil {
		msg = "There is already a user with the email address"
		return
	}
	f = true
	return
}

func createUUID() (uuid string, err error) {
	temp := make([]byte, 16) //128bits
	_, err = rand.Read(temp)
	if err != nil {
		uuid = ""
		return
	}
	temp[6] = (temp[6] & 0x0F) | 0x40 //version4
	temp[8] = (temp[8] & 0x3F) | 0x80 //variant1 (10xx)
	//or, temp[8] = (temp[8] | 0xC0) & 0xBF
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", temp[0:4], temp[4:6], temp[6:8], temp[8:10], temp[10:16])
	return
}

func encryptPassword(password string) (result string, salt string) {
	t := time.Now().String()
	t = t[:30] //take salt from the current time
	temp := md5.Sum([]byte(t))
	salt = fmt.Sprintf("%x", temp)
	hash := sha256.Sum256([]byte(password + salt))
	result = fmt.Sprintf("%x", hash)
	return
}

func createACode() (a_code string, err error) {
	temp := make([]byte, 16)
	_, err = rand.Read(temp)
	if err != nil {
		a_code = ""
		return
	}
	a_code = fmt.Sprintf("%x", temp)
	return
}

func createTempFile(user_data UserSignUp, code string, encrypted_pass string, salt string) (err error) {
	//create hashed email address
	hashed_email := fmt.Sprintf("%x", md5.Sum([]byte(user_data.email)))
	file, err := os.Create("temporary/" + hashed_email + ".txt")
	if err != nil {
		return
	}
	w := bufio.NewWriter(file)
	content := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s", code, user_data.uuid, user_data.first, user_data.last, user_data.email, formatBirthday(user_data.day, user_data.month, user_data.year), encrypted_pass, salt)
	w.WriteString(content)
	w.Flush()
	return
}

//format birthday into "1999-01-08"
func formatBirthday(day, month, year string) (formatted string) {
	day_i, _ := strconv.Atoi(day)
	month_i, _ := strconv.Atoi(month)

	if day_i < 10 {
		day = "0" + day
	}
	if month_i < 10 {
		month = "0" + month
	}
	formatted = year + "-" + month + "-" + day
	return
}

func createURL(email string, code string) (url string) {
	//create URL attached to email and used to confirm user when signing in
	url = conf.Path + "/confirm"
	hashed_email := fmt.Sprintf("%x", md5.Sum([]byte(email)))
	url = url + "/" + hashed_email + "/" + code
	return
}

func sendMail(url string, address string) (err error) {
	//authを成功させるためには、アカウントの設定で　安全性の低いアプリのアクセス　を有効にする必要がある。
	//this code is from gmail package
	content := fmt.Sprintf("Click URL below to confirm your registration.\n%s", url)
	email := gmail.Compose("Please Confirm your sign in", content)
	email.From = conf.Gmail
	email.Password = conf.Password

	// Defaults to "text/plain; charset=utf-8" if unset.
	email.ContentType = "text/html; charset=utf-8"

	email.AddRecipient(address)

	err = email.Send()
	return
}

func setCookie(uuid string, w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:     "_cookie",
		Value:    uuid,
		Domain:   conf.Path_without_port,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
}

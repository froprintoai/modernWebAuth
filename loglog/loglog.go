package loglog

import (
	"log"
	"os"
)

func init() {
	f, _ := os.Create("loglog/log.txt")
	f.Close()
}

func LogWarning(msg string, err error) {
	file, err := os.OpenFile("loglog/log.txt", os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln("Cannot Open Logfile:", err)
	}
	defer file.Close()
	logger := log.New(file, "Warning : ", log.Llongfile)
	logger.Println(msg, "error : ", err)
}

func LogWTF(msg string, err error) {
	file, err := os.OpenFile("loglog/log.txt", os.O_RDWR, 0666)
	if err != nil {
		log.Fatalln("Cannot Open Logfile:", err)
	}
	defer file.Close()
	logger := log.New(file, "WTF : ", log.Llongfile)
	logger.Fatalln(msg, "error : ", err)
}

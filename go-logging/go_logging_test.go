package go_logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_기본_Logger(t *testing.T) {
	log.Println("Logging") //2020/12/30 10:27:11 Logging
}

func Test_기본_Logger_날짜_시간_표시_X(t *testing.T) {
	log.SetFlags(0)
	log.Println("Logging") //Logging
}

func Test_기본_Logger_(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetPrefix("INFO: ")
	log.Println("Info")
}

type Logger struct {
	Trace *log.Logger
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger
}

var myLogger Logger

func logInit(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	myLogger.Trace = log.New(traceHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	myLogger.Info = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	myLogger.Warn = log.New(warningHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	myLogger.Error = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Test(t *testing.T) {
	logInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	myLogger.Info.Println("Starting the application...")
	myLogger.Info.Println("Something noteworthy happened")
	myLogger.Warn.Println("There is something you should know about")
	myLogger.Error.Println("Something went wrong")

	//INFO: 2020/12/30 15:04:22 go_logging_test.go:44: Starting the application...
	//INFO: 2020/12/30 15:04:22 go_logging_test.go:45: Something noteworthy happened
	//WARNING: 2020/12/30 15:04:22 go_logging_test.go:46: There is something you should know about
	//ERROR: 2020/12/30 15:04:22 go_logging_test.go:47: Something went wrong
}

var logger *log.Logger

func Test_Custom_Logger(t *testing.T) {
	logger = log.New(os.Stdout, "INFO: ", log.LstdFlags)
	logger.Println("Logging") //INFO: 2020/12/30 10:34:00 Logging
}

func Test_Custom_Logger_File(t *testing.T) {
	// 로그파일 오픈
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	logger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	printMsgLogger("test msg")
	logger.Println("End of Program")
}

func Test_기본_Logger_File(t *testing.T) {
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile) //로그 출력 위치를 파일로 변경

	printMsgLog("test msg")
	log.Println("End of Program")
}

func Test_Multiple_Outputs(t *testing.T) {
	logFile, err := os.OpenFile("logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)

	printMsgLog("test msg")
	log.Println("End of Program")
}

//todo : logger를 넘기려면 어떻게 인자를 작성해야 하나?
func printMsgLog(msg string) {
	log.Print(msg)
}

func printMsgLogger(msg string) {
	logger.Print(msg)
}

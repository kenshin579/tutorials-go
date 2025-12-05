package go_logging

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_Basic_Logger(t *testing.T) {
	log.Println("Logging") //2020/12/30 10:27:11 Logging
}

func Test_Fatal(t *testing.T) {
	log.Fatal("fatal") //메시지 출력 + os.Exit(1)
}

func Test_Panic(t *testing.T) {
	log.Panic("panic") //메시지 출력 + panic()
}

func Test_Basic_Logger_Flags_Setting_DateTime_Display_X(t *testing.T) {
	log.SetFlags(0)
	log.Println("Logging") //Logging
}

func Test_Basic_Logger_Flags_Setting2(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetPrefix("INFO: ")
	log.Println("Logging")

	//INFO: 2020/12/30 15:41:20 /Users/ykoh/GolandProjects/tutorials-go/go-logging/go_logging_test.go:23: Logging
}

func Test_Basic_Logger_File(t *testing.T) {
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

func printMsgLog(msg string) {
	log.Print(msg)
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

// todo : logger를 넘기려면 어떻게 인자를 작성해야 하나?
func printMsgLogger(msg string) {
	logger.Print(msg)
}

type Logger struct {
	Trace *log.Logger
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger
}

var myLogger Logger

func logInit(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	myLogger.Trace = log.New(traceHandle, "[TRACE] ", log.Ldate|log.Ltime|log.Lshortfile)
	myLogger.Info = log.New(infoHandle, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	myLogger.Warn = log.New(warningHandle, "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile)
	myLogger.Error = log.New(errorHandle, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Test(t *testing.T) {
	logInit(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	myLogger.Info.Println("Starting the application...")
	myLogger.Trace.Println("Something noteworthy happened")
	myLogger.Warn.Println("There is something you should know about")
	myLogger.Error.Println("Something went wrong")

	//[INFO] 2020/12/30 15:43:40 go_logging_test.go:46: Starting the application...
	//[WARNING] 2020/12/30 15:43:40 go_logging_test.go:48: There is something you should know about
	//[ERROR] 2020/12/30 15:43:40 go_logging_test.go:49: Something went wrong
}

package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mattbaird/gochimp"
)

type requestData struct {
	Name       string
	Location   string
	Speciality string
	Phone      string
	endSession bool
}

var (
	sendEmailChan = make(chan requestData, 10)
	mandrillAPI   *gochimp.MandrillAPI
	serverPort    = os.Getenv("FEEDBACK_PORT")
	targetEmail   = os.Getenv("FEEDBACK_TARGET_EMAIL")
)

func main() {
	var err error
	if mandrillAPI, err = gochimp.NewMandrill(os.Getenv("MANDRILL_KEY")); err != nil {
		fmt.Println("Unable to create mandrill client")
		os.Exit(1)
	}

	go sendEmail(sendEmailChan)

	if len(serverPort) == 0 {
		serverPort = "9000"
	}

	if len(targetEmail) == 0 {
		targetEmail = "sales@captureproof.com"
	}

	fmt.Println("Listening on port", serverPort)
	http.HandleFunc("/know-a-doctor", knowADoctor)
	http.ListenAndServe(":"+serverPort, nil)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	<-sigs

	fmt.Println("Exiting")
	sendEmailChan <- requestData{endSession: true}
}

func knowADoctor(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("%s [%s] %s\n", time.Now().Format(time.RFC822), req.Method, req.Header.Get("Origin"))

	origin := strings.Trim(req.Header.Get("Origin"), " ")
	sendCORSHeaders(w, origin)

	if req.Method == "OPTIONS" {
		return
	}

	if req.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !isValidHost(origin) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Invalid origin")
		return
	}

	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Error processing params: %s", err.Error())
		return
	}

	reqData := requestData{
		Name:       req.PostFormValue("name"),
		Location:   req.PostFormValue("location"),
		Speciality: req.PostFormValue("speciality"),
		Phone:      req.PostFormValue("phone"),
		endSession: false,
	}

	if invalidForm(reqData) {
		fmt.Println("Params are empty")
		w.WriteHeader(http.StatusOK)
		return
	}

	sendEmailChan <- reqData
	w.WriteHeader(http.StatusOK)
}

func sendCORSHeaders(w http.ResponseWriter, origin string) {
	if isValidHost(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Origin")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
}

func isValidHost(host string) (result bool) {
	result = strings.Contains(host, "captureproof.com") ||
		strings.Contains(host, "cp2.div-art.com.ua") ||
		strings.Contains(host, "localhost")
	return
}

func invalidForm(data requestData) (result bool) {
	result = len(strings.Trim(data.Name, " ")) == 0 &&
		len(strings.Trim(data.Location, " ")) == 0 &&
		len(strings.Trim(data.Speciality, " ")) == 0
	return
}

func sendEmail(queue chan requestData) {
	for {
		data := <-queue
		if data.endSession {
			break
		}
		sendMandrillEmail(data)
	}
}

func sendMandrillEmail(data requestData) {
	text := fmt.Sprintf(
		"A patient knows a doctor %s (%s), who is working at %s. Doctor's phone number is %s",
		data.Name,
		data.Speciality,
		data.Location,
		data.Phone)

	message := gochimp.Message{
		Text:      text,
		Subject:   "[marketing] Patient knows a doctor",
		FromEmail: "feedback@captureproof.com",
		FromName:  "Marketing Site",
		To: []gochimp.Recipient{
			gochimp.Recipient{
				Email: targetEmail,
				Name:  "Sales guys",
				Type:  "",
			},
		},
	}

	fmt.Println("Sending email: ", text)

	if _, err := mandrillAPI.MessageSend(message, false); err != nil {
		fmt.Printf("Error sending message: %s\n", err.Error())
	}
}

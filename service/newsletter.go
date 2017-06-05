package service

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/badoux/checkmail"

	gomail "gopkg.in/gomail.v2"

	"golang.org/x/crypto/bcrypt"

	restful "github.com/emicklei/go-restful"
	"github.com/ingmardrewing/gomicNewsletter/config"
	"github.com/ingmardrewing/gomicNewsletter/db"
)

type Content struct {
	Email string
}

type Newsletter struct {
	Body    string
	Subject string
}

type Msg struct {
	Text string
}

const (
	letterBytes   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 73 / letterIdxBits
)

var src = rand.NewSource(time.Now().UnixNano())

func getRandomString(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func getSubcriptionToken(length int) string {
	t := getRandomString(length)
	for db.TokenExists(t) {
		t = getRandomString(length)
	}
	return t
}

func getDeletionToken(length int) string {
	t := getRandomString(length)
	for db.DeletionTokenExists(t) {
		t = getRandomString(length)
	}
	return t
}

func NewNewsletterService() *restful.WebService {
	path := "/0.1/gomic/newsletter"

	add := "/add"
	delete := "/delete/{token}"
	verify := "/verify/{token}"
	send := "/send"

	service := new(restful.WebService)
	service.Path(path)

	log.Printf("Adding PUT route: %s\n", path+add)
	service.Route(service.PUT(add).Consumes(restful.MIME_JSON).To(Add))

	log.Printf("Adding GETroute: %s\n", path+delete)
	service.Route(service.GET(delete).Consumes(restful.MIME_JSON).To(Delete))

	log.Printf("Adding GET route: %s\n", path+verify)
	service.Route(service.GET(verify).Consumes(restful.MIME_JSON).To(Verify))

	log.Printf("Trigger Newsletter via POST route: %s\n", path+send)
	service.Route(service.POST(send).Consumes(restful.MIME_JSON).Filter(basicAuthenticate).To(Send))

	return service
}

func Verify(request *restful.Request, response *restful.Response) {
	token := request.PathParameter("token")
	msg := new(Msg)
	if db.TokenExists(token) {
		db.VerifySubscription(token)
		msg.Text = "Your email address has successfully been registered."
	} else {
		msg.Text = "Unfortunately the e-mail address could not be registered. Token not found: " + token
	}

	t, err := template.ParseFiles("verified.html")
	if err != nil {
		log.Fatalf("Template gave: %s", err)
	}
	t.Execute(response.ResponseWriter, msg)
}

func Add(request *restful.Request, response *restful.Response) {
	err, c := readContent(request)
	if err != nil {
		response.WriteErrorString(400, "400: Bad Request ("+err.Error()+")")
		return
	}
	token := getSubcriptionToken(60)
	deletion_token := getDeletionToken(60)
	msg := new(Msg)
	err = checkmail.ValidateFormat(c.Email)
	log.Println("Given Mail:")
	log.Println(c.Email)
	log.Println(err)
	if err != nil {
		msg.Text = "no address given"
	} else if db.AddressExists(c.Email) {
		msg.Text = "address already registered"
	} else {
		db.AddEmailAddress(c.Email, token, deletion_token)
		msg.Text = "Added " + c.Email + " Token: " + token

		link := "https://drewing.eu:16443/0.1/gomic/newsletter/verify/"
		link += token

		email_text := `Hi there,

I received word that you want to subscribe to the DevAbo.de newsletter.
All you need to do to complete the subscription is to click on the following link:

%s

In case this e-mail has reached you in error and you aren't interested, just delete this e-mail and you will not be bothered again.

Sincerely

Ingmar Drewing
`
		user, pass, host, port := config.GetSmtpCredentials()
		m := gomail.NewMessage()
		m.SetHeader("From", user)
		m.SetHeader("To", c.Email)
		m.SetHeader("Subject", "Newsletter Verification")
		m.SetBody("text/plain", fmt.Sprintf(email_text, link))

		portInt, _ := strconv.Atoi(port)

		d := gomail.NewDialer(host, portInt, user, pass)
		d.DialAndSend(m)

	}

	response.WriteEntity(msg)
}

func Delete(request *restful.Request, response *restful.Response) {
	token := request.PathParameter("token")
	msg := new(Msg)
	if db.TokenExists(token) {
		db.DeleteEmailAddressWithToken(token)
		msg.Text = "Deleted mailadress with token " + token
	} else {
		msg.Text = "No record with exists with token " + token
	}
	response.WriteEntity(msg)
}

func Send(request *restful.Request, response *restful.Response) {
	err, n := readNewsletter(request)
	if err != nil {
		response.WriteErrorString(400, "400: Bad Request ("+err.Error()+")")
		return
	}

	user, pass, host, port := config.GetSmtpCredentials()
	portInt, _ := strconv.Atoi(port)
	dialer := gomail.NewDialer(host, portInt, user, pass)
	recipients := db.GetNewsletterRecipients()

	for _, recipient := range recipients {
		txt := n.Body + "\n\nIf you wish to stop receiving e-mails like this one, just click on the following link: https://drewing.eu:16443/0.1/gomic/newsletter/delete/" + recipient.DeletionToken
		sendMail(user, dialer, n.Subject, txt, recipient.Email)
	}

	msg := new(Msg)
	msg.Text = "Newsletter sent"
	response.WriteEntity(msg)
}

func sendMail(user string, dialer *gomail.Dialer, subject string, body string, recipient string) {

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	dialer.DialAndSend(m)
}

func basicAuthenticate(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	err := authenticate(request)
	log.Println(err)
	if err != nil {
		response.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		response.WriteErrorString(401, "401: Not Authorized")
		return
	}

	chain.ProcessFilter(request, response)
}

func authenticate(req *restful.Request) error {
	user, pass, _ := req.Request.BasicAuth()
	given_pass := []byte(pass)
	stored_hash := []byte(config.GetPasswordHashForUser(user))
	//hash, _ := bcrypt.GenerateFromPassword(given_pass, coast)
	return bcrypt.CompareHashAndPassword(stored_hash, given_pass)
}

func checkContent(c *Content) error {
	msg := []string{}
	if len(c.Email) == 0 {
		msg = append(msg, "No E-Mail-Adress given")
	}

	if len(msg) > 0 {
		return errors.New(strings.Join(msg, ", "))
	}
	return nil
}

func checkNewsletter(n *Newsletter) error {
	msg := []string{}
	if len(n.Subject) == 0 {
		msg = append(msg, "No subject for newsletter given")
	}
	if len(n.Body) == 0 {
		msg = append(msg, "No text for newsletter given")
	}

	if len(msg) > 0 {
		return errors.New(strings.Join(msg, ", "))
	}
	return nil
}

func readNewsletter(request *restful.Request) (error, *Newsletter) {
	n := new(Newsletter)
	request.ReadEntity(n)
	err := checkNewsletter(n)
	if err != nil {
		return err, new(Newsletter)
	}
	return nil, n
}

func readContent(request *restful.Request) (error, *Content) {
	c := new(Content)
	request.ReadEntity(c)
	err := checkContent(c)
	if err != nil {
		return err, new(Content)
	}
	return nil, c
}

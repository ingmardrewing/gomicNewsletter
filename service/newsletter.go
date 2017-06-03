package service

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

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
	Txt string
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

func NewNewsletterService() *restful.WebService {
	path := "/0.1/gomic/newsletter"

	add := "/add"
	delete := "/delete/{token}"
	verify := "/verify/{token}"
	send := "/send"

	service := new(restful.WebService)
	service.
		Path(path).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	log.Printf("Adding PUT route: %s\n", path+add)
	service.Route(service.PUT(add).To(Add))

	log.Printf("Adding DELETE route: %s\n", path+delete)
	service.Route(service.DELETE(delete).To(Delete))

	log.Printf("Adding POST route: %s\n", path+verify)
	service.Route(service.POST(verify).To(Verify))

	log.Printf("Trigger Newsletter via POST route: %s\n", path+send)
	service.Route(service.POST(send).Filter(basicAuthenticate).To(Send))

	return service
}

func Verify(request *restful.Request, response *restful.Response) {
	token := request.PathParameter("token")
	msg := new(Msg)
	if db.TokenExists(token) {
		db.VerifySubscription(token)
		msg.Txt = "Successfully verified with token: " + token
	} else {
		msg.Txt = "Token not found: " + token
	}
	response.WriteEntity(msg)
}

func Add(request *restful.Request, response *restful.Response) {
	err, c := readContent(request)
	if err != nil {
		response.WriteErrorString(400, "400: Bad Request ("+err.Error()+")")
		return
	}
	token := getRandomString(255)
	msg := new(Msg)
	if !db.AddressExists(c.Email) {
		db.AddEmailAddress(c.Email, token)
		msg.Txt = "Added " + c.Email + " Token: " + token
	} else {
		msg.Txt = c.Email + " already registered"
	}
	response.WriteEntity(msg)
}

func Delete(request *restful.Request, response *restful.Response) {
	token := request.PathParameter("token")
	msg := new(Msg)
	if db.TokenExists(token) {
		db.DeleteEmailAddressWithToken(token)
		msg.Txt = "Deleted mailadress with token " + token
	} else {
		msg.Txt = "No record with exists with token " + token
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
	recipients := db.GetNewsletterRecipients()

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("Bcc", recipients...)
	m.SetHeader("Subject", n.Subject)
	m.SetBody("text/plain", n.Body)

	portInt, _ := strconv.Atoi(port)

	d := gomail.NewDialer(host, portInt, user, pass)
	err = d.DialAndSend(m)

	msg := new(Msg)
	if err != nil {
		msg.Txt = err.Error()
	}
	response.WriteEntity(msg)
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

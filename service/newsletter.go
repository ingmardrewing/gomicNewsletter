package service

import (
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/ingmardrewing/gomicNewsletter/db"
)

type Content struct {
	Email string
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

/*
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
*/

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

func readContent(request *restful.Request) (error, *Content) {
	c := new(Content)
	request.ReadEntity(c)
	err := checkContent(c)
	if err != nil {
		return err, new(Content)
	}
	return nil, c
}

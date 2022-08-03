package chat_repository

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	gomail "gopkg.in/mail.v2"

	"github.com/StarsPoker/loginBackend/consts"
	"github.com/StarsPoker/loginBackend/domain/one_time_password"
	"github.com/StarsPoker/loginBackend/domain/users"
	"github.com/StarsPoker/loginBackend/logger"
)

const MESSAGE = " Use seu token para entrar no nosso site do GrupoSx e lembre-se de não compartilhar com outras pessoas. Seu token de acesso é: "

func getSlice() string {
	message := []string{
		"Oi, que bom te ver aqui de novo.",
		"Olá, bem vindo novamente.",
		"Estamos felizes em te ver por aqui.",
		"Olá, como vai?",
		"Oi, seu token foi gerado com sucesso.",
		"É um prazer te ver por aqui.",
	}
	rand_id := rand.Intn(len(message))
	return message[rand_id]
}

func SendWhatsappMessage(otp one_time_password.OneTimePassword, user *users.User) {

	url_post := consts.WHATSAPP_TI_URL
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	content := fmt.Sprint(getSlice(), MESSAGE, "*", otp.Code, "*")
	data, err := json.Marshal(map[string]interface{}{
		"message":   content,
		"phone":     user.Contact,
		"messageId": "",
	})

	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{Transport: tr}
	req, reqError := http.NewRequest("POST", url_post, bytes.NewBuffer(data))
	req.Header.Set("content-type", "application/json")
	if reqError != nil {
		fmt.Println(reqError)
	}
	resp, errSendMessage := client.Do(req)

	if errSendMessage != nil {
		logger.Error("error when trying to send message annotated", errSendMessage)
	} else if resp.StatusCode != 201 {
		defer resp.Body.Close()
		logger.Error("error when trying to send message annotated", errSendMessage)
	} else {
		defer resp.Body.Close()
	}
}

func SendMail(otp one_time_password.OneTimePassword, user *users.User) {
	content := fmt.Sprint(getSlice(), MESSAGE, otp.Code, "\n\nAtenciosamente, equipe do GrupoSx.")

	m := gomail.NewMessage()

	m.SetHeader("From", "ti@sxgrupo.com.br")
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Token de acesso ao GrupoSx")
	m.SetBody("text/plain", content)

	d := gomail.NewDialer("smtppro.zoho.com", 465, "ti@sxgrupo.com.br", "LafeZXgVTU")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}

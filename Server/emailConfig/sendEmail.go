package emailConfig

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/mskKandula/model"
	"github.com/mskKandula/variables"
	"gopkg.in/gomail.v2"
)

func SendEmail(user model.BasicDetails) error {

	body, err := makeTemplate(user, "./Templates/registrationMailTemplate.html")

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", variables.SenderEmail)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Registration Successfull!")
	m.SetBody("text/html", body)

	d := gomail.NewPlainDialer("smtp.gmail.com", 587, variables.SenderEmail, variables.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Email Sent")

	return nil

}

func makeTemplate(profileObj model.BasicDetails, templatePath string) (string, error) {

	parsedTemplate, err := template.ParseFiles(templatePath)

	if err != nil {
		return "", err
	}

	var buff bytes.Buffer

	parseErr := parsedTemplate.Execute(&buff, profileObj)

	if parseErr != nil {
		return "", parseErr
	}

	body := buff.String()

	return body, nil
}

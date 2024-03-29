package emailConfig

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"

	"github.com/mskKandula/oes/api/model"
	"github.com/mskKandula/oes/util/variables"
	gomail "gopkg.in/gomail.v2"
)

func SendEmail(user model.Student) error {

	body, err := makeTemplate(user, "registrationMailTemplate.html")

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", variables.SenderEmail)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Registration Successful!")
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, variables.SenderEmail, variables.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} //comment it out for Production

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	fmt.Println("Email Sent")

	return nil

}

func makeTemplate(profileObj model.Student, templatePath string) (string, error) {

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

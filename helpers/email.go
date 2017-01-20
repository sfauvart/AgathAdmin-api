package helpers

import (
	"bytes"
	"github.com/sfauvart/Agathadmin-api/settings"
	"gopkg.in/gomail.v2"
	"html/template"
)

func SendEmail(templateFileName string, templateData interface{}, to string, locale string) error {
	if locale == "" {
		locale = settings.Get().DefaultLocale
	}
	body, err := ParseTemplate("./emails/"+locale+"/"+templateFileName+".html", templateData)
	if err != nil {
		return err
	}
	subject, err := ParseTemplate("./emails/"+locale+"/"+templateFileName+"_subject.txt", templateData)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", settings.Get().NoReplyFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := settings.GetSmtpDial()
	return d.DialAndSend(m)
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)

	templateData := struct {
		Data         interface{}
		BaseFrontUrl string
	}{
		Data:         data,
		BaseFrontUrl: settings.Get().BaseFrontUrl,
	}

	if err = t.Execute(buf, templateData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

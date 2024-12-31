package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"mime"
	"mime/multipart"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"
)

const envelopeWithAttachmentsTemplate = `From: {{.From}}
Date: {{.Time.Format "Mon Jan 2 15:04:05 MST 2006"}}
Subject: {{.Subject}}
To: {{.To}}
MIME-version: 1.0;
Content-Type: multipart/mixed; boundary="{{.Boundary}}"
{{$boundary:=.Boundary}}

--{{$boundary}}
Content-Type: text/html; charset="UTF-8"

{{.Data}}

{{range $key, $value := .Attachments}}
--{{$boundary}}
Content-Type: {{$value.MimeType}}
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="{{$key}}"

{{b64 $value.Data}}
{{end}}
--{{$boundary}}--
`

const envelopeTemplate = `From: {{.From}}
Date: {{.Time.Format "Mon Jan 2 15:04:05 MST 2006"}}
Subject: {{.Subject}}
To: {{.To}}
MIME-version: 1.0;
Content-Type: text/html; charset="UTF-8"


{{.Data}}`

type SMTPSettings struct {
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPass       string
	SMTPFrom       string
	SMTPSenderName string
}

type Attachment struct {
	Filename string
	Data     []byte
	MimeType string
}

type MailEnvelope struct {
	Time        time.Time
	Data        template.HTML
	Attachments map[string]Attachment
	fromName    string
	fromAddress string
	subject     string
	Boundary    string
	recipients  map[string]string
	MimeType    string
}

func SendTemplateMail(
	smtpSettings SMTPSettings,
	subject string,
	recipients map[string]string,
	data string,
	attachments map[string]Attachment,
) error {
	var useTemplate string
	var mimeType string

	if attachments != nil {
		useTemplate = envelopeWithAttachmentsTemplate
		for key, value := range attachments {
			item := attachments[key]
			item.MimeType = mime.TypeByExtension(filepath.Ext(value.Filename))
			attachments[key] = item
		}
	} else {
		useTemplate = envelopeTemplate
	}

	templateExecutable := template.New("mail_envelope_template.html")
	templateExecutable.Funcs(template.FuncMap{
		"b64": func(b []byte) template.HTML {
			return template.HTML(base64.StdEncoding.EncodeToString(b))
		},
	})

	loginTemplate, err := templateExecutable.Parse(useTemplate)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buffer)
	boundary := writer.Boundary()

	envelope := MailEnvelope{
		Time:        time.Now(),
		Data:        template.HTML(data),
		Attachments: attachments,
		fromName:    smtpSettings.SMTPSenderName,
		fromAddress: smtpSettings.SMTPFrom,
		subject:     subject,
		recipients:  recipients,
		Boundary:    boundary,
		MimeType:    mimeType,
	}

	err = loginTemplate.Execute(buffer, envelope)
	if err != nil {
		return err
	}

	smtpHost := fmt.Sprintf("%s:%s", smtpSettings.SMTPHost, smtpSettings.SMTPPort)
	smtpUser := smtpSettings.SMTPUser
	smtpAuth := &loginAuth{
		username: smtpUser,
		password: smtpSettings.SMTPPass,
	}

	// Here we are handling the custom version of the mailer, using the SMTP protocol directly
	c, err := smtp.Dial(smtpHost)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.Auth(smtpAuth)
	if err != nil {
		return err
	}

	log.Printf("ENVELOPE SMTP\n%s", string(buffer.Bytes()))

	for name, address := range envelope.recipients {
		err = c.Mail(smtpUser)
		if err != nil {
			return err
		}

		err = c.Rcpt(fmt.Sprintf(`"%s" <%s>`, name, address))
		if err != nil {
			//user dont have a email
			if address == "" {
				return nil
			}
			return err
		}

		wc, err := c.Data()
		if err != nil {
			return err
		}

		_, err = wc.Write(buffer.Bytes())
		if err != nil {
			return err
		}

		err = wc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// From returns `From` field of the SMTP envelope
func (e MailEnvelope) From() template.HTML {
	return template.HTML(fmt.Sprintf(`"%s" <%s>;`, e.fromName, e.fromAddress))
}

// RecipientsSlice returns a slice of the recipient contacts formatted as required by the SMTP envelope
func (e MailEnvelope) RecipientsSlice() []string {
	collected := make([]string, 0)

	for name, address := range e.recipients {
		collected = append(collected, fmt.Sprintf(`"%s" <%s>`, name, address))
	}

	return collected
}

// To returns `To` field of the SMTP envelope
func (e MailEnvelope) To() template.HTML {
	return template.HTML(strings.Join(e.RecipientsSlice(), ";"))
}

// Subject returns `Subject` field of the SMTP envelope
func (e MailEnvelope) Subject() template.HTML {
	return template.HTML(e.subject)
}

type loginAuth struct {
	username string
	password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unknown from server")
		}
	}
	return nil, nil
}

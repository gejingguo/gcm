package gcm

import (
	"fmt"
	"strings"
	"net/smtp"
	"bytes"
	"io/ioutil"
	"encoding/base64"
	"time"
	"path/filepath"
	"mime"
)

const (
	SmtpMailContentTypePlain = 1
	SmtpMailContentTypeHtml = 2
)

//
type Mail struct {
	Host string 		// 服务器地址, ip:port
	User string 		// 发送者帐号
	Pass string 		// 发送者密码, 发送吗
	//From string 		// 发送者名称
	To []string			// 接受者列表
	CC []string			// 抄送者列表
	BCC []string 		// 暗抄送者列表
	ContentType string 	// 内容类型
	Subject string		// 标题
	Body string 		// 内容
	Files []string 		// 附件列表
}

func (m *Mail) SetAccount(host string, user string, pass string) {
	//m.From = from
	m.Host = host
	m.User = user
	m.Pass = pass
}

func NewTextMail(subject string, body string) *Mail {
	mail := &Mail{}
	mail.Subject = subject
	mail.Body = body
	mail.ContentType = "text/plain"
	return mail
}

func NewHtmlMail(subject string, body string) *Mail {
	mail := &Mail{}
	mail.Subject = subject
	mail.Body = body
	mail.ContentType = "text/html"
	return mail
}

func NewMail(contentType string, subject string, body string) *Mail {
	mail := &Mail{}
	mail.Subject = subject
	mail.Body = body
	mail.ContentType = contentType
	return mail
}

func (m *Mail) AddReceiver(receiver string)  {
	if m.To == nil {
		m.To = make([]string, 0, 10)
	}
	m.To = append(m.To, receiver)
}

func (m *Mail) AddCopier(copy string) {
	if m.CC == nil {
		m.CC = make([]string, 0, 10)
	}
	m.CC = append(m.CC, copy)
}

func (m *Mail) AddFile(file string) {
	if m.Files == nil {
		m.Files = make([]string, 0, 10)
	}
	m.Files = append(m.Files, file)
}

func (m *Mail) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("From: " + m.User + "\r\n")
	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")
	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")
	if len(m.CC) > 0 {
		buf.WriteString("C: " + strings.Join(m.CC, ",") + "\r\n")
	}

	//fix  Encode
	var coder = base64.StdEncoding
	var subject = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.Subject)) + "?="
	buf.WriteString("Subject: " + subject + "\r\n")

	//if len(m.ReplyTo) > 0 {
	//	buf.WriteString("Reply-To: " + m.ReplyTo + "\r\n")
	//}

	buf.WriteString("MIME-Version: 1.0\r\n")
	boundary := "f46d043c813270fc6b04c2d223da"

	if len(m.Files) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}

	buf.WriteString(fmt.Sprintf("Content-Type: %s; charset=utf-8\r\n\r\n", m.ContentType))
	buf.WriteString(m.Body)
	buf.WriteString("\r\n")

	if len(m.Files) > 0 {
		for _, f := range m.Files {
			buf.WriteString("\r\n\r\n--" + boundary + "\r\n")
			data, err := ioutil.ReadFile(f)
			if err != nil {
				return nil, err
			}
			_, filename := filepath.Split(f)
			if false {
				buf.WriteString("Content-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + filename + "\"\r\n\r\n")

				buf.Write(data)
			} else {
				ext := filepath.Ext(filename)
				mimetype := mime.TypeByExtension(ext)
				if mimetype != "" {
					mime := fmt.Sprintf("Content-Type: %s\r\n", mimetype)
					buf.WriteString(mime)
				} else {
					buf.WriteString("Content-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")

				buf.WriteString("Content-Disposition: attachment; filename=\"=?UTF-8?B?")
				buf.WriteString(coder.EncodeToString([]byte(filename)))
				buf.WriteString("?=\"\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
				base64.StdEncoding.Encode(b, data)

				// write base64 content in lines of up to 76 chars
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}

			buf.WriteString("\r\n--" + boundary)
		}

		buf.WriteString("--")
	}

	return buf.Bytes(), nil
}


func (m *Mail) Send() error {
	hp := strings.Split(m.Host, ":")
	auth := smtp.PlainAuth("", m.User, m.Pass, hp[0])

	data, err := m.Bytes()
	if err != nil {
		return err
	}
	return smtp.SendMail(m.Host, auth, m.User, m.To, data)
}


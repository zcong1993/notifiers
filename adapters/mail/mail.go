package mail

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"time"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/zcong1993/notifiers/types"
)

// Client is mail gun client
type Client struct {
	domain     string
	privateKey string
	mg         mailgun.Mailgun
	recipient  string
	from       string
}

// NewClient construct a mail gun notifier client
func NewClient(domain, privateKey, recipient, from string) *Client {
	return &Client{
		domain:     domain,
		privateKey: privateKey,
		mg:         mailgun.NewMailgun(domain, privateKey),
		recipient:  recipient,
		from:       from,
	}
}

// Notify impl notifier's notify
func (mc *Client) Notify(msg *types.Message) error {
	var bf bytes.Buffer
	err := msgTpl.Execute(&bf, msg)
	if err != nil {
		return err
	}
	message := mc.mg.NewMessage(mc.from, msg.Title, msg.Content, mc.recipient)
	message.SetHtml(bf.String())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, id, err := mc.mg.Send(ctx, message)
	if err != nil {
		return err
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}

var msgTpl = template.Must(template.New("mail").Parse(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional //EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd" />
<html lang="en" xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
  <head> </head>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta name="x-apple-disable-message-reformatting" />
    <!--[if !mso]><!-->
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <!--<![endif]-->
    <style type="text/css">
      * {
        text-size-adjust: 100%;
        -ms-text-size-adjust: 100%;
        -moz-text-size-adjust: 100%;
        -webkit-text-size-adjust: 100%;
      }

      html {
        height: 100%;
        width: 100%;
      }

      body {
        height: 100% !important;
        margin: 0 !important;
        padding: 0 !important;
        width: 100% !important;
        mso-line-height-rule: exactly;
      }

      div[style*="margin: 16px 0"] {
        margin: 0 !important;
      }

      table,
      td {
        mso-table-lspace: 0pt;
        mso-table-rspace: 0pt;
      }

      img {
        border: 0;
        height: auto;
        line-height: 100%;
        outline: none;
        text-decoration: none;
        -ms-interpolation-mode: bicubic;
      }

      .ReadMsgBody,
      .ExternalClass {
        width: 100%;
      }

      .ExternalClass,
      .ExternalClass p,
      .ExternalClass span,
      .ExternalClass td,
      .ExternalClass div {
        line-height: 100%;
      }
    </style>
    <!--[if gte mso 9]>
      <style type="text/css">
      li { text-indent: -1em; }
      table td { border-collapse: collapse; }
      </style>
      <![endif]-->
    <title>{{ .Title }}</title>
    <!-- content -->
    <!--[if gte mso 9]><xml>
       <o:OfficeDocumentSettings>
        <o:AllowPNG/>
        <o:PixelsPerInch>96</o:PixelsPerInch>
       </o:OfficeDocumentSettings>
      </xml><![endif]-->
  </head>
  <body class="body" style="margin: 0; width: 100%;">
    <table class="bodyTable" role="presentation" width="100%" align="left" border="0" cellpadding="0" cellspacing="0" style="width: 100%; margin: 0;">
      <tr>
        <td class="body__content" align="left" width="100%" valign="top" style="color: #000000; font-family: Helvetica,Arial,sans-serif; font-size: 16px; line-height: 20px;">
          <p class="text p" style="display: block; margin: 14px 0; color: #000000; font-family: Helvetica,Arial,sans-serif; font-size: 16px; line-height: 20px;">{{ .Content }}</p> {{ if .URL }}
          <div class="secondary-button button" style="background-color: #6495ED;">
            <table role="presentation" width="100%" align="left" border="0" cellpadding="0" cellspacing="0">
              <tr>
                <td>
                  <table role="presentation" width="auto" align="center" border="0" cellspacing="0" cellpadding="0" class="button__table" style="margin: 0 auto;">
                    <tr>
                      <td align="center" class="button__cell" style="border-radius: 3px; padding: 6px 12px; background-color: #6495ED;" bgcolor="#6495ED"><a href="{{ .URL }}" class="button__link" style="color: #FFFFFF; text-decoration: none; background-color: #6495ED; display: inline-block;"><span class="button__text" style="color: #FFFFFF; text-decoration: none;">阅读全文</span></a></td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
          </div> {{ end }} </td>
      </tr>
    </table>
    <div style="display:none; white-space:nowrap; font-size:15px; line-height:0;">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; </div>
  </body>
</html>
`))

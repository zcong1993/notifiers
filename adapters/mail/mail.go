package mail

import (
	"bytes"
	"context"
	"fmt"
	"text/template"
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

func (mc *Client) GetName() string {
	return "mail"
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
      .markdown {
        word-wrap: break-word;
      }
      .markdown,
      .markdown h1,
      .markdown h2,
      .markdown h3,
      .markdown h4,
      .markdown h5,
      .markdown h6,
      .markdown pre,
      .markdown code,
      .markdown blockquote,
      .markdown em,
      .markdown strong,
      .markdown code {
        font-size: 14px;
        line-height: 20px;
        font-weight: normal;
        font-style: normal;
        font-family: consolas, monaco, courier, "courier new", monospace;
        color: #333333;
      }
      .markdown h1,
      .markdown h2,
      .markdown h3,
      .markdown h4,
      .markdown h5,
      .markdown h6,
      .markdown pre,
      .markdown code,
      .markdown blockquote,
      .markdown ol,
      .markdown ul,
      .markdown li,
      .markdown p,
      .markdown section,
      .markdown header,
      .markdown footer {
        float: none;
        margin: 0;
        padding: 0;
      }
      .markdown h1,
      .markdown p,
      .markdown ul,
      .markdown ol,
      .markdown pre,
      .markdown blockquote {
        margin-top: 20px;
        margin-bottom: 20px;
      }
      .markdown h1 {
        position: relative;
        display: inline-block;
        display: table-cell;
        padding: 20px 0 40px;
        margin: 0;
        overflow: hidden;
      }
      .markdown h1:after {
        content: "====================================================================================================";
        position: absolute;
        bottom: 20px;
        left: 0;
      }
      .markdown h1 + * {
        margin-top: 0;
      }
      .markdown h2,
      .markdown h3,
      .markdown h4,
      .markdown h5,
      .markdown h6 {
        position: relative;
        margin-bottom: 20px;
      }
      .markdown h2:before,
      .markdown h3:before,
      .markdown h4:before,
      .markdown h5:before,
      .markdown h6:before {
        content: "## ";
        display: inline;
      }
      .markdown h3:before {
        content: "### ";
      }
      .markdown h4:before {
        content: "#### ";
      }
      .markdown h5:before {
        content: "##### ";
      }
      .markdown h6:before {
        content: "###### ";
      }
      .markdown li {
        position: relative;
        display: block;
        padding-left: 34px;
        padding-left: 4ch;
      }
      .markdown li:after {
        position: absolute;
        top: 0;
        left: 0;
      }
      .markdown ul > li:after {
        content: "*";
      }
      .markdown ol {
        counter-reset: ol;
      }
      .markdown ol > li:after {
        content: counter(ol) ".";
        counter-increment: ol;
      }
      .markdown pre {
        margin-left: 34px;
        margin-left: 4ch;
      }
      .markdown blockquote {
        position: relative;
        padding-left: 17px;
        padding-left: 2ch;
        overflow: hidden;
      }
      .markdown blockquote:after {
        content: ">\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>\A>";
        white-space: pre;
        position: absolute;
        top: 0;
        left: 0;
        font-size: 14px;
        line-height: 20px;
      }
      .markdown strong:before,
      .markdown strong:after {
        content: "__";
        display: inline;
      }
      .markdown em:before,
      .markdown em:after {
        content: "*";
        display: inline;
      }
      .markdown a {
        text-decoration: none;
      }
      .markdown a:before {
        content: "[";
        display: inline;
        color: #333333;
      }
      .markdown a:after {
        content: "](" attr(href) ")";
        display: inline;
        color: #333333;
      }
      .markdown code {
        font-weight: 100;
      }
      .markdown code:before,
      .markdown code:after {
        content: "` + "`" + `";
        display: inline;
      }
      .markdown pre code:before,
      .markdown pre code:after {
        content: none;
      }
      .markdown hr {
        position: relative;
        height: 20px;
        font-size: 0;
        line-height: 0;
        overflow: hidden;
        border: 0;
        margin-bottom: 20px;
      }
      .markdown hr:after {
        content: "----------------------------------------------------------------------------------------------------";
        position: absolute;
        top: 0;
        left: 0;
        font-size: 14px;
        line-height: 20px;
        width: 100%;
        word-wrap: break-word;
      }
      @-moz-document url-prefix() {
        .markdown h1 {
          display: block;
        }
      }
      .markdown-ones ol > li:after {
        content: "1.";
      }
    </style>
    <!--[if gte mso 9]>
      <style type="text/css">
      li { text-indent: -1em; }
      table td { border-collapse: collapse; }
      </style>
      <![endif]-->
    <title>Welcome to HEML!</title>
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
          <div class="markdown container" style="margin: 0 auto; max-width: 600px; width: 100%; word-wrap: break-word; color: #333333; font-family: consolas,monaco,courier,'courier new',monospace; font-size: 14px; font-style: normal; font-weight: 400; line-height: 20px;">
          <!--[if mso | IE]>
            <table class="container__table__ie" role="presentation" border="0" cellpadding="0" cellspacing="0" style="margin-right: auto; margin-left: auto;width: 600px" width="600" align="center">
              <tr>
                <td> <![endif]-->
                  <table class="container__table" role="presentation" border="0" align="center" cellpadding="0" cellspacing="0" width="100%">
                    <tr class="container__row">
                      <td class="container__cell" width="100%" align="left" valign="top"> {{ .Content }} {{ if .URL }} <a href="{{ .URL }}" class="a" style="text-decoration: none;"><span class="a__text" style="text-decoration: none;">{{ .URL }}</span></a> {{ end }} </td>
                    </tr>
                  </table> <!--[if mso | IE]> </td>
              </tr>
            </table> <![endif]--> </div>
        </td>
      </tr>
    </table>
    <div style="display:none; white-space:nowrap; font-size:15px; line-height:0;">&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; </div>
  </body>
</html>`))

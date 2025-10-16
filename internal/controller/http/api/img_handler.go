package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

func renderHTMLToPNG(htmlContent string) ([]byte, error) {
	var buf []byte
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	// Запускаем задачи: открыть about:blank, вставить HTML через javascript, затем сделать снимок
	err := chromedp.Run(ctx,
		chromedp.EmulateViewport(500, 280),
		chromedp.Navigate("about:blank"),
		chromedp.Evaluate(`document.open();document.write(`+strconv.Quote(htmlContent)+`);document.close();`, nil),
		chromedp.Sleep(50*time.Millisecond), // ожидание рендера
		chromedp.FullScreenshot(&buf, 90),
	)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (h *Handler) getImgTicketStatus(c *gin.Context) {
	now := time.Now()
	formatted := now.Format("02.01.2006 15:04")
	statusMap := map[int]string{
		1: "Новая",
		2: "В работе (назначена)",
		3: "В работе (запланирована)",
		4: "Ожидающая",
		5: "Решена",
		6: "Закрыта",
	}
	statusClassMap := map[int]string{
		1: "pending",
		2: "pending",
		3: "pending",
		4: "pending",
		5: "completed",
		6: "closed",
	}
	groupCaption := "Текущая группа"
	ticketId := c.Param("id")
	ticket, err := h.uc.GetGLPITicketSimple(ticketId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Заявка не найдена", "error": err.Error()})
		return
	}
	if ticket.GroupName == "" {
		groupCaption = "Не назначена ни одной группе"
	}
	statusText, ok := statusMap[ticket.Status]
	if !ok {
		statusText = "Неизвестный статус"
	}
	statusClass, ok := statusClassMap[ticket.Status]
	if !ok {
		statusText = "pending"
	}

	htmlTemplate := `<html>
<head>
  <meta charset="UTF-8" />
  <style>
    body {
      margin: 0;
      background: #f5f5f5;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
        Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
      padding: 8px;
      box-sizing: border-box;
    }
    .card {
      width: 100%;
      height: 100%;
      background: #fff;
      border-radius: 4px;
      box-shadow: 0 2px 5px rgba(0,0,0,0.15);
      padding: 12px 18px 5px 18px;
      box-sizing: border-box;
      display: flex;
      flex-direction: column;
      position: relative;
    }
    .title {
      font-size: 24px;
      font-weight: 500;
      margin-bottom: 8px;
      color: #212121;
    }
    .status {
      display: inline-block;
      padding: 6px 16px;
      border-radius: 6px;
      font-size: 14px;
      font-weight: 500;
      text-transform: uppercase;
      letter-spacing: 1.1px;
      color: white;
      margin-bottom: 16px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      width: fit-content;
    }
    .status.completed {
      background-color: #4caf50;
      box-shadow: 0 4px 8px rgba(76, 175, 80, 0.4);
    }
    .status.pending {
      background-color:rgb(214, 130, 4);
      box-shadow: 0 4px 8px rgba(255, 152, 0, 0.4);
    }
	      .status.closed {
      background-color:rgb(90, 82, 70);
      box-shadow: 0 4px 8px rgba(37, 34, 30, 0.4);
    }
    .description-group {
      display: flex;
      flex-direction: column;
      align-items: flex-start;
      gap: 10px;
      margin-bottom: 40px;
      flex-grow: 1;
      width: 100%;
    }
    .description {
      font-size: 14px;
      color: #424242;
      line-height: 1.5;
      width: 100%;
    }
    .group {
      font-size: 14px;
      color: #9e9e9e;
      text-align: right;
      width: 100%;
    }
    .group .group-title {
      color: #646464;
      font-weight: 600;
      font-size: 14px;
      margin-bottom: 4px;
    }
    .footer2 {
      font-size: 12px;
      color:rgb(60, 78, 60);
      text-align: left;
      margin-top: auto;
	  padding-bottom:5px;
    }
    .footer {
      font-size: 12px;
      color: #9e9e9e;
      text-align: center;
      margin-top: auto;
	  padding-bottom:5px;
    }

    .prefix {
  font-size: 16px;        /* Меньше для текста "Заявка №" */
  color: #444444;
  text-transform: none;
letter-spacing: 4px;
  letter-spacing: normal;
}

.number {
  font-size: 20px;        /* Большой размер для номера */
  color: #777755;         /* Серый цвет */
  margin-left: 8px;       /* Немного отступа от префикса */
}
.little-number  {
  font-size: 10px;
}
  .little-grey-text  {
  font-size: 12px;
   color: #777755;
   padding-right:5px;
}
  </style>
</head>
<body>
  <div class="card">
      <div class="title">
      <span class="prefix">Заявка №</span><span class="number">` + strconv.Itoa(ticket.Id) + `</span>
     </div>
    <div class="status ` + statusClass + `">` + statusText + `</div>
    <div class="description-group">
      <div class="description">
        ` + ticket.Name + `
      </div>
      <div class="group">
        <div class="group-title">` + groupCaption + `</div>
        ` + ticket.GroupName + `
      </div>
    </div>
    <div class="footer2"><span class="little-grey-text">Организация:</span>` + ticket.Company + `</div>
    <div class="footer">Создано в системе support.rw, <span class="little-number">` + formatted + `</span></div>
  </div>
</body>
</html>
`

	pngBuf, err := renderHTMLToPNG(htmlTemplate)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Ошибка получения статуса заявки", "error": err.Error()})
		return
	}
	c.Header("Content-Type", "image/png")
	c.Status(http.StatusOK)
	c.Writer.Write(pngBuf)
}

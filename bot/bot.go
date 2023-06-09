package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"lunch_helper/constant"
	"lunch_helper/util"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type BotClient interface {
	ParseRequest(r *http.Request) ([]*linebot.Event, error)
	PushText(userID, text string)
	PushFlex(userID string, altText string, flexContainer linebot.FlexContainer)
	SendText(replyToken, text string)
	SendFlex(replyToken string, altText string, flexContainer linebot.FlexContainer)
	SendTextWithQuickReplies(replyToken, text string, replyItems *linebot.QuickReplyItems)
	GetProfile(userID string) *linebot.GetProfileCall
}

type LineBotClient struct {
	*linebot.Client
	channelToken string
}

type RichMenu_New struct {
	Size        linebot.RichMenuSize `json:"size"`
	Selected    bool                 `json:"selected"`
	Name        string               `json:"name"`
	ChatBarText string               `json:"chatBarText"`
	Areas       []AreaDetail_New     `json:"areas"`
}

type AreaDetail_New struct {
	Bounds linebot.RichMenuBounds `json:"bounds"`
	Action RichMenuAction_New     `json:"action"`
}

type RichMenuAction_New struct {
	Type            linebot.RichMenuActionType `json:"type"`
	URI             string                     `json:"uri,omitempty"`
	Text            string                     `json:"text,omitempty"`
	DisplayText     string                     `json:"displayText,omitempty"`
	Label           string                     `json:"label,omitempty"`
	Data            string                     `json:"data,omitempty"`
	Mode            string                     `json:"mode,omitempty"`
	Initial         string                     `json:"initial,omitempty"`
	Max             string                     `json:"max,omitempty"`
	Min             string                     `json:"min,omitempty"`
	RichMenuAliasID string                     `json:"richMenuAliasId,omitempty"`
	// 補上下面幾個新的
	InputOption string `json:"inputOption,omitempty"`
}

type CreateRichMenuResponse struct {
	RichMenuID string `json:"richMenuId"`
}

type RichMenuRequest struct {
	ImagePath      string
	AliasName      constant.RichMenuAliasName
	RichMenuStruct RichMenu_New
}

func NewBotClient(channelSecret, channelToken string) (*LineBotClient, error) {
	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		return nil, err
	}

	return &LineBotClient{bot, channelToken}, nil
}

func (bc *LineBotClient) SetWebHookUrl(apiBaseUrl, endPoint string) error {
	url, err := util.BindUrl(apiBaseUrl, endPoint)
	if err != nil {
		return errors.New("bind url error: " + err.Error())
	}
	_, err = bc.Client.SetWebhookEndpointURL(url).Do()
	if err != nil {
		return errors.New("SetWebHookUrl Error: " + err.Error())
	}
	return nil
}

func (bc *LineBotClient) SetupRichMenu(richMenuRequest ...RichMenuRequest) error {
	richMenuIdList := []string{}
	for _, request := range richMenuRequest {
		resp, err := bc.CreateRichMenu_New(request.RichMenuStruct)
		if err != nil {
			return err
		}
		log.Println(resp.RichMenuID)

		richMenuIdList = append(richMenuIdList, resp.RichMenuID)

		_, err = bc.UploadRichMenuImage(resp.RichMenuID, request.ImagePath).Do()
		if err != nil {
			return err
		}
	}

	// set default richmenu
	_, err := bc.SetDefaultRichMenu(richMenuIdList[0]).Do()
	if err != nil {
		return err
	}

	// create richmenu alias
	for index, request := range richMenuRequest {
		_, err := bc.CreateRichMenuAlias(string(request.AliasName), richMenuIdList[index]).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (bc *LineBotClient) ResetRichMenu() []error {
	aliasList, err := bc.GetRichMenuAliasList().Do()
	if err != nil {
		return []error{err}
	}

	errList := []error{}
	for _, alias := range aliasList {
		_, err := bc.DeleteRichMenuAlias(alias.RichMenuAliasID).Do()
		if err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) > 0 {
		return errList
	}

	richMenu, err := bc.GetRichMenuList().Do()
	if err != nil {
		return []error{err}
	}

	for _, rich := range richMenu {
		_, err := bc.DeleteRichMenu(rich.RichMenuID).Do()
		if err != nil {
			errList = append(errList, err)
		}
	}
	if len(errList) > 0 {
		return errList
	}

	return []error{}
}

// 原始createRichMenu內RichMenuAction沒有提供inputOption，因此重寫function
func (bc *LineBotClient) CreateRichMenu_New(richMenu RichMenu_New) (*CreateRichMenuResponse, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(richMenu)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.line.me/v2/bot/richmenu",
		&buf,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Bearer "+bc.channelToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	source, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	response := CreateRichMenuResponse{}
	if err := json.Unmarshal(source, &response); err != nil {
		return nil, fmt.Errorf("Unmarshal response body failed:", err)
	}

	return &response, nil
}

func (bc *LineBotClient) SendText(replyToken, text string) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text),
	).Do()
	if err != nil {
		log.Printf("SendText Error: %s", err)
	}
}

func (bc *LineBotClient) PushText(userID string, text string) {
	_, err := bc.PushMessage(userID, linebot.NewTextMessage(text)).Do()
	if err != nil {
		log.Printf("PushText Error: %s", err)
	}
}

func (bc *LineBotClient) SendTextWithQuickReplies(replyToken, text string, replyItems *linebot.QuickReplyItems) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(text).WithQuickReplies(replyItems),
	).Do()
	if err != nil {
		log.Printf("SendTextWithQuickReplies Error: %s", err)
	}
}

func (bc *LineBotClient) SendFlex(replyToken string, altText string, flexContainer linebot.FlexContainer) {
	_, err := bc.ReplyMessage(
		replyToken,
		linebot.NewFlexMessage(altText, flexContainer),
	).Do()
	if err != nil {
		log.Printf("SendFlex Error: %s", err)
	}
}

func (bc *LineBotClient) PushFlex(userID string, altText string, flexContainer linebot.FlexContainer) {
	_, err := bc.PushMessage(userID, linebot.NewFlexMessage(altText, flexContainer)).Do()
	if err != nil {
		log.Printf("PushFlex Error: %s", err)
	}
}

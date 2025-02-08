package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/cyclex/ambpi-core/api"
	"github.com/cyclex/ambpi-core/domain"
	"github.com/cyclex/ambpi-core/domain/model"
	"github.com/cyclex/ambpi-core/domain/repository"
	"github.com/cyclex/ambpi-core/pkg"
	"github.com/cyclex/ambpi-core/pkg/httprequest"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/google/uuid"
)

type chatUcase struct {
	m               repository.ModelRepository
	q               repository.QueueRepository
	urlSendMsg      string
	PhoneID         string
	AccessToken     string
	NumZonk         int
	NumLimitRedeem  int
	Rdb             *redis.Client
	UrlPush         string
	AccountID       string
	DivisionID      string
	AccessTokenPush string
}

var (
	berkesempatan     = "Selamat. kamu berkesempatan mendapatkan hadiah langsung %s dan Nomor Undian.\n\nData kamu akan kami validasi terlebih dahulu dalam 1x24 jam. Untuk memastikan kamu mendapatkan hadiah langsung %s dan nomor undian\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	zonk              = "Mohon maaf, kamu belum beruntung mendapatkan hadiah langsung. Namun kamu berkesempatan mendapatkan Nomor Undian yang akan diundi di akhir periode\n\nData kamu akan kami validasi terlebih dahulu dalam 1x24 jam. Untuk memastikan kamu mendapatkan nomor undian\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	programNotActive  = "Mohon Maaf Program Consumer AM 2025 belum dimulai.\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	programEnded      = "Mohon Maaf Program Consumer AM 2025 sudah berakhir.\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	invalidFormat     = "Maaf format pesan yang anda kirimkan tidak lengkap silahkan kirim foto struk pembelian dan caption\n*AM#NamaLengkap#NIK#Kota#Profesi (tukang/mandor/end konsumen)*\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	invalidPhoto      = "Maaf format foto yang kamu kirimkan salah, kirimkan hanya format image jpeg/jpg\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	invalidCaption    = "Maaf format caption foto yang kamu kirimkan salah. Format caption : *AM#NamaLengkap#NIK#Kota#Profesi (tukang/mandor/end konsumen)*\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	invalidNIK        = "Maaf NIK yang kamu kirimkan salah. Mohon kirim NIK yang sesuai.\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	invalidProfession = "Maaf profesi yang kamu masukan salah. Masukan profesi yang sesuai *(tukang/mandor/end konsumen)*"
	invalid           = "Untuk mengikuti program ini silahkan mengirimkan foto struk pembelian disertai caption *AM#NamaLengkap#NIK#Kota#Profesi (tukang/mandor/end konsumen)*\n\nInfo lebih lanjut: klik wa.me/6281225772215\nCS Office Hour (Senin s/d Jumat 08.00 - 17.30 WIB)\n\nHati-hati penipuan! Consumer Promo AM tidak dipungut biaya apapun!"
	maintenance       = "Mohon maaf, system Consumer promo sedang maintenance. Kamu dapat mengirimkan kembali struk pembelian beserta captionnya dalam 1x24 Jam kedepan\n\ninformasi lebih lanjut, hubungi CS :\nwa.me/6281225772215"

	// TODO replace system error to invalid
	systemError = "System error"
)

func NewChatUcase(m repository.ModelRepository, urlPush, urlSendMsg, phoneID, accessToken, accountID, divisionID, accessTokenPush string, queueRepo repository.QueueRepository, rdb *redis.Client) domain.ChatUcase {

	return &chatUcase{
		m:               m,
		urlSendMsg:      urlSendMsg,
		PhoneID:         phoneID,
		AccessToken:     accessToken,
		q:               queueRepo,
		Rdb:             rdb,
		UrlPush:         urlPush,
		AccountID:       accountID,
		DivisionID:      divisionID,
		AccessTokenPush: accessTokenPush,
	}
}

func (self *chatUcase) IncomingMessages(payload api.CproMessage) (trxChatBotID string, err error) {

	waID := payload.From
	sessID := uuid.NewString()

	isRedeem, outgoing, err := self.ReplyMessages(waID, payload)
	if isRedeem {
		return
	}

	res, statusCode, err := self.ChatToUserCoster(waID, outgoing, "text", "")
	if err != nil {
		err = errors.Wrap(err, "[usecase.IncomingMessages] ChatToUserCoster")
		return
	}

	if statusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("%s: %s", http.StatusText(statusCode), res))
		err = errors.Wrap(err, "[usecase.IncomingMessages]")
		return
	}

	fmt.Printf("%s", res)

	clog := model.ConversationsLog{
		SessionID: sessID,
		Incoming:  payload.Text.Body,
		WAID:      waID,
		Outgouing: outgoing,
	}
	err = self.m.CreateConversationsLog(clog)
	if err != nil {
		err = errors.Wrap(err, "[usecase.IncomingMessages] CreateConversationsLog")
	}

	return
}

func (self *chatUcase) ReplyMessages(waID string, payload api.CproMessage) (isRedeem bool, outgoing string, err error) {

	var isImage, isJPG bool

	if payload.Type == "image" {
		isImage = true
	}

	if !isImage {
		return false, invalid, nil
	}

	isActive := self.m.IsProgramActive("am")
	switch isActive {
	case 1:
		return false, programNotActive, nil
	case 2:
		return false, programEnded, nil
	default:
	}

	incoming := strings.ToLower(payload.Image.Caption)
	if len(incoming) == 0 {
		return false, invalidFormat, nil
	}

	if payload.Image.MimeType == "image/jpeg" {
		isJPG = true
	}

	if !isJPG {
		return false, invalidPhoto, nil
	}

	// AM#NamaLengkap#NIK#Kota#Profesi
	str := strings.Split(incoming, "#")
	if len(str) != 5 || str[0] != "am" {
		return false, invalidCaption, nil
	}

	if len(str[2]) != 16 || !pkg.IsInt(str[2]) || !pkg.IsLetter(str[1]) || !pkg.IsLetter(str[3]) {
		return false, invalidNIK, nil
	}

	var profession bool
	switch str[4] {
	case "tukang", "mandor", "end konsumen":
		profession = true
		break
	}

	if !profession {
		return false, invalidProfession, nil
	}

	msg := api.PayloadRedeem{
		WaID:          waID,
		Name:          str[1],
		Profession:    str[4],
		IsFormatValid: true,
		Raw:           incoming,
		NIK:           str[2],
		County:        str[3],
	}

	err = self.q.CreateQueueRedeem(model.QueueRedeem{
		ID:       primitive.NewObjectID(),
		TrxId:    uuid.NewString(),
		Messages: msg,
		State:    1,
	})

	if err != nil {
		err = errors.Wrap(err, "[usecase.ReplyMessages] CreateQueueRedeem")
		return false, invalid, err
	}

	return true, "", nil
}

func (self *chatUcase) DoRedeem(usUcode model.UsersUniqueCode) (outgoing string, err error) {

	var (
		cond  map[string]interface{}
		reply = usUcode.Raw
	)

	cond = map[string]interface{}{}

	p, err := self.m.FindActivePrizes(cond, true)
	if err != nil {
		if err.Error() != "record not found" {
			err = errors.Wrap(err, fmt.Sprintf("[usecase.DoRedeem] FindActivePrizes. cond:%+v, isActive:%v (waID:%s, reply:%s)", cond, true, usUcode.WaID, reply))
			return systemError, err
		}
		return maintenance, nil
	}

	outgoing, err = self.setAndCreate(p, usUcode)
	if err != nil {
		err = errors.Wrap(err, "[usecase.DoRedeem] setAndCreate")
		return systemError, err
	}

	return outgoing, nil

}

func (self *chatUcase) setAndCreate(p model.Prizes, u model.UsersUniqueCode) (outgoing string, err error) {

	var (
		cond map[string]interface{}
	)

	cond = map[string]interface{}{
		"id": p.ID,
	}

	updated := map[string]interface{}{
		"updated_at": time.Now().Local(),
		"is_used":    true,
	}

	_, err = self.m.SetPrizes(cond, updated)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[usecase.setAndCreate] SetPrizes. cond:%+v, (waID:%s, outgoing:%s)", cond, u.WaID, systemError))
		return systemError, err
	}

	outgoing = fmt.Sprintf(berkesempatan, p.Prize, p.Prize)
	if p.PrizeType == "zonk" {
		outgoing = zonk
		u.IsZonk = true
	}
	u.Reply = outgoing
	u.CreatedAt = time.Now().Local()
	_, err, pk := self.m.CreateUsersUniqueCodes(u)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[usecase.setAndCreate] CreateUsersUniqueCodes. (waID:%s, outgoing:%s)", u.WaID, systemError))
		return systemError, err
	}

	_, err = self.m.CreateRedeemPrizes(model.RedeemPrizes{
		PrizeID:           p.ID,
		UsersUniqueCodeID: pk,
		Msisdn:            u.WaID,
	})
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("[usecase.setAndCreate] CreateRedeemPrizes. (waID:%s, outgoing:%s)", u.WaID, systemError))
		return systemError, err
	}

	return outgoing, nil
}

func (self *chatUcase) ChatToUserCoster(waID, chat, types, media string) (res []byte, statusCode int, err error) {

	var payload interface{}
	url := self.urlSendMsg

	payload = api.ReqSendMessageText{
		XID:         uuid.NewString(),
		ChannelID:   "whatsapp-cloud",
		AccountID:   self.AccountID,
		DivisionID:  self.DivisionID,
		IsHelpdesk:  false,
		MessageType: "outbound",
		Data: api.Data{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               waID,
			Type:             types,
			Text: api.Text{
				Body: chat,
			},
		},
	}

	res, statusCode, err = httprequest.PostJson(url, payload, 15*time.Second, self.AccessToken)
	if err != nil {
		err = errors.Wrap(err, "[usecase.ChatToUserCoster] PostJson")
	}
	return
}

func (self *chatUcase) ChatToUser(waID string, chats []string, types, templateName string) (res []byte, statusCode int, err error) {
	const (
		TypePush    = "push"
		TypeChat    = "chat"
		AuthBearer  = "Bearer "
		ContentType = "application/json"
	)

	var payloadBytes []byte
	url := fmt.Sprintf("%s/v1/%s/messages", self.UrlPush, self.AccountID)
	param := []api.CproParameter{}

	for _, chat := range chats {
		param = append(param, api.CproParameter{Type: "text", Text: chat})
	}
	switch types {
	case TypePush:
		payload := api.CproPayloadPush{
			XID:  fmt.Sprintf("ambpi-%d", time.Now().UnixMilli()),
			To:   waID,
			Type: "template",
			Template: api.CproTemplate{
				Name: templateName,
				Language: api.CproLanguage{
					Code:   "id",
					Policy: "deterministic",
				},
				Components: []api.CproComponent{
					{
						Type:       "body",
						Parameters: param,
					},
				},
			},
		}
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, 0, pkg.WrapError(err, "[usecase.ChatToUser] Marshal")
		}

	default:
		return nil, 0, fmt.Errorf("[usecase.ChatToUser] unsupported type: %s", types)
	}

	accessToken := self.AccessTokenPush
	// accessToken, err := self.getAccessToken(context.Background())
	// if err != nil {
	// 	return nil, 0, pkg.WrapError(err, "[usecase.ChatToUser] getAccessToken")
	// }

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, 0, pkg.WrapError(err, "[usecase.ChatToUser] NewRequest")
	}

	req.Header.Add("Authorization", AuthBearer+accessToken)
	req.Header.Add("Content-Type", ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, pkg.WrapError(err, "[usecase.ChatToUser] Do")
	}
	defer resp.Body.Close()

	res, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, pkg.WrapError(err, "[usecase.ChatToUser] ReadAll")
	}

	fmt.Println(AuthBearer + accessToken)
	fmt.Printf("%s\n", payloadBytes)
	fmt.Println(url)
	return res, resp.StatusCode, nil
}

func (self *chatUcase) getAccessToken(ctx context.Context) (string, error) {
	token, err := self.GetToken(ctx)
	if err != nil {
		res, _, refreshErr := self.RefreshToken()
		if refreshErr != nil {
			return "", errors.Wrap(refreshErr, "[usecase.getAccessToken] RefreshToken")
		}
		return res.Token, nil
	}
	return token, nil
}

func (self *chatUcase) CreateQueueReply(data api.PayloadReply) (err error) {

	payload := model.QueueReply{
		ID:       primitive.NewObjectID(),
		TrxId:    uuid.New().String(),
		State:    1,
		Messages: data,
	}

	err = self.q.CreateQueueReply(payload)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateQueueReply] CreateQueueReply")
	}

	return
}

func (self *chatUcase) CreateConversationsLog(data model.ConversationsLog) (err error) {

	err = self.m.CreateConversationsLog(data)
	if err != nil {
		err = errors.Wrap(err, "[usecase.CreateConversationsLog] CreateConversationsLog")
	}

	return
}

func (self *chatUcase) RefreshToken() (apiResp api.RefreshResponse, statusCode int, err error) {

	url := "https://graph.jatismobile.com/login"

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		err = errors.Wrap(err, "[usecase.RefreshToken] NewRequest")
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", self.AccessTokenPush))

	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "[usecase.RefreshToken] Do")
		return
	}

	statusCode = resp.StatusCode
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "[usecase.RefreshToken] ReadAll")
		return
	}

	if err = json.Unmarshal(res, &apiResp); err != nil {
		err = errors.Wrap(err, "[usecase.RefreshToken] Unmarshal")
		return
	}

	return
}

func (self *chatUcase) SetToken(ctx context.Context, token string, ttl time.Duration) (err error) {

	self.Rdb.Set(ctx, "token", token, ttl).Err()
	if err != nil {
		err = errors.Wrap(err, "[usecase.SetToken] Set")
	}

	return
}

func (self *chatUcase) GetToken(ctx context.Context) (token string, err error) {

	token, err = self.Rdb.Get(ctx, "token").Result()
	if err != nil {
		err = errors.Wrap(err, "[usecase.GetToken] Get")
	}

	return
}

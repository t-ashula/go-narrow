package narrow

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// A Client fetch data from novel api
type Client struct {
	httpClient *http.Client
}

var userAgent = fmt.Sprintf("go-narrow/%s", Version)

// NewClient returns new novel api client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// Search returns narou API result
func (c *Client) Search(ctx context.Context, params *SearchParams) (*SearchResult, error) {
	u, err := params.ToURL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("user-agent", userAgent)

	req = req.WithContext(ctx)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result, err := parseResponse(body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseResponse(body []byte) (*SearchResult, error) {
	var responses []response
	err := json.Unmarshal(body, &responses)
	if err != nil {
		return nil, err
	}

	res := &SearchResult{
		AllCount:   *responses[0].AllCount,
		NovelInfos: toNovelInfos(responses[1:]),
	}
	return res, nil
}

func toNovelInfos(responses []response) []NovelInfo {
	novels := make([]NovelInfo, len(responses))
	for i, res := range responses {
		novels[i] = res.toNoevlInfo()
	}
	return novels
}

func (res *response) toNoevlInfo() NovelInfo {
	info := NovelInfo{}
	info.Title = res.Title
	info.NCode = res.NCode
	if res.UserID != nil {
		info.UserID = strp(fmt.Sprintf("%d", *res.UserID))
	}
	info.Writer = res.Writer
	info.Story = res.Story
	info.BigGenre = res.BigGenre
	info.Genre = res.Genre
	if res.Keyword != nil {
		info.Keywords = strings.Split(*res.Keyword, " ")
	} else {
		info.Keywords = nil
	}
	if res.GeneralFirstUp != nil {
		info.GeneralFirstUp = &(*res.GeneralFirstUp).Time
	}
	if res.GeneralLastUp != nil {
		info.GeneralLastUp = &(*res.GeneralLastUp).Time
	}
	info.NovelType = res.NovelType
	info.End = res.End
	info.GeneralAllNo = res.GeneralAllNo
	info.Length = res.Length
	info.Time = res.Time
	info.IsStop = nilOrTrue(res.IsStop)
	info.IsR15 = nilOrTrue(res.IsR15)
	info.IsBoysLove = nilOrTrue(res.IsBL)
	info.IsGirlsLove = nilOrTrue(res.IsGL)
	info.IsZankoku = nilOrTrue(res.IsZankoku)
	info.IsTensei = nilOrTrue(res.IsTensei)
	info.IsTenni = nilOrTrue(res.IsTenni)
	info.PCOrK = res.PcOrK
	info.GlobalPoint = res.GlobalPoint
	info.FavNovelCount = res.FavNovelCnt
	info.ReviewCount = res.ReviewCnt
	info.AllPoint = res.AllPoint
	info.AllHyokaCount = res.AllHyokoCnt
	info.SasieCount = res.SasieCnt
	info.KaiwaRitu = res.KaiwaRitu
	if res.NovelUpdatedAt != nil {
		info.NovelUpdatedAt = &(*res.NovelUpdatedAt).Time
	}
	if res.UpdatedAt != nil {
		info.UpdatedAt = &(*res.UpdatedAt).Time
	}
	return info
}

type response struct {
	// number of all matched novels
	AllCount *int `json:"allcount"`

	// novel title
	Title *string `json:"title"`
	// N-code
	NCode *string `json:"ncode"`
	// user id
	UserID *int `json:"userid"`
	// Writer Name
	Writer *string `json:"writer"`
	// あらすじ
	Story *string `json:"story"`
	// Large category
	BigGenre *int `json:"biggenre"`
	// genre
	Genre *int `json:"genre"`
	// based on (not used)
	Gensaku *string `json:"gensaku"`
	// Keyword (space separated)
	Keyword *string `json:"keyword"`
	// 初回掲載日
	GeneralFirstUp *novelTime `json:"general_firstup"` // `YYYY-MM-DD HH:MM:SS` in JST
	// 最終掲載日
	GeneralLastUp *novelTime `json:"general_lastup"` // `YYYY-MM-DD HH:MM:SS` in JST
	// 連載 短編
	NovelType *int `json:"novel_type"` // 1:serialized, 2:short story
	// 未完結
	End *int `json:"end"` // 0:finished, 1:running
	// 全掲載部分数です。短編の場合は1です。
	GeneralAllNo *int `json:"general_all_no"`
	// 小説文字数
	Length *int `json:"length"`
	// 読了時間（分）小説文字数÷500を切り上げした数値
	Time *int `json:"time"`
	// 長期連載停止中なら1、それ以外は0です。
	IsStop *int `json:"isstop"`
	// 登録必須キーワードに「R15」が含まれる場合は1、それ以外は0です。
	IsR15 *int `json:"isr15"`
	// 登録必須キーワードに「ボーイズラブ」が含まれる場合は1、それ以外は0です。
	IsBL *int `json:"isbl"`
	// 登録必須キーワードに「ガールズラブ」が含まれる場合は1、それ以外は0です。
	IsGL *int `json:"isgl"`
	// 登録必須キーワードに「残酷な描写あり」が含まれる場合は1、それ以外は0です。
	IsZankoku *int `json:"iszankoku"`
	// 登録必須キーワードに「異世界転生」が含まれる場合は1、それ以外は0です。
	IsTensei *int `json:"istensei"`
	// 登録必須キーワードに「異世界転移」が含まれる場合は1、それ以外は0です。
	IsTenni *int `json:"istenni"`
	// 1はケータイのみ、2はPCのみ、3はPCとケータイで投稿された作品です。
	PcOrK *int `json:"pc_or_k"`
	// 総合評価ポイント(=(ブックマーク数×2)+評価点)
	GlobalPoint *int `json:"global_point"`
	//ブックマーク数
	FavNovelCnt *int `json:"fav_novel_cnt"`
	// レビュー数
	ReviewCnt *int `json:"review_cnt"`
	// 評価点
	AllPoint *int `json:"all_point"`
	// 評価者数
	AllHyokoCnt *int `json:"all_hyoka_cnt"`
	// 挿絵の数
	SasieCnt *int `json:"sasie_cnt"`
	// 会話率
	KaiwaRitu *int `json:"kaiwaritu"`
	// 小説の更新日時
	NovelUpdatedAt *novelTime `json:"novelupdated_at"`
	// 最終更新日時 (注意：システム用で小説更新時とは関係ありません)
	UpdatedAt *novelTime `json:"updated_at"`
}

type novelTime struct {
	time.Time
}

const novelTimeFormat = "2006-01-02 15:04:05"

func (nt *novelTime) UnmarshalJSON(b []byte) (err error) {
	nt.Time = time.Time{}
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return
	}
	nt.Time, err = time.ParseInLocation(novelTimeFormat, s, loc)
	return
}

func intp(val int) *int       { return &val }
func strp(str string) *string { return &str }
func boolp(b bool) *bool      { return &b }
func nilOrTrue(ip *int) *bool {
	if ip == nil {
		return nil
	}
	return boolp(*ip == 1)
}

package narrow

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_parseResponse(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *SearchResult
		wantErr bool
	}{
		{"normal response",
			args{[]byte(`[
			{
			 "allcount":12345
			},{
			 "title":"\u5F92\u7136\u8349",
			 "ncode":"N0000AA","userid":1234567,
			 "writer":"\u5409\u7530\u517C\u597D",
			 "story":"\u3064\u308C\u3065\u308C\u306A\u308B\u307E\u307E\u306B\u3072\u3050\u3089\u3057\u786F\u306B\u5411\u304B\u3044\u3066",
			 "biggenre":99,"genre":9903,"gensaku":"",
			 "keyword":"R15 \u5e74\u306e\u5dee \u5192\u967a",
			 "general_firstup":"2019-05-06 18:39:05","general_lastup":"2019-08-16 08:35:52",
			 "novel_type":1,"end":1,
			 "general_all_no":45,"length":130245,"time":261,"isstop":0,
			 "isr15":1,"isbl":0,"isgl":0,"iszankoku":1,"istensei":0,"istenni":0,
			 "pc_or_k":2,"global_point":154,"fav_novel_cnt":33,"review_cnt":0,"all_point":88,"all_hyoka_cnt":10,"sasie_cnt":0,
			 "kaiwaritu":43,"novelupdated_at":"2019-08-16 08:35:52","updated_at":"2019-08-16 08:38:16"
			},{
			 "title":"\u6E90\u6C0F\u7269\u8A9E",
			 "ncode":"N9999ZZ","userid":9876543,"writer":"purple",
			 "story":"\u3044\u3065\u308C\u306E\u5FA1\u6642\u306B\u304B\u3001\u5973\u5FA1\u3001\u66F4\u8863\u3042\u307E\u305F\u3055\u3076\u3089\u3072\u305F\u307E\u3072\u3051\u308B\u306A\u304B\u306B\u3001\u3044\u3068\u3084\u3080\u3054\u3068\u306A\u304D\u969B\u306B\u306F\u3042\u3089\u306C\u304C\u3001\u3059\u3050\u308C\u3066\u6642\u3081\u304D\u305F\u307E\u3075\u3042\u308A\u3051\u308A\u3002",
			 "biggenre":1,"genre":102,"gensaku":"","keyword":"\u7570\u4E16\u754C\u8EE2\u79FB \u6850\u58FA",
			 "general_firstup":"2018-02-19 19:22:14","general_lastup":"2019-08-16 08:35:42","novel_type":1,"end":1,
			 "general_all_no":359,"length":361837,"time":724,"isstop":0,"isr15":0,"isbl":0,"isgl":0,"iszankoku":0,"istensei":0,"istenni":1,
			 "pc_or_k":2,"global_point":132,"fav_novel_cnt":49,"review_cnt":0,"all_point":34,"all_hyoka_cnt":5,"sasie_cnt":0,
			 "kaiwaritu":30,"novelupdated_at":"2019-08-16 08:35:42","updated_at":"2019-08-16 08:38:16"
			}]`)},
			&SearchResult{
				AllCount: 12345,
				NovelInfos: []NovelInfo{
					{
						Title:          strp("徒然草"),
						NCode:          strp("N0000AA"),
						UserID:         strp("1234567"),
						Writer:         strp("吉田兼好"),
						Story:          strp("つれづれなるままにひぐらし硯に向かいて"),
						BigGenre:       intp(99),
						Genre:          intp(9903),
						Keywords:       []string{"R15", "年の差", "冒険"},
						GeneralFirstUp: jstDate(2019, 5, 6, 18, 39, 5, 0),
						GeneralLastUp:  jstDate(2019, 8, 16, 8, 35, 52, 0),
						NovelType:      intp(1),
						End:            intp(1),
						GeneralAllNo:   intp(45),
						Length:         intp(130245),
						Time:           intp(261),
						IsStop:         boolp(false),
						IsR15:          boolp(true),
						IsBoysLove:     boolp(false),
						IsGirlsLove:    boolp(false),
						IsZankoku:      boolp(true),
						IsTensei:       boolp(false),
						IsTenni:        boolp(false),
						PCOrK:          intp(2),
						GlobalPoint:    intp(154),
						FavNovelCount:  intp(33),
						ReviewCount:    intp(0),
						AllPoint:       intp(88),
						AllHyokaCount:  intp(10),
						SasieCount:     intp(0),
						KaiwaRitu:      intp(43),
						NovelUpdatedAt: jstDate(2019, 8, 16, 8, 35, 52, 0),
						UpdatedAt:      jstDate(2019, 8, 16, 8, 38, 16, 0),
					},
					{
						Title:          strp("源氏物語"),
						NCode:          strp("N9999ZZ"),
						UserID:         strp("9876543"),
						Writer:         strp("purple"),
						Story:          strp("いづれの御時にか、女御、更衣あまたさぶらひたまひけるなかに、いとやむごとなき際にはあらぬが、すぐれて時めきたまふありけり。"),
						BigGenre:       intp(1),
						Genre:          intp(102),
						Keywords:       []string{"異世界転移", "桐壺"},
						GeneralFirstUp: jstDate(2018, 2, 19, 19, 22, 14, 0),
						GeneralLastUp:  jstDate(2019, 8, 16, 8, 35, 42, 0),
						NovelType:      intp(1),
						End:            intp(1),
						GeneralAllNo:   intp(359),
						Length:         intp(361837),
						Time:           intp(724),
						IsStop:         boolp(false),
						IsR15:          boolp(false),
						IsBoysLove:     boolp(false),
						IsGirlsLove:    boolp(false),
						IsZankoku:      boolp(false),
						IsTensei:       boolp(false),
						IsTenni:        boolp(true),
						PCOrK:          intp(2),
						GlobalPoint:    intp(132),
						FavNovelCount:  intp(49),
						ReviewCount:    intp(0),
						AllPoint:       intp(34),
						AllHyokaCount:  intp(5),
						SasieCount:     intp(0),
						KaiwaRitu:      intp(30),
						NovelUpdatedAt: jstDate(2019, 8, 16, 8, 35, 42, 0),
						UpdatedAt:      jstDate(2019, 8, 16, 8, 38, 16, 0),
					},
				},
			},
			false,
		},
		{"title only",
			args{[]byte(`[{"allcount": 123},{"title":"AAA"}, {"title":"BBB"}]`)},
			&SearchResult{AllCount: 123, NovelInfos: []NovelInfo{{Title: strp("AAA")}, {Title: strp("BBB")}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseResponse(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				gstr, err := json.MarshalIndent(got, " ", "  ")
				if err != nil {
					gstr = []byte(fmt.Sprintf("%+v", got))
				}
				wstr, err := json.MarshalIndent(tt.want, " ", "  ")
				if err != nil {
					wstr = []byte(fmt.Sprintf("%+v", tt.want))
				}
				t.Errorf("parseResponse() = %s, want %s", gstr, wstr)
				// t.Errorf("parseResponse() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func jstDate(year int, month time.Month, day, hour, min, sec, nsec int) *time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil
	}
	t := time.Date(year, month, day, hour, min, sec, nsec, jst)
	return &t
}

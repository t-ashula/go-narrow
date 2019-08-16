package narrow

import "time"

// SearchResult contains fetch result
type SearchResult struct {
	// AllCount is number of novels that search parameter matched, not `len(NovelInfos)`
	AllCount int
	// NovelInfos are contains search result
	NovelInfos []NovelInfo
}

// NovelInfo contains a novel info
type NovelInfo struct {
	// 小説名
	Title *string
	// Nコード
	NCode *string
	// ユーザーID (R18 では取得できないので*)
	UserID *string
	// 作者名
	Writer *string
	// あらすじ
	Story *string
	// 大ジャンル
	BigGenre *int
	// ジャンル
	Genre *int
	//// 原作(未使用項目)
	//// Gensaku string
	// キーワード
	Keywords []string
	// 初回掲載日
	GeneralFirstUp *time.Time
	// 最終掲載日
	GeneralLastUp *time.Time
	// 連載か短編か
	NovelType *int
	// 連載中か
	End *int

	// 全掲載部分数
	GeneralAllNo *int
	// 文字数
	Length *int
	// 読了想定時間(文字数/500)
	Time *int
	// 長期連載停止中
	IsStop *bool
	// 必須キーワードに「R15」を含む(R18 では取得できないので*)
	IsR15 *bool
	// 必須キーワードにボーイズラブを含む
	IsBoysLove *bool
	// 必須キーワードにガールズラブを含む
	IsGirlsLove *bool
	// 必須キーワードに残酷な描写ありを含む
	IsZankoku *bool
	// 必須キーワードに異世界転生を含む
	IsTensei *bool
	// 必須キーワードに異世界転移を含む
	IsTenni *bool
	// pc_or_k
	PCOrK *int

	// 総合評価ポイント
	GlobalPoint *int
	// (R18)ブックマーク数
	FavNovelCount *int
	// レビュー数
	ReviewCount *int
	// 評価点
	AllPoint *int
	// 評価者数
	AllHyokaCount *int
	// 挿絵の数
	SasieCount *int
	// 会話率(%)
	KaiwaRitu *int
	// 小説の更新日
	NovelUpdatedAt *time.Time
	// システム用最終更新日時
	UpdatedAt *time.Time

	// R18 側結果かどうか
	fromSearchX bool
}

// OutputField represents 'of' parameter
type OutputField int

//
const (
	OutputFieldAll OutputField = iota
	OutputFieldTitle
	OutputFieldNCode
	OutputFieldUserID
	OutputFieldWriter
	OutputFieldStory
	OutputFieldBigGenre
	OutputFieldGenre
	OutputFieldKeyword
	OutputFieldGeneralFirstUp
	OutputFieldGeneralLastUp
	OutputFieldNovelType
	OutputFieldEnd
	OutputFieldGeneralAllNo
	OutputFieldLength
	OutputFieldTime
	OutputFieldIsStop
	OutputFieldIsR15
	OutputFieldIsBL
	OutputFieldIsGL
	OutputFieldIsZankoku
	OutputFieldIsTensei
	OutputFieldIsTenni
	OutputFieldPcOrK
	OutputFieldGlobalPoint
	OutputFieldFavNovelCount
	OutputFieldReviewCount
	OutputFieldAllPoint
	OutputFieldAllHyokaCount
	OutputFieldSasieCount
	OutputFieldKaiwaritu
	OutputFieldNovelUpdatedAt
	OutputFieldUpdatedAt
)

// SearchField stands for search word field
type SearchField int

// search fields
const (
	SearchFieldAll SearchField = iota
	SearchFieldTitle
	SearchFieldStory
	SearchFieldKeyword
	SearchFieldWriter
)

// BigGenre stands for large category
type BigGenre int

// big genres
const (
	BigGenreAll      BigGenre = 0
	BigGenreRenai    BigGenre = 1
	BigGenreFantasy  BigGenre = 2
	BigGenreBungei   BigGenre = 3
	BigGenreSF       BigGenre = 4
	BigGenreOther    BigGenre = 99
	BigGenreNonGenre BigGenre = 98
)

// Genre is genre
type Genre int

// genres
const (
	GenreAll Genre = 0

	GenreRenaiIsekai   Genre = 101
	GenreRenaiGenjitsu Genre = 102

	GenreFantasyHighFantasy Genre = 201
	GenreFantasyLowFantasy  Genre = 202

	GenreBungeiJunbungaku Genre = 301
	GenreBungeiHumanDrama Genre = 302
	GenreBungeiHistory    Genre = 303
	GenreBungeiMistrey    Genre = 304
	GenreBungeiHorror     Genre = 305
	GenreBungeiAction     Genre = 306
	GenreBungeiComedy     Genre = 307

	GenreSFVRGame Genre = 401
	GenreSFSpace  Genre = 402
	GenreSFSF     Genre = 403
	GenreSFPanic  Genre = 404

	GenreOtherFairyTale Genre = 9901
	GenreOtherPoetry    Genre = 9902
	GenreOtherEssei     Genre = 9903
	GenreOtherReplay    Genre = 9904
	GenreOtherOther     Genre = 9999

	GenreNongenreNongenre Genre = 9801
)

// NovelState represent status and type of novel
type NovelState int

const (
	// NovelStateAll is default status
	NovelStateAll NovelState = iota
	// NovelStateShortStory means t：短編
	NovelStateShortStory
	// NovelStateRensaiRunning means r：連載中
	NovelStateRensaiRunning
	// NovelStateRensaiEnded means er：完結済連載小説
	NovelStateRensaiEnded
	// NovelStateRensaiAll means re：すべての連載小説(連載中および完結済)
	NovelStateRensaiAll
	// NovelStateShortAndRensaiEnded means ter：短編と完結済連載小説
	NovelStateShortAndRensaiEnded
)

// Buntai stands indent style, empty lines rate. experimental
type Buntai int

// buntais
const (
	// 文体指定なし
	BuntaiAll Buntai = 0
	// 字下げされておらず、連続改行が多い作品
	BuntaiNoIndentManyEmptyLines Buntai = 1
	// 字下げされていないが、改行数は平均な作品
	BuntaiNoIndentAveraegEmpytLines Buntai = 2
	// 字下げが適切だが、連続改行が多い作品
	BuntaiIndentManyEmptyLines Buntai = 4
	// 字下げが適切でかつ改行数も平均な作品
	BuntaiIndentAverageEmptyLines Buntai = 6
)

// StopState is 連載停止中作品に関する指定
type StopState int

const (
	// StopStateAll is 指定なし
	StopStateAll StopState = 0
	// StopStateExclude is 長期連載停止中を除きます
	StopStateExclude StopState = 1
	// StopStateOnly  is 長期連載停止中のみ取得します
	StopStateOnly StopState = 2
)

// PickupState is for ispickup param
type PickupState int

// pickup state
const (
	PickupStateNone PickupState = iota
	PickupStateNot
	PickupStatePickup
)

// LastUpType is 最終掲載日(general_lastup)で抽出
type LastUpType int

//
const (
	// LastUpTypeNone is default, no condition
	LastUpTypeNone LastUpType = iota
	// LastUpTypeTimeStamp means Unixtime stamp
	LastUpTypeTimeStamp
	LastUpTypeThisWeek
	LastUpTypeLastWeek
	LastUpTypeSevenDay
	LastUpTypeThisMonth
	LastUpTypeLastMonth
)

// OrderItem is `order` param
type OrderItem int

const (
	// OrderItemNew is 新着順
	OrderItemNew OrderItem = iota
	// OrderItemFavNovelCount is ブックマーク数の多い順
	OrderItemFavNovelCount
	// OrderItemReviewCount is レビュー数の多い順
	OrderItemReviewCount
	// OrderItemHyoka is 総合ポイントの高い順
	OrderItemHyoka
	// OrderItemHyokaAsc is 総合ポイントの低い順
	OrderItemHyokaAsc
	// OrderItemImpressionCount is 感想の多い順
	OrderItemImpressionCount
	// OrderItemHyokaCount is 評価者数の多い順
	OrderItemHyokaCount
	// OrderItemHyokaCountAsc is 評価者数の少ない順
	OrderItemHyokaCountAsc
	// OrderItemWeekly is 週間ユニークユーザの多い順
	OrderItemWeekly
	// OrderItemLengthDesc is 小説本文の文字数が多い順
	OrderItemLengthDesc
	// OrderItemLengthAsc is 小説本文の文字数が少ない順
	OrderItemLengthAsc
	// OrderItemNCodeDesc is 新着投稿順
	OrderItemNCodeDesc
	// OrderItemOld is 更新が古い順
	OrderItemOld
)

// SearchParams contains API search parameters
type SearchParams struct {
	outputFields         []OutputField
	words                []string
	notWords             []string
	searchFields         []SearchField
	bigGenres            []BigGenre
	notBigGenres         []BigGenre
	genres               []Genre
	notGenres            []Genre
	userIDs              []int
	requiredKeywordFlags map[requiredKeyword]bool
	lengths              minmaxPair
	kaiwaritus           minmaxPair
	sasies               minmaxPair
	readTimes            minmaxPair
	ncodes               []string
	state                NovelState
	buntais              []Buntai
	stopState            StopState
	pickupState          PickupState
	lastUp               LastUpType
	lastUpStart          time.Time
	lastUpEnd            time.Time

	limit  int
	offset int
	order  OrderItem
}

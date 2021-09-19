package route

import (
	"Bancho/userDB"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func modeParser(ms *string) (i int) {
	switch *ms {
	case "osu":
		i = 0
	case "taiko":
		i = 1
	case "fruits":
		i = 2
	case "mania":
		i = 3
	default:
		i = 0
	}
	return
}

func typeParser(ts *string) (s string){
	switch *ts {

	case "country":
		s = "PP"
	case "performance":
		s = "PP"
	case "score":
		s = "RANKED_SCORE"
	default:
		s = "PP"
	}
	return
}

func pageParser(s string) (ss string) {
	atoi, err := strconv.Atoi(s)
	if err != nil || atoi <= 0 {
		return "LIMIT 50"
	}
	return fmt.Sprintf("LIMIT %d,50", atoi*50)
}

// /rankings/{mode}/{type}
func Rankings(c echo.Context) (err error) {

	var req request
	if err = c.Bind(&req); err != nil || req.Mode == "" || req.Type == "" {
		_ = c.NoContent(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	req.Cursor.Id = c.Request().URL.Query().Get("cursor[_id]")
	req.Cursor.Score = c.Request().URL.Query().Get("cursor[_score]")
	req.Cursor.Page = c.Request().URL.Query().Get("cursor[page]")
	// USER_ID, RANKED_SCORE, PLAYCOUNT, TOTAL_SCORE, TOTAL_HITS, LEVEL, PLAYTIME, AVG_ACCURACY, PP, REPLAY_WATCHED_COUNT
	//TODO 나중에 다이나믹쿼리로 변경
	rows, err := userDB.Maria.Query(`
SELECT  A.GLOBAL_RANK,A.USER_ID, A.RANKED_SCORE, A.PLAYCOUNT, A.TOTAL_SCORE, A.TOTAL_HITS, A.LEVEL, A.PLAYTIME, A.AVG_ACCURACY, A.PP, A.REPLAY_WATCHED_COUNT,
B.SS, B.SSH, B.S, B.SH, B.A from 
	(SELECT 
    ROW_NUMBER() OVER(ORDER BY `+typeParser(&req.Type)+` DESC) AS GLOBAL_RANK,USER_ID, MODE, RELAX, RANKED_SCORE, PLAYCOUNT, TOTAL_SCORE, TOTAL_HITS, LEVEL, PLAYTIME, AVG_ACCURACY, PP, REPLAY_WATCHED_COUNT
    FROM BANCHO.USER_STATUS_MODE where USER_ID in (select id from Ainu.users where ban_datetime = 0) AND MODE = ? AND RELAX = ? ORDER BY `+
		typeParser(&req.Type)+` DESC `+pageParser(req.Cursor.Page)+`) A
LEFT JOIN 
    (SELECT 
    USER_ID, MODE, RELAX, SS, SSH, S, SH, A, B, C, F
    FROM BANCHO.USER_GRADE_COUNT) B on A.USER_ID = B.USER_ID AND A.MODE = B.MODE AND A.RELAX = B.RELAX;
`,

		modeParser(&req.Mode),req.Relax,
	)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	defer rows.Close()
	var resp RankingsStruct

	for rows.Next() {
		var r RankingStruct
		err = rows.Scan(&r.GlobalRank, &r.User.Id, &r.RankedScore, &r.PlayCount, &r.TotalScore, &r.TotalHits, &r.Level.Current, &r.PlayTime, &r.HitAccuracy, &r.Pp, &r.ReplaysWatchedByOthers,
			&r.GradeCounts.Ss,&r.GradeCounts.Ssh,&r.GradeCounts.S,&r.GradeCounts.Sh,&r.GradeCounts.A,
			) // TODO
		if err != nil {
			c.NoContent(http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		resp.Ranking = append(resp.Ranking, r)

	}

	return c.JSON(http.StatusOK,resp)
}

type request struct {
	Mode string `param:"mode"`
	Type string `param:"type"`

	Relax     bool   `query:"relax"`
	Country   int    `query:"country"`
	Cursor    struct{
		Id string
		Score string
		Page string
	} `query:"cursor"`
	Filter    string `query:"filter"`
	Spotlight string `query:"spotlight"`
	Variant   string `query:"variant"`
}

type RankingsStruct struct {
	Cursor struct {
		Page int `json:"page"`
	} `json:"cursor"`
	Ranking []RankingStruct `json:"ranking"`
	Total   int             `json:"total"`
}
type RankingStruct struct {
	Level struct {
		Current  *int `json:"current"`
		Progress *int `json:"progress"`
	} `json:"level"`
	GlobalRank             *int     `json:"global_rank"`
	Pp                     *float64 `json:"pp"`
	RankedScore            *int64   `json:"ranked_score"`
	HitAccuracy            *float64 `json:"hit_accuracy"`
	PlayCount              *int     `json:"play_count"`
	PlayTime               *int     `json:"play_time"`
	TotalScore             *int64   `json:"total_score"`
	TotalHits              *int     `json:"total_hits"`
	MaximumCombo           *int     `json:"maximum_combo"`
	ReplaysWatchedByOthers *int     `json:"replays_watched_by_others"`
	IsRanked               *bool    `json:"is_ranked"`
	GradeCounts            struct {
		Ss  *int `json:"ss"`
		Ssh *int `json:"ssh"`
		S   *int `json:"s"`
		Sh  *int `json:"sh"`
		A   *int `json:"a"`
	} `json:"grade_counts"`
	User struct {
		AvatarUrl     *string      `json:"avatar_url"`
		CountryCode   *string      `json:"country_code"`
		DefaultGroup  *string      `json:"default_group"`
		Id            *int         `json:"id"`
		IsActive      *bool        `json:"is_active"`
		IsBot         *bool        `json:"is_bot"`
		IsDeleted     *bool        `json:"is_deleted"`
		IsOnline      *bool        `json:"is_online"`
		IsSupporter   *bool        `json:"is_supporter"`
		LastVisit     *time.Time  `json:"last_visit"`
		PmFriendsOnly *bool        `json:"pm_friends_only"`
		ProfileColour *interface{} `json:"profile_colour"`
		Username      *string      `json:"username"`
		Country       struct {
			Code *string `json:"code"`
			Name *string `json:"name"`
		} `json:"country"`
		Cover struct {
			CustomUrl *string `json:"custom_url"`
			Url       *string  `json:"url"`
			Id        *string `json:"id"`
		} `json:"cover"`
	} `json:"user"`
}

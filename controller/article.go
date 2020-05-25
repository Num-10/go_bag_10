package controller

import (
	"blog_go/model"
	"blog_go/util/e"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Index(c *gin.Context) {
	RequestColumn(85)
	e.Json(c, &e.Return{Code: e.SERVICE_SUCCESS})
}

type requestColumnData struct {
	Cid          int  `json:"cid"`
	WithGroupbuy bool `json:"with_groupbuy"`
}

func RequestColumn(cid int) {
	url := "https://time.geekbang.org/serv/v1/column/intro"

	data, _ := json.Marshal(requestColumnData{
		Cid:          cid,
		WithGroupbuy: true,
	})
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	reqest.Header.Add("Content-Type", "application/json")
	reqest.Header.Add("Host", "time.geekbang.org")
	reqest.Header.Add("Origin", "https://time.geekbang.org")
	reqest.Header.Add("User-Agent", " Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	myMap := make(map[string]interface{})
	json.Unmarshal(respBody, &myMap)
	fmt.Println("cid: ", cid, " code: ", myMap["code"])
	if myMap["code"].(float64) != 0 {
		return
	}
	switch myMap["data"].(type) {
	case []interface{}:
		return
	}
	val := myMap["data"].(map[string]interface{})
	column := &model.Columns{
		ColumnID:          int(val["id"].(float64)),
		ColumnTitle:       val["column_title"].(string),
		ColumnSubtitle:    val["column_subtitle"].(string),
		ColumnType:        int(val["column_type"].(float64)),
		ColumnPrice:       val["column_price"].(float64),
		ColumnPriceMarket: val["column_price_market"].(float64),
		ColumnBeginTime:   int(val["column_begin_time"].(float64)),
		ColumnEndTime:     int(val["column_end_time"].(float64)),
		ColumnSku:         int(val["column_sku"].(float64)),
		ColumnCoverInner:  val["column_cover_inner"].(string),
		ColumnCoverWxlite: val["column_cover_wxlite"].(string),
		AuthorName:        val["author_name"].(string),
		AuthorIntro:       val["author_intro"].(string),
		ArticleDoneCount:  int(val["article_count"].(float64)),
		ArticleTotalCount: int(val["article_total_count"].(float64)),
	}
	find := model.Columns{}
	find.Find(map[string]interface{}{"column_id": cid}, "")
	if find.ID > 0 && column.ArticleDoneCount == find.ArticleDoneCount {
		return
	}
	if find.ID > 0 {
		err = column.Update(map[string]interface{}{"id": find.ID}, column)
	} else {
		err = column.ColumnAdd()
	}
	if err != nil {
		fmt.Println(err)
	}
	RequestArticle(cid)
	return
}

type requestArticleData struct {
	Cid    int    `json:"cid"`
	Order  string `json:"order"`
	Prev   int    `json:"prev"`
	Sample bool   `json:"sample"`
	Size   int    `json:"size"`
}

func RequestArticle(cid int) {
	url := "https://time.geekbang.org/serv/v1/column/articles"

	data, _ := json.Marshal(requestArticleData{
		Cid:    cid,
		Order:  "earliest",
		Prev:   0,
		Sample: false,
		Size:   1000,
	})
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	reqest.Header.Add("Content-Type", "application/json")
	reqest.Header.Add("Host", "time.geekbang.org")
	reqest.Header.Add("Origin", "https://time.geekbang.org")
	reqest.Header.Add("User-Agent", " Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.129 Safari/537.36")

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		return
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	myMap := make(map[string]interface{})
	json.Unmarshal(respBody, &myMap)
	if myMap["code"].(float64) != 0 {
		return
	}
	if myMap["data"].(map[string]interface{})["page"].(map[string]interface{})["count"].(float64) == 0 {
		return
	}
	for _, val := range myMap["data"].(map[string]interface{})["list"].([]interface{}) {
		map_val := val.(map[string]interface{})
		article := &model.Articles{}
		article.ArticleID = int(map_val["id"].(float64))
		article.Find(map[string]interface{}{"article_id": article.ArticleID}, "")
		if article.ID > 0 {
			continue
		}
		article.ColumnID = cid
		article.ArticleTitle = map_val["article_title"].(string)
		article.ArticleSummary = map_val["article_summary"].(string)
		err = article.Create()
		if err != nil {
			fmt.Println("新增articles系列错误：", err)
			continue
		}
	}
}

type columnList struct {
	ColumnID          int           `json:"column_id"`
	ColumnTitle       string        `json:"column_title"`
	ColumnSubtitle    string        `json:"column_subtitle"`
	ColumnType        string        `json:"column_type"`
	ColumnPrice       float64       `json:"column_price"`
	ColumnPriceMarket float64       `json:"column_price_market"`
	ColumnBeginTime   string        `json:"column_begin_time"`
	ColumnEndTime     string        `json:"column_end_time"`
	ColumnSku         int           `json:"column_sku"`
	ColumnCoverInner  string        `json:"column_cover_inner"`
	ColumnCoverWxlite string        `json:"column_cover_wxlite"`
	AuthorName        string        `json:"author_name"`
	AuthorIntro       string        `json:"author_intro"`
	ArticleDoneCount  int           `json:"article_done_count"`
	ArticleTotalCount int           `json:"article_total_count"`
	ArticleList       []articleList `json:"article_list"`
}

type articleList struct {
	Link           string `json:"link"`
	ArticleTitle   string `json:"article_title"`
	ArticleSummary string `json:"article_summary"`
}

func List(c *gin.Context) {
	pageStr := c.Query("page")
	page := 1
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	pageSizeStr := c.Query("page_size")
	page_size := 1000
	if pageSizeStr != "" {
		page_size, _ = strconv.Atoi(pageSizeStr)
	}

	lists := []columnList{}
	list := columnList{}
	col := model.Columns{}
	cols := []model.Columns{}
	count := 0
	col.GetList(map[string]interface{}{}, map[string]interface{}{"page": page, "page_size": page_size, "order": "column_id asc"}, &cols, &count)
	for _, value := range cols {
		endTime := ""
		if value.ColumnEndTime != 0 {
			endTime = time.Unix(int64(value.ColumnEndTime), 0).Format("2006-01-02")
		}
		articles := []articleList{}
		article := articleList{}
		articleCol := model.Articles{}
		articleCols := []model.Articles{}
		articleCol.GetList(map[string]interface{}{"column_id": value.ColumnID}, map[string]interface{}{"order": "article_id asc"}, &articleCols, &count)
		for _, val := range articleCols {
			article = articleList{
				Link:           "https://time.geekbang.org/column/article/" + strconv.Itoa(val.ArticleID),
				ArticleTitle:   val.ArticleTitle,
				ArticleSummary: val.ArticleSummary,
			}
			articles = append(articles, article)
		}
		list = columnList{
			ColumnID:          value.ColumnID,
			ColumnTitle:       value.ColumnTitle,
			ColumnSubtitle:    value.ColumnSubtitle,
			ColumnType:        "value.ColumnType",
			ColumnPrice:       value.ColumnPrice / 100,
			ColumnPriceMarket: value.ColumnPriceMarket / 100,
			ColumnBeginTime:   time.Unix(int64(value.ColumnBeginTime), 0).Format("2006-01-02"),
			ColumnEndTime:     endTime,
			ColumnSku:         value.ColumnSku,
			ColumnCoverInner:  value.ColumnCoverInner,
			ColumnCoverWxlite: value.ColumnCoverWxlite,
			AuthorName:        value.AuthorName,
			AuthorIntro:       value.AuthorIntro,
			ArticleDoneCount:  value.ArticleDoneCount,
			ArticleTotalCount: value.ArticleTotalCount,
			ArticleList:       articles,
		}
		lists = append(lists, list)
	}

	e.Json(c, &e.Return{Code: e.SERVICE_SUCCESS, Data: lists})
}

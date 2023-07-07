package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Config struct {
	Token    string
	Modified string
	Ver      string
}

type Channel struct {
	ID         int
	Post_title string
}

type ChannelListResponse struct {
	Status string
	Data   struct {
		Posts []Channel
	}
}

type ChannelUrl struct {
	Type     string
	Url      string
	Priority string
}

type ChannelResponse struct {
	Data struct {
		Streams []ChannelUrl
	}
}

var config Config

const IndexPageUrl = "https://tingfm.com/region/cnr"
const PlayListUrl = "https://tingfm.com/wp-json/query/wnd_posts?is_main_query=true&type=radio&without_content=true&paged=1&update_post_term_cache=false&update_post_meta_cache=false&posts_per_page=30&_term_region=3"
const StreamUrl = "https://tingfm.com/wp-json/query/wndt_streams?post_id=4&in_web=true"

/*
* Parse Stream-Token from html
* Init config
 */
func initConfig() {
	resp, _ := http.Get(IndexPageUrl)
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	regex := regexp.MustCompile(`var wndt = (.*);`)
	ss := regex.FindStringSubmatch(bodyStr)
	json.Unmarshal([]byte(ss[1]), &config)
}

func main() {
	initConfig()
	url := fetchChannelUrl()
	fmt.Println(url)
}

/*
* Fetch playlist
 */
func fetchPlayList() {
	resp, err := http.Get(PlayListUrl)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var resp1 ChannelListResponse
	json.Unmarshal(body, &resp1)
	fmt.Println(resp1)
}

/*
* Fetch stream url
 */
func fetchChannelUrl() string {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", StreamUrl, nil)
	reqest.Header.Add("Stream-Token", config.Token)
	if err != nil {
		panic(err)
	}

	response, err := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	var channelResp ChannelResponse
	if nil != json.Unmarshal(body, &channelResp) {
		panic("JSON Unmarshal error")
	}
	return channelResp.Data.Streams[0].Url
}

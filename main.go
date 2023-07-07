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

func initConfig() {
	resp, _ := http.Get("https://tingfm.com/region/cnr")
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

func fetchPlayList() {
	url := "https://tingfm.com/wp-json/query/wnd_posts?is_main_query=true&type=radio&without_content=true&paged=1&update_post_term_cache=false&update_post_meta_cache=false&posts_per_page=30&_term_region=3"
	resp, _ := http.Get(url)
	body, _ := io.ReadAll(resp.Body)
	var resp1 ChannelListResponse
	json.Unmarshal(body, &resp1)
	fmt.Println(resp1)
}

func fetchChannelUrl() string {
	url := "https://tingfm.com/wp-json/query/wndt_streams?post_id=4&in_web=true"
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Stream-Token", config.Token)
	if err != nil {
		panic(err)
	}
	response, _ := client.Do(reqest)
	defer response.Body.Close()
	aa, _ := io.ReadAll(response.Body)
	var channelResp ChannelResponse
	json.Unmarshal(aa, &channelResp)
	return channelResp.Data.Streams[0].Url
}

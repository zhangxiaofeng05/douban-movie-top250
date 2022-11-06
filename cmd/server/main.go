package main

import (
	"github.com/zhangxiaofeng05/douban-movie-top250/parser"
	"github.com/zhangxiaofeng05/douban-movie-top250/utils"
	"time"
)

func Start() {
	var movies []parser.Movie

	pages := parser.GetPages(utils.BaseUrl)
	for _, page := range pages {
		list := parser.ParseMovie(page)
		movies = append(movies, list...)

		// mock browser
		time.Sleep(2 * time.Second)
	}
}

func main() {
	Start()
}

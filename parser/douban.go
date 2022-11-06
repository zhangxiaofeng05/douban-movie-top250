package parser

import (
	"fmt"
	"github.com/zhangxiaofeng05/douban-movie-top250/utils"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
)

type Movie struct {
	Title    string
	Subtitle string
	Other    string
	Desc     string
	Year     string
	Area     string
	Tag      string
	Star     string
	Comment  string
	Quote    string
}

type Page struct {
	Page int
	Url  string
}

func GetPages(url string) []Page {
	log.Println("start get all page info")
	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("User-Agent", browser.Random())

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("get url:%s error", url)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return parsePage(doc)
}

func parsePage(doc *goquery.Document) []Page {
	// current page(first page)
	pages := []Page{
		{
			Page: 1,
			Url:  "",
		},
	}

	doc.Find("#content > div > div.article > div.paginator > a").Each(func(i int, s *goquery.Selection) {
		page, _ := strconv.Atoi(s.Text())
		url, _ := s.Attr("href")
		pages = append(pages, Page{
			Page: page,
			Url:  url,
		})
	})
	log.Printf("pages num is %v", len(pages))
	return pages
}

func ParseMovie(page Page) []Movie {
	movies := make([]Movie, 0)
	fullUrl := fmt.Sprintf("%s%s", utils.BaseUrl, page.Url)
	log.Printf("parser page index:%d url:%s", page.Page, fullUrl)

	client := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, fullUrl, nil)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("User-Agent", browser.Random())

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		log.Fatalf("get url:%s error", fullUrl)
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".hd a span").Eq(0).Text()

		subTitle := s.Find(".hd a span").Eq(1).Text()
		subTitle = strings.TrimLeft(subTitle, " +/+")

		other := s.Find(".hd a span").Eq(2).Text()
		other = strings.TrimLeft(other, " +/+")

		desc := strings.TrimSpace(s.Find(".bd p").Eq(0).Text())
		descInfo := strings.Split(desc, "\n")
		desc = descInfo[0]

		movieDesc := strings.Split(descInfo[1], "/")
		year := movieDesc[0]
		area := movieDesc[1]
		tag := movieDesc[2]

		star := s.Find(".bd .star .rating_num").Text()

		comment := strings.TrimSpace(s.Find(".bd .star span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")

		quote := s.Find(".quote .inq").Text()

		movie := Movie{
			Title:    title,
			Subtitle: subTitle,
			Other:    other,
			Desc:     desc,
			Year:     year,
			Area:     area,
			Tag:      tag,
			Star:     star,
			Comment:  comment,
			Quote:    quote,
		}

		log.Printf("i: %d, movie:%+v", i, movie)
		movies = append(movies, movie)
	})
	return movies
}

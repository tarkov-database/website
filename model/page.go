package model

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/tarkov-database/website/model/api"
	"github.com/tarkov-database/website/model/item"
	"github.com/tarkov-database/website/model/location"
	"github.com/tarkov-database/website/version"

	"github.com/google/logger"
)

const (
	entityTypeItem     = "item"
	entityTypeLocation = "location"
)

var host string

func init() {
	if env := os.Getenv("HOST"); len(env) > 0 {
		host = env
	} else {
		logger.Warning("Host is not set!")
	}
}

type Page struct {
	App  *version.Application
	API  *api.API
	Host string
	Path string
	URI  string
}

func CreatePage(u *url.URL) *Page {
	return &Page{
		App:  version.App,
		Host: host,
		Path: u.Path,
		URI:  u.RequestURI(),
	}
}

func CreatePageWithAPI(u *url.URL) (*Page, error) {
	p := CreatePage(u)

	var err error

	p.API, err = api.GetAPI()
	if err != nil {
		return p, err
	}

	return p, nil
}

type IndexPage struct {
	*Page
}

func (p *Page) GetIndex() *IndexPage {
	return &IndexPage{p}
}

type EntityPage struct {
	*Page
}

type ItemPage struct {
	*EntityPage
	Item item.Entity
}

func (p *Page) Item(e item.Entity) *ItemPage {
	return &ItemPage{EntityPage: &EntityPage{p}, Item: e}
}

type LocationPage struct {
	*EntityPage
	Location *location.Location
}

func (p *Page) Location(loc *location.Location) *LocationPage {
	return &LocationPage{EntityPage: &EntityPage{p}, Location: loc}
}

type EntityList struct {
	*Page
	Type       string
	IsSearch   bool
	Keyword    string
	TotalCount int64
	PageCount  int64
	PageNumber int64
	PageNext   *Pagination
	PagePrev   *Pagination
}

type Pagination struct {
	Number int64
	URL    string
}

const itemLimit = 100

func (l *EntityList) GetPagination() {
	if l.TotalCount > itemLimit && l.URI != "" {
		u, err := url.Parse(l.URI)
		if err != nil {
			logger.Error(err)
			return
		}
		query := u.Query()

		if len(query.Get("p")) == 0 {
			query.Set("p", "")
		}
		page := &query["p"][0]

		var p int64 = 1
		if len(*page) != 0 {
			p, err = strconv.ParseInt(*page, 10, 0)
			if err != nil {
				logger.Error(err)
				return
			}
		}
		if p < 1 {
			p = 1
		}

		total := l.TotalCount / itemLimit
		if (l.TotalCount % itemLimit) != 0 {
			total = total + 1
		}

		var next int64
		if total > p {
			next = p + 1
		}

		var prev int64
		if p > 1 {
			prev = p - 1
		}

		l.PageNumber = p

		*page = strconv.FormatInt(next, 10)
		l.PageNext = &Pagination{
			Number: next,
			URL:    fmt.Sprintf("%s?%s", u.Path, query.Encode()),
		}

		*page = strconv.FormatInt(prev, 10)
		l.PagePrev = &Pagination{
			Number: prev,
			URL:    fmt.Sprintf("%s?%s", u.Path, query.Encode()),
		}
	}
}

type ItemList struct {
	*EntityList
	List []item.Entity
}

func (p *Page) ItemResult(res item.EntityResult, kw string, search bool) *ItemList {
	l := &ItemList{
		EntityList: &EntityList{
			Type:       entityTypeItem,
			Page:       p,
			IsSearch:   search,
			Keyword:    kw,
			TotalCount: res.GetCount(),
			PageCount:  int64(len(res.GetEntities())),
		},
		List: res.GetEntities(),
	}

	l.GetPagination()

	return l
}

type LocationList struct {
	*EntityList
	List []location.Location
}

func (p *Page) LocationResult(res *location.LocationResult, kw string, search bool) *LocationList {
	l := &LocationList{
		EntityList: &EntityList{
			Type:       entityTypeLocation,
			Page:       p,
			IsSearch:   search,
			Keyword:    kw,
			TotalCount: res.Count,
			PageCount:  int64(len(res.Items)),
		},
		List: res.Items,
	}

	l.GetPagination()

	return l
}

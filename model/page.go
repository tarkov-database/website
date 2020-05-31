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
	Entity interface{}
}

func (p *Page) Entity(e interface{}) *EntityPage {
	return &EntityPage{Page: p, Entity: e}
}

type EntityList struct {
	*Page
	Type        EntityType
	IsSearch    bool
	Keyword     string
	TotalCount  int64
	PageCount   int64
	PageTotal   int64
	PageCurrent int64
	PageNext    *Pagination
	PagePrev    *Pagination
	Filter      Filter
	List        interface{}
}

type Pagination struct {
	Number int64
	URL    string
}

const (
	itemLimit = 100
	pageKey   = "p"
)

func (l *EntityList) GetPagination() error {
	if l.TotalCount > itemLimit && l.URI != "" {
		u, err := url.ParseRequestURI(l.URI)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidInput, err)
		}

		query := u.Query()

		var p int64 = 1
		if v := query.Get(pageKey); v != "" && v != "1" {
			p, err = strconv.ParseInt(v, 10, 0)
			if err != nil {
				return fmt.Errorf("%w: %s", ErrInvalidInput, err)
			}
		}

		if p < 1 {
			p = 1
		}

		l.PageTotal = l.TotalCount / itemLimit
		if (l.TotalCount % itemLimit) != 0 {
			l.PageTotal = l.PageTotal + 1
		}

		var next int64
		if l.PageTotal > p {
			next = p + 1
		}

		var prev int64
		if p > 1 {
			prev = p - 1
		}

		l.PageCurrent = p

		query.Set(pageKey, strconv.FormatInt(next, 10))
		l.PageNext = &Pagination{
			Number: next,
			URL:    u.Path + "?" + query.Encode(),
		}

		query.Set(pageKey, strconv.FormatInt(prev, 10))
		l.PagePrev = &Pagination{
			Number: prev,
			URL:    u.Path + "?" + query.Encode(),
		}
	}

	return nil
}

func (p *Page) Result(res interface{}, kw string) (*EntityList, error) {
	el := &EntityList{
		Page:    p,
		Keyword: kw,
	}

	switch v := res.(type) {
	case item.EntityResult:
		el.Type = TypeItem
		el.TotalCount = v.GetCount()
		list := v.GetEntities()
		el.List = list
		el.PageCount = int64(len(list))
		el.Filter = v.GetKind().GetFilter()
	case *location.LocationResult:
		el.Type = TypeLocation
		el.TotalCount = v.Count
		el.List = v.Items
		el.PageCount = int64(len(v.Items))
		el.Filter = location.GetFilter()
	case []*SearchResult:
		el.IsSearch = true
		el.TotalCount = int64(len(v))
		el.List = v
		el.PageCount = el.TotalCount
	}

	if err := el.GetPagination(); err != nil {
		return nil, err
	}

	return el, nil
}

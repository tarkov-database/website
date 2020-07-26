package item

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/tarkov-database/website/core/api"
	"github.com/tarkov-database/website/core/search"

	"github.com/google/logger"
)

const defaultSort = "name"

func GetItem(id objectID, kind Kind) (Entity, error) {
	entity, err := kind.GetEntity()
	if err != nil {
		return entity, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = api.GET(ctx, fmt.Sprintf("/item/%s/%s", kind, id), &api.Options{}, entity); err != nil {
		return entity, err
	}

	return entity, nil
}

func GetItems(kind Kind, opts *api.Options) (EntityResult, error) {
	result, err := kind.GetEntityResult()
	if err != nil {
		return result, err
	}

	if opts.Sort == "" {
		opts.Sort = defaultSort
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := api.GET(ctx, fmt.Sprintf("/item/%s", kind), opts, result); err != nil {
		return result, err
	}

	return result, nil
}

func GetItemsByID(ids string, kind Kind, opts *api.Options) (EntityResult, error) {
	if _, ok := opts.Filter["id"]; !ok {
		opts.Filter = make(map[string]string)
	}

	opts.Filter["id"] = ids

	result, err := GetItems(kind, opts)
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetItemsBySearch(text string, limit int) (EntityResult, error) {
	opts := &api.Options{
		Sort:  "",
		Limit: limit,
		Filter: map[string]string{
			"search": text,
		},
	}

	result, err := KindCommon.GetEntityResult()
	if err != nil {
		return result, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err = api.GET(ctx, "/item", opts, result); err != nil {
		return result, err
	}

	return result, nil
}

func GetItemList(items ItemList, sort string) map[Kind][]Entity {
	ch := make(chan EntityResult)
	wg := &sync.WaitGroup{}

	for k, v := range items {
		chunks := toQueryArray(v)
		wg.Add(len(chunks))

		for _, c := range chunks {
			go func(k Kind, ids string) {
				defer wg.Done()

				res, err := GetItemsByID(ids, k, &api.Options{Limit: 100, Sort: sort})
				if err != nil {
					logger.Error(err)
					return
				}

				ch <- res
			}(k, c)
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := make(map[Kind][]Entity)
	for r := range ch {
		items := r.GetEntities()
		kind := items[0].GetKind()
		result[kind] = append(result[kind], items...)
	}

	return result
}

const arrayLimit = 100

func toQueryArray(arr []objectID) []string {
	length := len(arr)
	count := length / arrayLimit

	if count*arrayLimit != length {
		count++
	}

	qa := make([]string, count)
	if full := count - 1; full > 0 {
		for i := 0; i < full; i++ {
			s := i * arrayLimit
			e := s + arrayLimit
			qa[i] = strings.Join(arr[s:e], ",")
		}
	}

	s := (count - 1) * arrayLimit
	qa[count-1] = strings.Join(arr[s:], ",")

	return qa
}

func Search(term string, limit int, kind *Kind) (*search.Result, error) {
	if kind != nil {
		term = fmt.Sprintf("kind:%s AND %s", kind, term)
	}

	query := &search.Query{
		Term:  term,
		Index: search.IndexItem,
	}

	opts := &search.Options{
		Limit: limit,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := search.Search(ctx, query, opts)
	if err != nil {
		return result, err
	}

	return result, nil
}

func SearchByName(term string, limit int, kind *Kind) (*search.Result, error) {
	if kind != nil {
		term = fmt.Sprintf("kind:%s AND name:%s", kind, term)
	} else {
		term = fmt.Sprintf("name:%s", term)
	}

	query := &search.Query{
		Term:  term,
		Index: search.IndexItem,
	}

	opts := &search.Options{
		Limit: limit,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := search.Search(ctx, query, opts)
	if err != nil {
		return result, err
	}

	return result, nil
}

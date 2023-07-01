package view

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/tarkov-database/website/model/item"

	"github.com/goccy/go-json"
	"github.com/google/logger"
)

var errTemplateNotExist = errors.New("template doesn't exist")

const templateDir = "view/templates"

var (
	templates map[string]*template.Template
	funcMap   = template.FuncMap{
		// Helpers
		"addFloat":         addFloat,
		"subtractFloat":    subtractFloat,
		"multiplyFloat":    multiplyFloat,
		"divideFloat":      divideFloat,
		"decimalToPercent": decimalToPercent,
		"toTitle":          strings.Title,
		"camelToTitle":     camelToTitle,
		"formatTime":       formatTime,
		"localeString":     localeString,
		"setQuery":         setQuery,
		"hasQuery":         hasQuery,
		"hasPrefix":        hasPrefix,
		"queryEscape":      url.QueryEscape,
		"toBase64":         toBase64,

		"appendStaticHash": appendStaticHash,
		"categoryToName":   item.CategoryToDisplayName,
		"kindToCategory":   item.KindToCategory,

		// Get entities
		"getItem":                getItem,
		"getAmmunitionByCaliber": getAmmunitionByCaliber,
		"getAmmunitionRangeData": getAmmunitionRangeData,
		"resolveItemList":        resolveItemList,
		"resolveSlots":           resolveSlots,
	}
)

func init() {
	f, err := os.ReadFile(templateDir + "/index.json")
	if err != nil {
		logger.Fatal(err)
	}

	index := make(map[string][]string)
	if err = json.Unmarshal(f, &index); err != nil {
		logger.Fatal(err)
	}

	templates = make(map[string]*template.Template)

	ch := make(chan *template.Template, 1)
	wg := &sync.WaitGroup{}
	wg.Add(len(index))

	for k, v := range index {
		go loadTemplate(k, v, ch, wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for t := range ch {
		templates[t.Name()] = t
	}
}

func loadTemplate(name string, files []string, ch chan *template.Template, wg *sync.WaitGroup) {
	for i, k := range files {
		files[i] = templateDir + k
	}

	ch <- template.Must(template.New(name).Funcs(funcMap).ParseFiles(files...))
	wg.Done()
}

func RenderHTML(t string, d interface{}, w http.ResponseWriter) {
	tmpl, ok := templates[t]
	if !ok {
		switch strings.Split(t, "_")[0] {
		case "item":
			tmpl = templates["item_common"]
		case "grid":
			tmpl = templates["grid_generic"]
		case "table":
			tmpl = templates["table_generic"]
		default:
			logger.Error(errTemplateNotExist)
			tmpl = templates["status_500"]
		}
	}

	if err := tmpl.ExecuteTemplate(w, "base", d); err != nil {
		logger.Error(err)
		http.Error(w, fmt.Sprint("500 Internal Server Error"), http.StatusInternalServerError)
	}
}

func RenderJSON(data interface{}, status int, w http.ResponseWriter) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		logger.Error(err)
	}
}

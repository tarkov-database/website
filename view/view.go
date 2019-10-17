package view

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/tarkov-database/website/model/item"

	"github.com/google/logger"
)

var (
	templates map[string]*template.Template

	errTemplateNotExist = errors.New("template doesn't exist")
)

func init() {
	templates = make(map[string]*template.Template)

	funcMap := template.FuncMap{
		// Helpers
		"addFloat":         addFloat,
		"subtractFloat":    subtractFloat,
		"multiplyFloat":    multiplyFloat,
		"toTitle":          strings.Title,
		"formatTime":       formatTime,
		"localeString":     localeString,
		"setQuery":         setQuery,
		"hasQuery":         hasQuery,
		"decimalToPercent": decimalToPercent,
		"staticHashShort":  staticHashShort,
		"categoryToName":   item.CategoryToDisplayName,
		"kindToCategory":   item.KindToCategory,

		// Get entities
		"getItem":                getItem,
		"getAmmunitionByCaliber": getAmmunitionByCaliber,
		"resolveItemList":        resolveItemList,
		"resolveSlots":           resolveSlots,
	}

	// General
	templates["index"] = template.Must(template.New("index").Funcs(funcMap).ParseFiles(
		"view/templates/index/index.gohtml", "view/templates/index/base.gohtml",
		"view/templates/index/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["list"] = template.Must(template.New("list").Funcs(funcMap).ParseFiles(
		"view/templates/list/list.gohtml", "view/templates/list/base.gohtml",
		"view/templates/list/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))

	// Item
	templates["item_common"] = template.Must(template.New("item_common").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_common.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_ammunition"] = template.Must(template.New("item_ammunition").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_ammunition.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_firearm"] = template.Must(template.New("item_firearm").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_firearm.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_melee"] = template.Must(template.New("item_melee").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_melee.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_magazine"] = template.Must(template.New("item_magazine").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_magazine.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_modification"] = template.Must(template.New("item_modification").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_modification.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_grenade"] = template.Must(template.New("item_grenade").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_grenade.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_armor"] = template.Must(template.New("item_armor").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_armor.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_backpack"] = template.Must(template.New("item_backpack").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_backpack.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_clothing"] = template.Must(template.New("item_clothing").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_clothing.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_tacticalrig"] = template.Must(template.New("item_tacticalrig").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_tacticalrig.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_medical"] = template.Must(template.New("item_medical").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_medical.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_food"] = template.Must(template.New("item_food").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_food.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_container"] = template.Must(template.New("item_container").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_container.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["item_headphone"] = template.Must(template.New("item_headphone").Funcs(funcMap).ParseFiles(
		"view/templates/item/item.gohtml", "view/templates/item/item_headphone.gohtml",
		"view/templates/item/base.gohtml", "view/templates/item/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))

	// Table
	templates["table_generic"] = template.Must(template.New("table_generic").Funcs(funcMap).ParseFiles(
		"view/templates/table/table_generic.gohtml", "view/templates/table/base.gohtml",
		"view/templates/table/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["table_ammunition"] = template.Must(template.New("table_ammunition").Funcs(funcMap).ParseFiles(
		"view/templates/table/table_ammunition.gohtml", "view/templates/table/base.gohtml",
		"view/templates/table/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["table_armor"] = template.Must(template.New("table_armor").Funcs(funcMap).ParseFiles(
		"view/templates/table/table_armor.gohtml", "view/templates/table/base.gohtml",
		"view/templates/table/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["table_magazine"] = template.Must(template.New("table_magazine").Funcs(funcMap).ParseFiles(
		"view/templates/table/table_magazine.gohtml", "view/templates/table/base.gohtml",
		"view/templates/table/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))

	// Status
	templates["status_404"] = template.Must(template.New("status_404").Funcs(funcMap).ParseFiles(
		"view/templates/status/404.gohtml", "view/templates/status/base.gohtml",
		"view/templates/status/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["status_500"] = template.Must(template.New("status_500").Funcs(funcMap).ParseFiles(
		"view/templates/status/500.gohtml", "view/templates/status/base.gohtml",
		"view/templates/status/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["status_503"] = template.Must(template.New("status_503").Funcs(funcMap).ParseFiles(
		"view/templates/status/503.gohtml", "view/templates/status/base.gohtml",
		"view/templates/status/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))

	templates["about"] = template.Must(template.New("about").Funcs(funcMap).ParseFiles(
		"view/templates/about/about.gohtml", "view/templates/about/base.gohtml",
		"view/templates/about/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
	templates["projects"] = template.Must(template.New("projects").Funcs(funcMap).ParseFiles(
		"view/templates/projects/projects.gohtml", "view/templates/projects/base.gohtml",
		"view/templates/projects/head.gohtml", "view/templates/header.gohtml", "view/templates/footer.gohtml"))
}

func Render(t string, d interface{}, w http.ResponseWriter) {
	tmpl, ok := templates[t]
	if !ok {
		switch strings.Split(t, "_")[0] {
		case "item":
			tmpl = templates["item_common"]
		case "table":
			tmpl = templates["table_generic"]
		default:
			logger.Error(errTemplateNotExist)
			tmpl = templates["status_500"]
			return
		}
	}

	if err := tmpl.ExecuteTemplate(w, "base", d); err != nil {
		logger.Error(err)
		http.Error(w, fmt.Sprint("500 Internal Server Error"), http.StatusInternalServerError)
	}
}

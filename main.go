package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/ledongthuc/goterators"
)

const BASE_URL = "https://sinete.com.br"

type Category struct {
	url string
}

type Product struct {
	sku         string
	ean         string
	name        string
	description string
	images      []string
	brand       string
	factory     string
}

func getCategories(doc soup.Root) []Category {
	allCategoriesElem := doc.FindAll("a", "class", "menu-link")
	allCategoriesLink := goterators.Map(allCategoriesElem, func(item soup.Root) Category {
		return Category{url: item.Attrs()["href"]}
	})
	return goterators.Filter(allCategoriesLink, func(item Category) bool {
		return item.url != ""
	})
}

func getProductsLinkFromCategory(category Category) []string {
	page := 1

	var allProductsLink []string

	for {
		client := &http.Client{}
		var data = strings.NewReader(`ajax_nav=1`)
		req, err := http.NewRequest("POST", category.url+"?p="+strconv.Itoa(page), data)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("accept", "*/*")
		req.Header.Set("accept-language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
		req.Header.Set("cookie", "PHPSESSID=a60d7a683a4f06038c1a0cd9546e577c; form_key=gIbYSYa3BgwEVmMd; form_key=gIbYSYa3BgwEVmMd; mage-messages=; cookies-policy=1; private_content_version=913083a594f90562bc29be3ff5648957; mage-cache-storage=%7B%7D; mage-cache-storage-section-invalidation=%7B%7D; mage-cache-sessid=true; section_data_ids=%7B%22cart%22%3A1714505382%7D; recently_viewed_product=%7B%7D; recently_viewed_product_previous=%7B%7D; recently_compared_product=%7B%7D; recently_compared_product_previous=%7B%7D; product_data_storage=%7B%7D")
		req.Header.Set("dnt", "1")
		req.Header.Set("origin", "https://sinete.com.br")
		req.Header.Set("referer", "https://sinete.com.br/medicamentos.html")
		req.Header.Set("sec-ch-ua", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Linux"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
		req.Header.Set("x-requested-with", "XMLHttpRequest")
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)

		pageJson := struct {
			CategoryProducts string `json:"category_products"`
			CatalogLeftnav   string `json:"catalog_leftnav"`
			PageMainTitle    string `json:"page_main_title"`
			UpdatedUrl       string `json:"updated_url"`
		}{}

		err = json.Unmarshal(bodyBytes, &pageJson)

		if err != nil {
			log.Fatal(err)
		}

		doc := soup.HTMLParse(pageJson.CategoryProducts)

		productsLink := goterators.Map(doc.FindAll("a", "class", "product-item-link"), func(item soup.Root) string {
			return item.Attrs()["href"]
		})

		fmt.Printf("%s | page: %d | products: %d\n", category.url, page, len(productsLink))

		allProductsLink = append(allProductsLink, productsLink...)

		page += 1
	}
}

func main() {
	resp, err := soup.Get(BASE_URL)

	if err != nil {
		log.Fatal(err)
	}

	doc := soup.HTMLParse(resp)

	allCategories := getCategories(doc)

	for _, category := range allCategories {
		allProductsLink := getProductsLinkFromCategory(category)
		fmt.Println(allProductsLink)
		break
	}
}

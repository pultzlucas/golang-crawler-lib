package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/ledongthuc/goterators"
)

const BASE_URL = "https://pokemondb.net/pokedex/national"

type Pokemon struct {
	id    string
	name  string
	types []string
	image string
}

func main() {
	resp, err := soup.Get(BASE_URL)

	if err != nil {
		os.Exit(1)
	}

	doc := soup.HTMLParse(resp)
	pokeCards := doc.FindAll("div", "class", "infocard")

	var allPokemons []Pokemon

	for _, card := range pokeCards {
		pokeId := card.Find("span", "class", "infocard-lg-data").Find("small").Text()
		pokeName := card.Find("a", "class", "ent-name").Text()
		pokeTypes := goterators.Map(card.FindAll("a", "class", "itype"), func(item soup.Root) string {
			return item.Text()
		})

		pokeImage := card.Find("img", "class", "img-sprite").Attrs()["src"]

		pokemon := Pokemon{
			id:    pokeId,
			name:  pokeName,
			types: pokeTypes,
			image: pokeImage,
		}

		allPokemons = append(allPokemons, pokemon)
	}

	f, err := os.Create("output_pokemons.csv")
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	defer csvWriter.Flush()

	for _, pokemon := range allPokemons {
		row := []string{pokemon.id, pokemon.name, pokemon.image, strings.Join(pokemon.types, "|")}
		if err := csvWriter.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}

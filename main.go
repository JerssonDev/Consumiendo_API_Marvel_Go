package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type marvelRes struct {
	Code            int    `json:"code"`
	Status          string `json:"status"`
	Copyright       string `json:"copyright"`
	AttributionText string `json:"attributionText"`
	AttributionHTML string `json:"attributionHTML"`
	Etag            string `json:"etag"`
	Data            data   `json:"data"`
}

type data struct {
	Offset  int       `json:"offset"`
	Limit   int       `json:"limit"`
	Total   int       `json:"total"`
	Count   int       `json:"count"`
	Results []results `json:"results"`
}

type results struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Modified    string    `json:"modified"`
	Thumbnail   thumbnail `json:"thumbnail"`
	ResourceURI string    `json:"resourceURI"`
	Comics      comics    `json:"comics"`
	Series      series    `json:"series"`
	Stories     stories   `json:"stories"`
	Events      events    `json:"events"`
	Urls        []urls    `json:"urls"`
}

type thumbnail struct {
	Path      string `json:"path"`
	Extension string `json:"extension"`
}

type comics struct {
	Available     int          `json:"available"`
	CollectionURI string       `json:"collectionURI"`
	ItemsComic    []itemsComic `json:"items"`
	Returned      int          `json:"returned"`
}

type itemsComic struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
}

type series struct {
	Available     int           `json:"available"`
	CollectionURI string        `json:"collectionURI"`
	ItemsSeries   []itemsSeries `json:"items"`
	Returned      int           `json:"returned"`
}

type itemsSeries struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
}

type stories struct {
	Available     int            `json:"available"`
	CollectionURI string         `json:"collectionURI"`
	ItemsStories  []itemsStories `json:"items"`
	Returned      int            `json:"returned"`
}

type itemsStories struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
	Type        string `json:"type"`
}

type events struct {
	Available     int           `json:"available"`
	CollectionURI string        `json:"collectionURI"`
	ItemsEvents   []itemsEvents `json:"items"`
	Returned      int           `json:"returned"`
}

type itemsEvents struct {
	ResourceURI string `json:"resourceURI"`
	Name        string `json:"name"`
}

type urls struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

func main() {

	fmt.Println("\n¡Bienvenido! ¡Vamos a consultar la API de Marvel!")
	hash := "37cbc8678064f8a9c30209cf100546f0"
	apiKey := "01fc08457a033330da59f035421d552c"
	url := fmt.Sprintf("https://gateway.marvel.com:443/v1/public/characters?ts=1&apikey=%s&hash=%s&limit=20", apiKey, hash)

	var idSelect int

	fmt.Println("\n *-*-* Menu *-*-*")
	fmt.Println("1.- Digite el nombre del SuperHeroe para mostrar información.")
	fmt.Println("2.- Listar los 20 primeros Registros")
	fmt.Println("")

	fmt.Print("Opcion: ")
	fmt.Scanln(&idSelect)

	switch idSelect {
	case 1:
		fmt.Print("\nIngrese el nombre del SuperHeroe: ")

		reader := bufio.NewReader(os.Stdin)            // crea una nueva instancia para input
		entrada, _ := reader.ReadString('\n')          // Leer hasta el separador de salto de línea
		eleccion := strings.TrimRight(entrada, "\r\n") // remueve los saltos de linea

		nombre := fmt.Sprintf("&nameStartsWith=%s", urlEncoded(eleccion))

		var buffer bytes.Buffer

		buffer.WriteString(url)
		buffer.WriteString(nombre)

		newUrl := buffer.String()

		//fmt.Println(nombre)
		//fmt.Println(newUrl)

		fmt.Println("\n... Mostrando Data ...")

		fmt.Println(getResponseOp1(newUrl))
		fmt.Println("")

	case 2:
		fmt.Println("\nMostrando los Primeros 20 Registros...")

		fmt.Println("\n... Mostrando Data ...")

		fmt.Println("")

		getResponseOp2(url)

		fmt.Println("")

	default:
		fmt.Println("Valor fuera de rango")
	}

}

func getStations(body []byte) (*marvelRes, error) {
	var s *marvelRes = new(marvelRes)
	err := json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return s, err
}

func urlEncoded(str string) string {

	t := &url.URL{Path: str}
	return t.String()

}

func getResponseOp1(str string) string {

	var res string

	resp, err := http.Get(str)

	if err != nil {
		fmt.Printf("La Solicitud HTTP ha fallado: %s\n", err)
	} else {

		data, _ := ioutil.ReadAll(resp.Body)
		s, err := getStations([]byte(data))
		if err != nil {
			res = err.Error()
		} else {
			var nombreComic, nombreSeries, nombreStories, nombreEvents string

			if s.Data.Total > 0 {
				id := s.Data.Results[0].ID
				nombre := s.Data.Results[0].Name
				descripcion := s.Data.Results[0].Description
				modificado := s.Data.Results[0].Modified
				if s.Data.Results[0].Comics.Available > 0 {
					nombreComic = s.Data.Results[0].Comics.ItemsComic[0].Name
				} else {
					nombreComic = "No Tiene"
				}
				if s.Data.Results[0].Series.Available > 0 {
					nombreSeries = s.Data.Results[0].Series.ItemsSeries[0].Name
				} else {
					nombreSeries = "No Tiene"
				}
				if s.Data.Results[0].Stories.Available > 0 {
					nombreStories = s.Data.Results[0].Stories.ItemsStories[0].Name
				} else {
					nombreStories = "No Tiene"
				}
				if s.Data.Results[0].Events.Available > 0 {
					nombreEvents = s.Data.Results[0].Events.ItemsEvents[0].Name
				} else {
					nombreEvents = "No Tiene"
				}

				url := s.Data.Results[0].Urls[0].URL
				res = fmt.Sprintf("\nCon su id %d Se encontro al SuperHeroe %s su descripción es la siguiente: %s \nLa fecha de modificacion es: %s \nTiene => Comic: %s, Serie: %s, Stories: %s, Events: %s; mas informacion en el siguiente enlace: %s", id, nombre, descripcion, modificado, nombreComic, nombreSeries, nombreStories, nombreEvents, url)

			} else {
				res = "No se pudieron encontrar datos del SuperHeroe buscado"
			}
		}
	}
	return res
}

func getResponseOp2(str string) {

	resp, err := http.Get(str)

	if err != nil {
		fmt.Printf("La Solicitud HTTP ha fallado: %s\n", err)
	} else {

		data, _ := ioutil.ReadAll(resp.Body)
		s, err := getStations([]byte(data))

		if err != nil {
			fmt.Println(err)
		} else {

			i := 0
			for i < s.Data.Count {
				fmt.Println(" - " + s.Data.Results[i].Name)
				i = i + 1
			}

		}
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

// код писать тут
const accessTocken  = "123b4j23k4hvjjgvk421g34k"

func SearchServer(limit int, offset int, query string, orderField string, orderBy int) ([]Row, string){
	filePath := "dataset.xml"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, `{"Error": "BadOrderField"}`
	}
	all := ReadClients(file)
	result := make([]Row, 0)

	for _, val := range all {

		if strings.Contains(val.FirstName, query) || strings.Contains(val.About, query) {
			result = append(result, val)
		}
	}

	if orderField == ""{
		orderField = "Name"
	}

	if orderField != "Name" && orderField != "Id" && orderField != "Age"{
		fmt.Println("error happen")
		return nil, `{"Error": "BadOrderField"}`
	}

	switch orderBy {
	case -1:
		sort.Slice(result, func(i, j int) bool {
			switch orderField {
			case "Name":
				return result[i].Name > result[j].Name
			case "Age":
				return result[i].Age > result[j].Age
			default:
				return result[i].Id > result[j].Id
			}
		})
	case 1:
		sort.Slice(result, func(i, j int) bool {
			switch orderField {
			case "Name":
				return result[i].Name < result[j].Name
			case "Age":
				return result[i].Age < result[j].Age
			default:
				return result[i].Id < result[j].Id
			}
		})
	case 0:

	default:
		return nil, `{"Error": "ErrorBadOrderField"}`
	}

	return result, ""
}

type Row struct {
	Id int `xml:"id"`
	Age  int `xml:"age"`
	FirstName string `xml:"first_name"`
	LastName string `xml:"last_name"`
	Gender string `xml:"gender"`
	About string `xml:"about"`
	Name string
}

func ReadClients(file []byte) []Row {

	input := bytes.NewReader(file)
	decoder := xml.NewDecoder(input)
	clients := make([]Row, 0, 10)


	for {
		tok, tokenErr := decoder.Token()
		if tokenErr != nil && tokenErr != io.EOF {
			fmt.Println("error happen", tokenErr)
			break
		} else if tokenErr == io.EOF {
			break
		}
		if tok == nil {
			fmt.Println("t is nil break")
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			//fmt.Println(tok.Name)
			if tok.Name.Local == "row"{
				client := Row{}
				err := decoder.DecodeElement(&client, &tok)
				if err != nil{
					fmt.Println("error happend", err)
				}
				client.Name = client.FirstName+client.LastName
				clients = append(clients, client)
			}
		}

	}

	return clients
}

func WrapJson(clients []Row) string{
	var jsonData = "["
	for _, val := range clients{
		user := User{
			val.Id,
			val.Name,
			val.Age,
			val.About,
			val.Gender,
		}
		data, _ := json.Marshal(user)
		jsonData += string(data)
	}
	jsonData+="]"
	return jsonData
}

func HandlerFunc(w http.ResponseWriter, r *http.Request){

	limit, err := strconv.Atoi(r.FormValue("limit"))
	offset, err:= strconv.Atoi(r.FormValue("offset"))
	query := r.FormValue("query")
	orderField := r.FormValue("order_field")
	orderBy, err := strconv.Atoi(r.FormValue("order_by"))

	if query[0] == '~'{
		time.Sleep(2*time.Second)
	}


	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if orderBy == 121 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	at := r.Header.Get("AccessToken")
	if at != accessTocken{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	clients, errr := SearchServer(limit, offset, query, orderField,  orderBy)

	jsonResp := WrapJson(clients)

	if errr != "" {
		if query[0] == '_'{
			errr+="cdscsd"
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errr))
		return
	}

	if query[0] == '_'{
		jsonResp+="cdscsd"
	}
	w.Write([]byte(jsonResp))
}



func TestFindUsers(t *testing.T){

	srs := []SearchRequest{
		{
			0,
			0,
			"Hilda",
			"Id",
			1,
		},
		{
			0,
			0,
			"Hilda",
			"Id",
			100,
		},
		{
			0,
			0,
			"Hilda",
			"Id",
			100,
		},
		{
			0,
			0,
			"Hilda",
			"Id",
			100,
		},
		{
			0,
			0,
			"Hilda",
			"Id",
			100,
		},
		{
			0,
			0,
			"Hilda",
			"Id",
			0,
		},
		{
			2,
			0,
			"Hilda",
			"Id",
			0,
		},
		{
			26,
			0,
			"Hilda",
			"Id",
			0,
		},
		{
			0,
			-1,
			"Hilda",
			"Id",
			0,
		},
		{
			-1,
			0,
			"Hilda",
			"Id",
			0,
		},
		{
			200,
			0,
			"Hilda",
			"Id",
			0,
		},
		{
			0,
			0,
			"",
			"",
			101,
		},
		{
			0,
			0,
			"\nn12",
			"OLOLO",
			0,
		},
		{
			0,
			0,
			"\nn12",
			"OLOLO",
			121,
		},
		{
			0,
			0,
			"",
			"OLOLO",
			0,
		},
		{
			0,
			0,
			"_nhb_",
			"Age",
			0,
		},
		{
			0,
			0,
			"_nn12",
			"OLOLO",
			120,
		},
		{
			0,
			0,
			"~nn12",
			"OLOLO",
			0,
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(HandlerFunc))
	defer ts.Close()
	url := ts.URL

	sc := &SearchClient{
		"",
		"",
	}

	for idx, val := range srs {
		sc.URL = url
		sc.AccessToken = accessTocken
		if idx == 0{
			sc.AccessToken = "cscscscdsc"
		}
		if idx == 1{
			sc.URL = "cdsncjncla"
		}

		result, err := sc.FindUsers(val)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}

}
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

//structure of Article
type Article struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	SubTitle           string `json:"subtitle"`
	Content            string `json:"content"`
	Creation_Timestamp string `json:"creation"`
}

//articleHandlers struct is used to handle mutex and store
type articleHandlers struct {
	sync.Mutex
	store map[string]Article // store where articles is stored
}

//articles is used to check what method will be implemented
func (h *articleHandlers) articles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.GetAllArticles(w, r)
		return
	case "POST":
		h.AddArticle(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed")) //given mehod is neither post nor get
		return
	}
}

//GetAllArticles is used for getting all the articles
func (h *articleHandlers) GetAllArticles(w http.ResponseWriter, r *http.Request) {
	articles := make([]Article, len(h.store)) //make slice of articles  to display all articles

	// Add Pagination to get the articles by limit
	Limit := r.URL.Query().Get("limit") //get the limit by url
	limit, _ := strconv.Atoi(Limit)     //convert it into integer
	if limit == 0 {
		limit = 3 //if value of limit is not in url then default value of limit is 3
	}
	Offset := r.URL.Query().Get("offset") //get the offset by url
	offset, _ := strconv.Atoi(Offset)     //convert it into integer
	if offset == 0 {
		offset = 1 //if value of offset is not in url then default value of offset is 0
	}
	page_limit := limit - offset //how many articles will be sent
	h.Lock()
	i := 0
	for _, article := range h.store {
		articles[i] = article
		i++
		if i > page_limit {
			break
		}
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(articles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//AddArticle is used to add articles in store by post method
func (h *articleHandlers) AddArticle(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ht := r.Header.Get("content-type")
	if ht != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', have '%s'", ht)))
		return
	}

	var article Article
	err = json.Unmarshal(bodyBytes, &article)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	article.ID = fmt.Sprint((rand.Intn(100)))        //used to give id i.e between 0 to 100
	article.Creation_Timestamp = time.Now().String() // It gives the current timestamp
	h.Lock()
	h.store[article.ID] = article //store the article
	defer h.Unlock()
}

//GetArticle is used to get an article by ID
func (h *articleHandlers) GetArticle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/") //string is split to get the id by url
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Lock()
	article, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

//SearchArticle is used to search an article by term
func (h *articleHandlers) SearchArticle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Lock()
	term := r.URL.Query().Get("q") //term is retrieved by url
	for _, val := range h.store {  //store is iterated for compare the value with term
		if (val.Title == term) || (val.SubTitle == term) || (val.Content == term) { //term here could be title,subtitle or content
			jsonBytes, _ := json.Marshal(val)
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonBytes)
			break //get out from the loop when we get the matching article
		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	}
	h.Unlock()

}

//newarticleHandlers is used to create an articleHandler for calling the methods
func newarticleHandlers() *articleHandlers {
	return &articleHandlers{
		store: map[string]Article{},
	}

}

//main func is used to start the execution of program
func main() {
	articleHandlers := newarticleHandlers()
	http.HandleFunc("/articles", articleHandlers.articles)
	http.HandleFunc("/articles/", articleHandlers.GetArticle)
	http.HandleFunc("/articles/search", articleHandlers.SearchArticle)

	err := http.ListenAndServe(":8000", nil) //bind the address to 8000 port
	if err != nil {
		panic(err) //if panic occurs then  program will stop execution
	}
}

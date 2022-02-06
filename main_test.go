package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

//TestAddArticle is used to test an article whether is created or not
func TestAddArticle(t *testing.T) {

	var jsonStr = []byte(`{"id":"40","title":"mockingbird","subtitle":"ancient_mariner","content":"tale_of_two_birds","creation":"2022-02-06 13:19:53.809890243 +0530 IST m=+646.059253026"}`)

	req, err := http.NewRequest("POST", "/articles", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err) // if  error occurs while reading the request
	}
	rr := httptest.NewRecorder()                       //create a newrecorder to record the response.
	req.Header.Set("Content-Type", "application/json") //set content-type
	handler := http.HandlerFunc(newarticleHandlers().articles)
	handler.ServeHTTP(rr, req)                      //reading all the http requests and send to the response body
	if status := rr.Code; status != http.StatusOK { //check the status code
		t.Errorf("return status code wrong: have %v require %v",
			status, http.StatusOK)
	}

}

//Test Function for get an article
func TestGetArticle(t *testing.T) {

	req, err := http.NewRequest("GET", "/articles/81", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newarticleHandlers().articles)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("return status code wrong: have %v require %v",
			status, http.StatusOK)
	}

}

//TestGetAllArticles is used to check all articles
func TestGetAllArticles(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newarticleHandlers().GetAllArticles)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("return status code wrong: have %v require %v",
			status, http.StatusOK)
	}

}

//Test function for searching an article
func TestSearchArticle(t *testing.T) {
	req, err := http.NewRequest("GET", "/articles/search?q=tail", nil) //search an article by title,subtitle or content
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(newarticleHandlers().SearchArticle)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("return status code wrong: have %v require %v",
			status, http.StatusOK)
	}

}

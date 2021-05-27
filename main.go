package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Blog struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Time        string `json:"time"`
	Description string `json:"description"`
}

var blogs []Blog = make([]Blog, 0, 100)
var counter int = 1

// create a new blog
func createBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var blog Blog
	_ = json.NewDecoder(r.Body).Decode(&blog)

	blog.Id = counter
	counter++
	blog.Time = time.Now().Format("2006-01-02 15:04:05")

	blogs = append(blogs, blog)

	json.NewEncoder(w).Encode(blog)
}

// get all blogs
func getAllBlogs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Blog []Blog `json:"blogs"`
	}{blogs})
}

// get specific blog
func getBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	value, _ := strconv.Atoi(params["id"])

	for _, blog := range blogs {
		if blog.Id == value {
			json.NewEncoder(w).Encode(struct {
				Blog Blog `json:"blog"`
			}{blog})
			return
		}
	}
	json.NewEncoder(w).Encode(struct {
		Blog [0]bool `json:"blog"`
	}{})
}

// update specific blog
func updateBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var modifiedBlog Blog
	_ = json.NewDecoder(r.Body).Decode(&modifiedBlog)

	params := mux.Vars(r)
	value, _ := strconv.Atoi(params["id"])

	var temp *Blog

	for index := range blogs {
		if blogs[index].Id == value {
			// modify this blog
			temp = &blogs[index]
			(*temp).Title = modifiedBlog.Title
			(*temp).Author = modifiedBlog.Author
			(*temp).Time = time.Now().Format("2006-01-02 15:04:05")
			(*temp).Description = modifiedBlog.Description

			json.NewEncoder(w).Encode(struct {
				Blog Blog `json:"blog"`
			}{(*temp)})
			return
		}
	}
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

// delete a specific blog
func deleteBlog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	value, _ := strconv.Atoi(params["id"])

	index := -1
	for location, blog := range blogs {
		if blog.Id == value {
			index = location
			break
		}
	}
	if index == -1 {
		http.Error(w, "Blog not found", http.StatusNotFound)
	} else {
		blogs = append(blogs[:index], blogs[index+1:]...)
		http.Error(w, "Delete success", http.StatusOK)
	}

}

func main() {
	r := mux.NewRouter()
	fmt.Println("Starting the server...")

	blogs = append(blogs, Blog{Title: "PSA: You can probably try Gmail’s new integrated Chat now", Author: "Mitchell Clark", Time: time.Now().Format("2006-01-02 15:04:05"), Description: "Sample description of post", Id: 1})
	counter++

	r.HandleFunc("/api/blogs", createBlog).Methods("POST")
	r.HandleFunc("/api/blogs", getAllBlogs).Methods("GET")
	r.HandleFunc("/api/blogs/{id}", getBlog).Methods("GET")
	r.HandleFunc("/api/blogs/{id}", updateBlog).Methods("PUT")
	r.HandleFunc("/api/blogs/{id}", deleteBlog).Methods("DELETE")

	http.ListenAndServe(":8080", r)
}

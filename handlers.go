package main

import (
	"net/http"
)

// When a client requests the root path "/", serve the given file.
func ServeIndex(file string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, file)
		} else {
			http.NotFound(w, r)
		}
	})
}

// Allow clients to request this static file, full path included.
//
// Ex: ServeFile("images/cat.jpg") allows GET /images/cat.jpg
func ServeFile(path string) {
	http.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}

// Calls the given function when the client path starts with 'prefix', but
// doesn't expose this prefix to the calling function (i.e. r.URL.Path doesn't
// show it). Useful for simple extraction of parameters from path.
//
// Ex: HandleFuncStripped("/call/", func (...) {
//     param = r.URL.Path
//     print(param)
// }
//
// GET /call/value -> prints "value"
func HandleFuncStripped(prefix string, function http.HandlerFunc) {
	http.Handle(prefix, http.StripPrefix(prefix, http.HandlerFunc(function)))
}

func main() {
	ServeIndex("../web-interact/index.html")
	ServeFile("README.md")
	panic(http.ListenAndServe(":8080", nil))
}

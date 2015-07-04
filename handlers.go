package main

import (
	"net/http"
	"path"
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
func ServeFile(filePath string) {
	http.HandleFunc("/"+filePath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filePath)
	})
}

// Allow clients to request any file from the given directory.
//
// Ex: ServeDir("images") allows GET /images/cat.jpg
func ServeDir(dirPath string) {
	http.Handle(path.Clean("/"+dirPath), http.FileServer(http.Dir(dirPath)))
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
	ServeDir(".")
	panic(http.ListenAndServe(":8080", nil))
}

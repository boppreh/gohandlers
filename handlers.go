package main

import (
	"net/http"
	"os"
	"path"
	"strings"
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

// Allow clients to request any file from the given directory. Does *NOT* allow
// directory listing.
//
// Ex: ServeDir("images") allows GET /images/cat.jpg
func ServeDir(dirPath string) {
	// In case the user passes something like ".".
	cleanPath := path.Clean("/" + dirPath)
	http.HandleFunc(cleanPath, func(w http.ResponseWriter, r *http.Request) {
		p := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if f, err := os.Stat(p); err != nil || f.IsDir() {
			http.NotFound(w, r)
		} else {
			http.ServeFile(w, r, p)
		}
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

// Starts listening for connections on the given port, on all interfaces.
// This is a blocking call.
func Start(port string) {
	panic(http.ListenAndServe(":"+port, nil))
}

func main() {
	ServeDir(".")
	Start("8080")
}

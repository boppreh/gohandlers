package handlers

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strings"
)

var idChars = []rune("0123456789abcdef")
var idLength = 32

// Generates a 32 character random hexadecimal number.
func randId() string {
	b := make([]rune, idLength)
	for i := range b {
		b[i] = idChars[rand.Intn(len(idChars))]
	}
	return string(b)
}

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
	cleanPath := path.Clean("/"+dirPath) + "/"
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
    if !strings.HasSuffix(prefix, "/") {
        prefix += "/"
    }
	http.Handle(prefix, http.StripPrefix(prefix, http.HandlerFunc(function)))
}

// Allow users to POST files to the given URL path, saving these uploaded files
// in the 'storageDir'. For each file uploaded we call
// 'callback(savedPath, originalFilename)', where 'savedPath' is the path for
// the uploaded file stored in disk and 'originalFilename' is the name the user
// sent.
//
// The 'savedPath' retains the extension from 'originalFilename', but with a
// random name to avoid conflicts.
func AllowUpload(urlPath, formKey, storageDir string, callback func(string, string)) {
	http.HandleFunc(urlPath, func(w http.ResponseWriter, r *http.Request) {
		infile, header, err := r.FormFile(formKey)
		if err != nil {
			http.Error(w, "Error parsing uploaded file: "+err.Error(), http.StatusBadRequest)
			return
		}

		originalFilename := header.Filename
		ext := path.Ext(originalFilename)
		storagePath := path.Join(storageDir, randId()+ext)
		outfile, err := os.Create(storagePath)
		if err != nil {
			http.Error(w, "Error creating file: "+err.Error(), http.StatusBadRequest)
			return
		}

		_, err = io.Copy(outfile, infile)
		if err != nil {
			http.Error(w, "Error saving file: "+err.Error(), http.StatusBadRequest)
			return
		}

		callback(storagePath, originalFilename)
	})
}

// Starts listening for connections on the given port, on all interfaces.
// This is a blocking call.
func Start(port string) {
	panic(http.ListenAndServe(":"+port, nil))
}

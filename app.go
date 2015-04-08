/*
AUTHOR:  Daniel Geers
MATR.NR: 2986942
================================================================================

# DOCUMENTATION / SERVER API: #

### baseurl = localhost, port = 8080
### API-Style: <baseurl>:<port>/<command>/<argument>

<baseurl>:<port>/delete/<filename> ** deletes a file from upload directory,
									  prints an error if not existing
<baseurl>:<port>/search/<string>   ** searches for a <string> in all files
									  within the upload dir. returns all hits.
									  ignores all files except txt
<baseurl>:<port>/list			   ** file operations
<baseurl>:<port>/files             ** simply lists all files (you may download
									  single files using right mouse button / dl
<baseurl>:<port>/upload            ** landing page with menu and upload form.
								      automatically extracts *.zip and *.tar.gz
									  files into the webserver's base upload dir
<baseurl>:<port>/zip/<filename>	   ** compresses a file into format zip and
      	  						      shows a download link
<baseurl>:<port>/tar/<filename>	   ** compresses a file into tar.gz format and
									  shows a download link
<baseurl>:<port>/download/
<filename>	                       ** forces download in browser

================================================================================
*/

package main

import (
	"api"
	"config"
	"html/template"
	"htmltemplate"
	"net/http"
)

// GLOBAL:
var DIR = config.DIR // uploads dir

func main() {
	// upload files
	http.HandleFunc("/upload", api.UploadHandler)

	// Download raw file
	http.HandleFunc("/download/", api.DownloadHandler)

	// search function
	http.HandleFunc("/search/", api.SearchHandler)

	// Download as tar file
	http.HandleFunc("/tar/", api.TarHandler)
	// untar files
	http.HandleFunc("/untar/", api.UntarHandler)

	// Download as zip file
	http.HandleFunc("/zip/", api.ZipHandler)

	// unzip files
	http.HandleFunc("/unzip/", api.UnzipHandler)

	// delete a file
	http.HandleFunc("/delete/", api.DeleteHandler)

	// list files
	http.HandleFunc("/list", api.Files)

	//static file handler - just to compare - not "officially implemented" ;)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("uploads"))))

	http.Handle("/", http.HandlerFunc(Index))

	//start server and isten on port 8080, backup files on server error
	http.ListenAndServe(":8080", nil)

}

const sfield = `
<form action="/" name=sfunction method="GET"><input maxLength=1024 size=70
name=s value="" title="Search for it"><input type=submit
value="search" name="searchfield">
</form>
{{if .}}
<a href="search/{{.}}">click here to see all entries for »{{.}}«</a>
<br>
{{end}}
`

var searchtempl = template.Must(template.New("searchfield").Parse(sfield))

func Index(w http.ResponseWriter, r *http.Request) {
	htmltemplate.Head(w, r)
	searchtempl.Execute(w, r.FormValue("s"))
	htmltemplate.Foot(w, r)
}

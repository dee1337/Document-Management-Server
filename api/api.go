package api

import (
	"bytes"
	"compression"
	"config"
	"fmt"
	"html/template"
	"htmltemplate"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"search"
	"strconv"
	"strings"
	"time"
)

var DIR = config.DIR

//Compile templates on start
var templates = template.Must(template.ParseFiles("tmpl/upload.html"))

//Display the named template
func WriteToHtml(w http.ResponseWriter, tmpl string, data interface{}) {
	templates.ExecuteTemplate(w, tmpl+".html", data)
}

//This is where the action happens.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	htmltemplate.Head(w, r)
	switch r.Method {

	case "GET":
		WriteToHtml(w, "upload", nil)
	case "POST":
		//parse the multipart form in the request
		err := r.ParseMultipartForm(100000)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m := r.MultipartForm       // reference to the parsed multipart form
		files := m.File["myfiles"] // get the *fileheaders
		for i, _ := range files {
			//for each fileheader, get a handle to the actual file
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			destination, err := os.Create(DIR + files[i].Filename)
			defer destination.Close()

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// zip & tar.gz handling
			if (strings.HasSuffix(files[i].Filename, "zip") && strings.HasSuffix(files[i].Filename, "tar.gz")) == false {
				if _, err := io.Copy(destination, file); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			/********** ZIP DETECTION *****************/
			if strings.HasSuffix(files[i].Filename, "zip") {
				//create destination file making sure the path is writeable.
				destination, err := os.Create(DIR)
				defer destination.Close()
				if err != nil {
					log.Println(err)
				}
				log.Println("Debug Message: special file extension detected: zip")
				// display(w, "upload", "file extension seems to be zip!!")
				log.Println("unpacking " + (files[i].Filename) + " into " + DIR)
				//unzipFile(DIR + (files[i].Filename))
				compression.Unzip(DIR+(files[i].Filename), DIR)
				//log.Println("deleting " + (files[i].Filename))
			}

			/********** TAR.GZ DETECTION *****************/
			if strings.HasSuffix(files[i].Filename, ".tar.gz") {
				//create destination file making sure the path is writeable.
				destination, err := os.Create(DIR)
				defer destination.Close()
				if err != nil {
					log.Println(err)
				}
				log.Println("Debug Message: special file extension detected: .tar.gz")
				// display(w, "upload", "file extension seems to be zip!!")
				log.Println("unpacking " + (files[i].Filename) + " into " + DIR)
				//unzipFile(DIR + (files[i].Filename))
				compression.UnpackTar(DIR + files[i].Filename)
				//fmt.Println("deleting " + (files[i].Filename))
			}

		}
		//display success message.
		// display(w, "upload", "Upload successful :)  ")
		fmt.Fprintf(w, "Upload successful :)")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	htmltemplate.Foot(w, r)
}

// API to unpack tar files
func UntarHandler(w http.ResponseWriter, r *http.Request) {
	fname := r.URL.Path[len("/untar/"):]
	compression.UnpackTar(DIR + fname)
}

// API to unpack zip files
func UnzipHandler(w http.ResponseWriter, r *http.Request) {
	fname := r.URL.Path[len("/unzip/"):]
	compression.Unzip(DIR+fname, DIR)
}

func TarHandler(w http.ResponseWriter, r *http.Request) {
	htmltemplate.Head(w, r)
	fname := r.URL.Path[len("/tar/"):]
	files := []string{DIR + fname}
	compression.CreateTar(DIR+fname+".tar.gz", files)
	fmt.Fprintf(w, "the compressed file can be downloaded <a href='/download/%s' target='_blank'>here</a> !<br> Do you want to <a href='/delete/%s.zip'> delete %s </a> now?<br><a href='/list'>back</a>", fname+".tar.gz", fname+".tar.gz", fname+".tar.gz")
	htmltemplate.Foot(w, r)
}

func ZipHandler(w http.ResponseWriter, r *http.Request) {
	fname := r.URL.Path[len("/zip/"):]
	files := []string{DIR + fname}
	compression.CreateZip(DIR+fname+".zip", files)
	htmltemplate.Head(w, r)
	fmt.Fprintf(w, "the zipped file can be downloaded <a href='/download/%s' target='_blank'>here</a> !<br> Do you want to <a href='/delete/%s.zip'> delete %s </a> now?", fname+".zip", fname+".zip", fname+".zip")
	htmltemplate.Foot(w, r)
}

// handler that deletes single files by calling /delete/filename
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	f := r.URL.Path[len("/delete/"):]
	err := os.Remove(DIR + f)
	htmltemplate.Head(w, r)
	if err != nil {
		fmt.Fprintf(w, "error: %s (USE your BROWSERs BACK BUTTON)", err)
		return
	} else {
		fmt.Fprintf(w, "%s was deleted successfully :) <a href='/list'>back</a> ", f)
	}
	htmltemplate.Foot(w, r)

}

/* search handler: parses everything after /search/* & uses it as searchterm
on all files within basedir */
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	hits := 0
	searchterm := r.URL.Path[len("/search/"):]
	htmltemplate.Head(w, r)
	fmt.Fprintf(w, "searching for »"+searchterm+"« in %s directory<br>", DIR)

	files, _ := ioutil.ReadDir(DIR)
	for _, f := range files {
		fmt.Fprintf(w, "<ul>")
		if search.Search(searchterm, f.Name()) == true {
			hits++
			fmt.Fprintf(w, "<li>found in %s </li>", f.Name())
		} else {
			//fmt.Println(searchterm, "NOT found in ", f.Name())
		}
		fmt.Fprintf(w, "</ul>")
	}
	fmt.Fprintf(w, "<p>(<b>»"+searchterm+"« was found %s times.</b>) <a href='/list'> back to files</p>", strconv.Itoa(hits))
	htmltemplate.Foot(w, r)
}

/* list the full directory */
func Files(w http.ResponseWriter, r *http.Request) {
	htmltemplate.Head(w, r)
	files, _ := ioutil.ReadDir(DIR)
	for _, f := range files {
		fmt.Fprintf(w, "[<a href='/download/%s'><i class='fa fa-download'></i></a>] [<a href='/tar/%s'><i class='fa fa-file-archive-o'></i></a>] [<a href='/delete/%s'><i class='fa fa-trash-o'></i></a>]  %s<br>", f.Name(), f.Name(), f.Name(), f.Name())
	}
	htmltemplate.Foot(w, r)
}

/*
var Buffer = make([]string, 1000)

// helper traversing the directory
func WalkFn(p string, f os.FileInfo, err error) error {
	i := 0
	current, err := os.Stat(p)
	if current.IsDir() == false {
		Buffer[i] = p
		i++
		//fmt.Fprintf(w, "[<a href='/download/%s'><i class='fa fa-download'></i></a>] [<a href='/tar/%s'><i class='fa fa-file-archive-o'></i></a>] [<a href='/delete/%s'><i class='fa fa-trash-o'></i></a>]  %s<br>", p, p, p, p) //f.Name(), f.Name(), f.Name(), f.Name())
	}
	//fmt.Printf("Visited: %s\n", p)
	return nil
}

func Files2(w http.ResponseWriter, r *http.Request) {
	htmltemplate.Head(w, r)

	err := filepath.Walk(DIR, WalkFn)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		fmt.Fprintf(w, "[<a href='/download/%s'><i class='fa fa-download'></i></a>] [<a href='/tar/%s'><i class='fa fa-file-archive-o'></i></a>] [<a href='/delete/%s'><i class='fa fa-trash-o'></i></a>]  %s<br>", buffer[i], buffer[i], buffer[i], buffer[i]) //f.Name(), f.Name(), f.Name(), f.Name())
	}

	htmltemplate.Foot(w, r)
}
*/

/* forces download of a file into the browser without asking */
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fname := r.URL.Path[len("/download/"):]
	url := DIR + fname
	data, err := ioutil.ReadFile(url)
	if err != nil {
		panic(err)
	}
	// trick to force a browser do instantly download the file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fname)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	http.ServeContent(w, r, DIR, time.Now(), bytes.NewReader(data))
}

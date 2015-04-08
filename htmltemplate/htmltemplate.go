package htmltemplate

import (
	"fmt"
	"net/http"
)

func Head(w http.ResponseWriter, r *http.Request) {
	head := `
    <!DOCTYPE html>
    <html lang="en">
        <head>
        <meta http-equiv='Content-Type' content='text/html; charset=utf-8'>
         <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.3.0/css/font-awesome.min.css">
        </head>

        <nav>

            <a href="/upload"> <i class="fa fa-upload"> upload</i></a>
            <a href="/list"> <i class="fa fa-list"> show files</i></a>
            <a href="/"> <i class="fa fa-search"> search</i></a>
            <!-- <a href="/"> <i class="fa fa-search"> search</i></a> -->

        <br><br>
        </nav>

    `
	fmt.Fprintf(w, head)
}

func UploadForm(w http.ResponseWriter, r *http.Request) {
	form := `<div class="container">
      <h3>Upload Files</h3>
      <div class="message">{{.}}</div>
      <form class="form-signin" method="post" action="/upload" enctype="multipart/form-data">
          <fieldset>
            <input type="file" name="myfiles" id="myfiles" multiple="multiple">
            <input type="submit" name="submit" value="Submit">
        </fieldset>
      </form>
    </div>`
	fmt.Fprintf(w, form)

}

func Foot(w http.ResponseWriter, r *http.Request) {
	foot := `</html>`
	fmt.Fprintf(w, foot)
}

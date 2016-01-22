package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const port = ":3001"
const index = `
<!doctype html>
<form>
<input type=text name=url placeholder="Enter a URL, e.g. http://...." />
<button type=submit>Submit</button>
</form>
<pre>
{{ .Result }}
</pre>
`

var tmpl = template.Must(template.New("index").Parse(index))

type View struct {
	Result string
}

func main() {
	log.Println("starting on port ", port)
	http.HandleFunc("/", proxyHandler)
	http.ListenAndServe(port, nil)
}

func renderIndex(w io.Writer, result string) {
	tmpl.Execute(w, View{Result: result})
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	//show index page
	if len(url) == 0 {
		renderIndex(w, "")
		log.Println("GET / => home")
		return
	}

	resp, err := http.Get(url)

	if err != nil {
		log.Printf("ERROR GET /?url=%s => home\n", url)
		renderIndex(w, fmt.Sprintf("ERROR:%s", err))
		return
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("ERROR GET /?url=%s => home\n", url)
		renderIndex(w, fmt.Sprintf("ERROR:%s", err))
		return
	}

	renderIndex(w, string(data))
}

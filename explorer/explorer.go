package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	blockchain "github.com/devholic77/duckcoin/blockchanin"
)

const (
	tempalteDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// data := homeData{"Home", blockchain.BlockChain().AllBlock()}
	// templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.BlockChain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
	templates.ExecuteTemplate(rw, "add", nil)
}

func Start(port int) {
	handler := http.NewServeMux()
	templates = template.Must(template.ParseGlob(tempalteDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(tempalteDir + "partials/*.gohtml"))
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)
	fmt.Printf("Explorer Listening on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

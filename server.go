package main

import (
	"errors"
	"flag"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
)

var format = render.New()

type Response struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

var (
	port     = flag.String("port", "3201", "service port")
	fontdpi  = flag.Float64("fontdpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "./fonts/DroidSansMono.ttf", "filename of the ttf font")
	fontsize = flag.Float64("fontsize", 12, "font size in points")
)

func main() {
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/preview", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		q := r.Form

		width, _ := strconv.Atoi(q.Get("width"))
		height, _ := strconv.Atoi(q.Get("height"))
		source := q.Get("source")
		target := q.Get("target")

		ext := strings.ToLower(path.Ext(source))
		base := path.Base(source)
		if base == "Dockerfile" {
			ext = ".docker"
		}

		var err error
		switch ext {
		case ".jpg", ".png", ".jpeg":
			err = getImagePreview(source, target, width, height)

		case ".txt":
			err = genTxtPreview(source, target, width, height)

		case ".json", ".js", ".css", ".html", ".htm", ".yaml", ".yml", ".docker", ".xml", ".ts", ".md", ".ini", ".java", ".go", ".sql", ".sh":
			err = genCodePreview(source, target, width, height)

		case ".xls", ".xlsx":
			err = genOfficeOtherPreview(source, target, width, height)

		case ".doc", ".docx":
			err = genOfficeDocPreview(source, target, width, height)

		case ".tiff", ".svg", ".pdf", ".gif", ".webp":
			err = genImagePreview(source, target, width, height)

		default:
			err = errors.New("unsupported file type")
		}

		if err == nil {
			format.JSON(w, 200, Response{Status: true})
		} else {
			format.JSON(w, 200, Response{Status: false, Error: err.Error()})
		}
	})

	log.Printf("starting service ad port %s", *port)
	http.ListenAndServe(":"+*port, r)

}

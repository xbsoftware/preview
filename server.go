package main

import (
	"bytes"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/unrolled/render"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/google/gops/agent"
	"github.com/jinzhu/configor"
)

var format = render.New()

type Response struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

type AppConfig struct {
	Port        string `default:"3201"`
	Key         string
	Text        TextConfig
	UploadLimit int64 `default:"33000000"`
}

type TextConfig struct {
	FontDPI  float64 `default:"72"`
	FontSize float64 `default:"12"`
	FontFile string  `default:"./fonts/DroidSansMono.ttf"`
}

var Config AppConfig

func main() {
	configor.New(&configor.Config{ENVPrefix: "APP"}).Load(&Config, "config.yml")

	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/preview", func(w http.ResponseWriter, r *http.Request) {
		if Config.Key != "" {
			apiKey := r.URL.Query().Get("key")
			if Config.Key != apiKey {
				format.Text(w, 500, "Access denied")
				return
			}
		}

		r.Body = http.MaxBytesReader(w, r.Body, Config.UploadLimit)
		r.ParseMultipartForm(Config.UploadLimit)

		file, handler, err := r.FormFile("file")
		if err != nil {
			format.Text(w, 500, "Can't accept the file")
			return
		}
		defer file.Close()

		q := r.Form

		width, _ := strconv.Atoi(q.Get("width"))
		height, _ := strconv.Atoi(q.Get("height"))
		name := handler.Filename

		ext := strings.ToLower(path.Ext(handler.Filename))
		if handler.Filename == "Dockerfile" {
			ext = ".docker"
		}

		var writer = bytes.Buffer{}
		var outType = "png"
		switch ext {
		case ".txt":
			err = genTxtPreview(file, &writer, width, height)

		case ".json", ".js", ".css", ".html", ".htm", ".yaml", ".yml", ".docker", ".xml", ".ts", ".md", ".ini", ".java", ".go", ".sql", ".sh":
			err = genCodePreview(file, &writer, name, width, height)

		case ".xls", ".xlsx":
			err = genOfficeOtherPreview(file, &writer, name, width, height)

		case ".doc", ".docx":
			err = genOfficeDocPreview(file, &writer, name, width, height)

		case ".jpg", ".png", ".jpeg", ".tiff", ".svg", ".pdf", ".gif", ".webp":
			err = genImagePreview(file, &writer, width, height)
			outType = "jpg"
		default:
			err = errors.New("unsupported file type")
		}

		if err == nil {
			if outType == "png" {
				w.Header().Add("Content-type", "image/png")
			} else {
				w.Header().Add("Content-type", "image/jpg")
			}
			writer.WriteTo(w)
		} else {
			format.Text(w, 500, err.Error())
		}
	})

	log.Printf("starting service ad port %s", Config.Port)
	http.ListenAndServe(":"+Config.Port, r)

}

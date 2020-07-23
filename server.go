package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/unrolled/render"

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

	illegalName := regexp.MustCompile(`[^[:alnum:]-.]`)
	r.Post("/convert", func(w http.ResponseWriter, r *http.Request) {
		file, name, err := incomingFile(w, r)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}
		defer file.Close()

		writer := bytes.Buffer{}
		outName := illegalName.ReplaceAllString(r.Form.Get("name"), "")
		outType := illegalName.ReplaceAllString(r.Form.Get("type"), "")

		if outName == "" || outType == "" {
			format.Text(w, 500, "Invalid conversion target")
			return
		}

		err = convertOffice(file, &writer, name, outType)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}

		w.Header().Add("Content-type", "application/"+outType)
		w.Header().Add("Content-Disposition", "attachment; filename=\""+outName+"\"")
		writer.WriteTo(w)
	})

	r.Post("/preview", func(w http.ResponseWriter, r *http.Request) {
		file, name, err := incomingFile(w, r)
		if err != nil {
			format.Text(w, 500, err.Error())
			return
		}
		defer file.Close()

		q := r.Form

		width, _ := strconv.Atoi(q.Get("width"))
		height, _ := strconv.Atoi(q.Get("height"))

		ext := strings.ToLower(path.Ext(name))
		if name == "Dockerfile" {
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

func incomingFile(w http.ResponseWriter, r *http.Request) (io.ReadCloser, string, error) {
	if Config.Key != "" {
		apiKey := r.URL.Query().Get("key")
		if Config.Key != apiKey {
			return nil, "", errors.New("Access denied")
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, Config.UploadLimit)
	r.ParseMultipartForm(Config.UploadLimit)

	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, "", errors.New("Can't accept the file")
	}

	return file, handler.Filename, nil
}

package doc

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

const Prefix = "/docs"

//go:generate go run openapi_gen.go
func Handler() http.Handler {
	docsURL := Prefix

	specFilename := "openapi.yaml"
	specURL := fmt.Sprintf("%s/%s", docsURL, specFilename)

	specFile, _ := openapi.Open(specFilename)
	defer specFile.Close()
	specData, _ := ioutil.ReadAll(specFile)

	buf := new(bytes.Buffer)
	docTmpl := template.Must(template.New("openapi").Parse(docTemplate))
	_ = docTmpl.Execute(buf, struct {
		Title, RedocURL, SpecURL string
	}{
		Title:    "Payments API Documentation",
		RedocURL: "https://cdn.jsdelivr.net/npm/redoc@2.0.0-rc.8/bundles/redoc.standalone.js",
		SpecURL:  specURL,
	})
	docData := buf.Bytes()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case specURL:
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(specData)
		case docsURL:
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(docData)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})
}

const docTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.2/swagger-ui.css" >
    <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }
      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }
      body
      {
        margin:0;
        background: #fafafa;
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.2/swagger-ui-bundle.js"> </script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.2/swagger-ui-standalone-preset.js"> </script>
    <script>
    window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        url: {{ .SpecURL }},
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })
      // End Swagger UI call region
      window.ui = ui
    }
  </script>
  </body>
</html>`

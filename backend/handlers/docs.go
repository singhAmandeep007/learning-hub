package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"learninghub/config"
	"learninghub/constants"
	contract "learninghub/contract"
)

const swaggerUIHTMLTemplate = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Learning Hub API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
      html, body {
        margin: 0;
        padding: 0;
      }
      #swagger-ui {
        min-height: 100vh;
      }
    </style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui',
        deepLinking: true,
        {{SUBMIT_METHODS}}
        presets: [SwaggerUIBundle.presets.apis],
      });
    </script>
  </body>
</html>
`

// ServeSwaggerUI serves the Swagger UI for API documentation.
func ServeSwaggerUI(c *gin.Context) {
	submitMethodsConfig := ""
	if config.AppConfig != nil && config.AppConfig.ENV_MODE == constants.EnvModeProd {
		submitMethodsConfig = "supportedSubmitMethods: [],"
	}

	html := strings.Replace(swaggerUIHTMLTemplate, "{{SUBMIT_METHODS}}", submitMethodsConfig, 1)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// ServeOpenAPISpec serves the OpenAPI specification in YAML format.
func ServeOpenAPISpec(c *gin.Context) {
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", contract.OpenAPISpecYAML)
}

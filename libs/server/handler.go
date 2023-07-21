package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"simplerest/libs/settings"
	"strings"
)

const (
	maxInt                 = ^int64(0)
	ContentTypeFormEncoded = "application/x-www-form-urlencoded"
	AcceptAny              = "*/*"
	AcceptTextYAML         = "text/yaml"
	AcceptCSV              = "text/csv"
)

type resourceHandler struct {
	res settings.Resource
	db  *sqlx.DB
}

func (h *resourceHandler) failure(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": err.Error(),
	})
}

func (h *resourceHandler) success(c *gin.Context, data interface{}) {
	accept := c.GetHeader("accept")
	cb := c.JSON
	if strings.Contains(accept, gin.MIMEJSON) {
		cb = c.JSON
	} else if strings.Contains(accept, gin.MIMETOML) {
		cb = c.TOML
	} else if strings.Contains(accept, gin.MIMEYAML) || strings.Contains(accept, AcceptTextYAML) {
		cb = c.YAML
	} else if strings.Contains(accept, gin.MIMEXML) || strings.Contains(accept, gin.MIMEXML2) {
		cb = c.XML
	} else if strings.Contains(accept, gin.MIMEHTML) {
		cb = func(code int, obj any){
      c.HTML(code, h.res.Template, obj)
    }
	} else if strings.Contains(accept, AcceptCSV) {
		cb = func(code int, obj any) {
			c.Render(code, render.Data{
				ContentType: "text/csv",
				Data:        []byte("not implemented yet"),
			})
		}
	} else {
		cb = c.JSON
	}
	cb(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func (h *resourceHandler) params(c *gin.Context) gin.H {
	params := gin.H{}

	// from request payload - whether it is a form
	if strings.Contains(c.ContentType(), ContentTypeFormEncoded) {
		c.Request.ParseMultipartForm(maxInt)
		for k, v := range c.Request.PostForm {
			params[k] = strings.Join(v, ",")
		}
	}

	// from request payload - whether it is a json
	if strings.Contains(c.ContentType(), gin.MIMEJSON) {
		if c.Request.Body != nil {
			if data, err := io.ReadAll(c.Request.Body); err == nil {
				var result map[string]interface{}
				json.Unmarshal([]byte(data), &result)
				for k, v := range result {
					params[k] = v
				}
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	// from URL query params
	for k, v := range c.Request.URL.Query() {
		params[k] = strings.Join(v, ",")
	}

	// from URL params
	for _, p := range c.Params {
		params[p.Key] = p.Value
	}

	// from previous middlewares
	if username, found := c.Get("username"); found {
		params["__USERNAME"] = username
	}

	return params
}

func (h *resourceHandler) query(c *gin.Context) {
	params := h.params(c)
	rows, err := h.db.NamedQuery(h.res.Query, params)
	if err != nil {
		h.failure(c, err)
		return
	}
	results := []gin.H{}
	for rows.Next() {
		result := gin.H{}
		if err := rows.MapScan(result); err != nil {
			h.failure(c, err)
			return
		}
		results = append(results, result)
	}
	rows.Close()
	h.success(c, results)
}

func (h *resourceHandler) exec(c *gin.Context) {
	results := gin.H{}
	params := h.params(c)
	result, err := h.db.NamedExec(h.res.Query, params)
	if err != nil {
		h.failure(c, err)
		return
	}

	insertedID, err := result.LastInsertId()
	if err == nil {
		results["inserted_id"] = insertedID
	} else {
		results["inserted_id"] = err
	}

	affectedRows, err := result.RowsAffected()
	if err == nil {
		results["affected_rows"] = affectedRows
	} else {
		results["affected_rows"] = err
	}
	h.success(c, results)
}

func (h *resourceHandler) handle(c *gin.Context) {
	switch h.res.Method {
	case http.MethodGet:
		h.query(c)
	case http.MethodPost, http.MethodPut, http.MethodDelete:
		h.exec(c)
	default:
		c.JSON(http.StatusNotImplemented, gin.H{
			"success": false,
			"message": h.res.Method + " not implemented",
		})
	}
}

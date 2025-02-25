package swagger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Sicarii-NaumaN/poloniex-gateway-api/tools/logger"
	"net/http"
	"sort"
	"strings"
)

// Handler struct
type Handler struct {
	HandlerFunc http.HandlerFunc
	Path        string
	Method      string
	Description string
	Opts        []Option

	ResponseBody interface{}
	RequestBody  interface{}

	IsResponseFile   bool
	ResponseMimeType mimeType
	Tag              string
}

// Parameter struct
type Parameter struct {
	Name string
	Type swaggerFieldType
	// only for Type array
	ItemsType   swaggerFieldType
	Required    bool
	Description string
}

func (h *Handler) AppendRequestBody(schema string) {
	h.Opts = append(h.Opts, bodyOpts{
		BodySchema: schema,
	})
}

func (h *Handler) AddMimeTypeProduce() string {
	if h.ResponseMimeType == "" {
		h.ResponseMimeType = MimeJson
	}

	mimes := []string{
		string(MimeJson),
	}

	if MimeJson != h.ResponseMimeType {
		mimes = append(mimes, string(h.ResponseMimeType))
	}

	if len(mimes) == 2 {
		return fmt.Sprintf(responseProduce, fmt.Sprintf(`"%s", "%s"`, mimes[0], mimes[1]))
	}
	return fmt.Sprintf(responseProduce, fmt.Sprintf(`"%s"`, mimes[0]))
}

func GenerateDoc(nameAPI string, handlers []Handler) (string, error) {
	var (
		err error

		handlersByPath = make(map[string][]Handler, len(handlers))
		paths          = make([]string, 0, len(handlers))
		definitions    = make([]interface{}, 0, len(handlers))
	)

	for _, h := range handlers {
		if _, ok := handlersByPath[h.Path]; !ok {
			handlersByPath[h.Path] = []Handler{h}
		} else {
			handlersByPath[h.Path] = append(handlersByPath[h.Path], h)
		}
	}
	// Default err response
	definitions = append(definitions, RpcStatus{})

	for pathHandler, hdls := range handlersByPath {
		var handlersPath = make([]string, 0, len(hdls))
		for _, h := range hdls {
			var (
				commonParams   = make([]string, 0)
				responseSchema = EmptyObject
			)

			if h.ResponseBody != nil {
				definitions = append(definitions, h.ResponseBody)
				responseSchema = GetDefSchema(h.ResponseBody)
			}

			if h.IsResponseFile {
				responseSchema = GetFileResponse()
			}
			commonParams = append(commonParams, h.AddMimeTypeProduce())

			if h.RequestBody != nil {
				definitions = append(definitions, h.RequestBody)
				var reqSchema string
				reqSchema = GetDefSchema(h.RequestBody)
				h.AppendRequestBody(reqSchema)
			}

			tag := nameAPI
			if len(h.Tag) != 0 {
				tag = h.Tag
			}
			consumesProduces, params := MergeOptionsJSON(h.Opts...)
			commonParams = append(commonParams, params)
			commonParams = append(consumesProduces, commonParams...)
			var p = fmt.Sprintf(
				HandlerRaw1,
				strings.ToLower(h.Method),
				h.Description,
				strings.Join(commonParams, ","),
				responseSchema,
				GetDefSchema(RpcStatus{}),
				"",
				tag,
			)
			handlersPath = append(handlersPath, p)
		}
		paths = append(paths, fmt.Sprintf(HandlerBrakesRaw, pathHandler, strings.Join(handlersPath, ",")))
	}

	sort.Slice(paths, func(i, j int) bool {
		return paths[i] < paths[j]
	})

	defs := BuildDefs(definitions...)
	data, err := json.Marshal(defs)
	if err != nil {
		logger.Fatal(context.Background(), err)
	}

	return fmt.Sprintf(SwaggerJSON, nameAPI, strings.Join(paths, ","), string(data)), nil
}

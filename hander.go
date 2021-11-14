package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
	ErrorCode   int    `json:"error_code"`
	RequestId   string `json:"request_id"`
	OriginalURI string `json:"original_uri"`
	RayId       string `json:"ray_id"`   
	ClientIp 		string `json:"client_ip"`
	Message     string `json:"message"`
}

var ErrorMap = map[int]string{
	403: "Access Denied",
	404: "Not Found",
	413: "Request Too Large",
	502: "Bad Gateway",
	503: "Service Unavailable Error",
}

func errorHandler(t *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(FormatHeader, r.Header.Get(FormatHeader))
		w.Header().Set(CodeHeader, r.Header.Get(CodeHeader))
		w.Header().Set(ContentType, r.Header.Get(ContentType))
		w.Header().Set(OriginalURI, r.Header.Get(OriginalURI))
		w.Header().Set(RequestId, r.Header.Get(RequestId))

		format := r.Header.Get(FormatHeader)
		if format != "application/json" {
			format = DefaultFormat
		}

		w.Header().Set(ContentType, format)

		errCode := r.Header.Get(CodeHeader)
		code, err := strconv.Atoi(errCode)
		if err != nil {
			code = 404
		}
		w.WriteHeader(code)

		message, ok := ErrorMap[code]
		if !ok {
			message = "Unknown Error"
		}

		resp := Response{
			ErrorCode:   code,
			RequestId:   r.Header.Get(RequestId),
			OriginalURI: r.Header.Get(OriginalURI),
			ClientIp:    r.Header.Get(ClientIp),
			RayId:  		 r.Header.Get(RayId),
			Message:     message,
		}

		if format == "application/json" {
			respContent, err := json.Marshal(&resp)
			if err != nil {
				log.Printf("Marshal json error: %v\n", err)
				return
			}
			if _, err = w.Write(respContent); err != nil {
				log.Printf("Write response failed with err %v\n", err)
			}
			return
		}

		if err = t.ExecuteTemplate(w, DefaultErrorTemplateName, resp); err != nil {
			log.Printf("Execute template failed with error: %v\n", err)
		}
	}
}

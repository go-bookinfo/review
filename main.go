package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

type ratings []struct {
	Id   int `json:"Id"`
	Star int `json:"Star"`
}

type reviewes struct {
	Id       int    `json:"Id"`
	Star     int    `json:"Star"`
	Reviewer string `json:"Reviewer"`
	Review   string `json:"Review"`
	Color    string `json:"color"`
}

func Init() (opentracing.Tracer, io.Closer) {
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	tracer, closer := jaeger.NewTracer(
		"",
		jaeger.NewConstSampler(false),
		jaeger.NewNullReporter(),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
	)
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

func Extract(r *http.Request) (string, opentracing.SpanContext, error) {
	requestID := r.Header.Get("x-request-id")
	spanCtx, err :=
		opentracing.GlobalTracer().Extract(
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(r.Header))
	return requestID, spanCtx, err
}

func Inject(spanContext opentracing.SpanContext, request *http.Request, requestID string) error {
	request.Header.Add("x-request-id", requestID)
	return opentracing.GlobalTracer().Inject(
		spanContext,
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func main() {
	_, closer := Init()
	defer closer.Close()

	reviewer1 := reviewes{
		Id:       1,
		Star:     0,
		Reviewer: "Reviewer1",
		Review:   "An extremely entertaining play by Shakespeare. The slapstick humour is refreshing!",
		Color:    "red",
	}
	reviewer2 := reviewes{
		Id:       2,
		Star:     0,
		Reviewer: "Reviewer2",
		Review:   "Absolutely fun and entertaining. The play lacks thematic depth when compared to other plays by Shakespeare.",
		Color:    "",
	}
	http.HandleFunc("/review", func(w http.ResponseWriter, r *http.Request) {
		requestID, ctx, _ := Extract(r)
		req, _ := http.NewRequest("GET", "http://rating/rating", nil)
		Inject(ctx, req, requestID)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", string(data))

		var rat ratings
		json.Unmarshal(data, &rat)

		for _, v := range rat {
			if v.Id == 1 {
				reviewer1.Star = v.Star
			} else if v.Id == 2 {
				reviewer2.Star = v.Star
			}
		}

		reviewer := []reviewes{reviewer1, reviewer2}
		bs, err := json.Marshal(reviewer)
		if err != nil {
			fmt.Println(err)
		}
		w.Write(bs)
	})
	http.ListenAndServe(":80", nil)

}

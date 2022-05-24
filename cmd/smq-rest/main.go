package main

/*
   POST   /mq/:channel               (publish new message)
   POST   /mq/:channel/:subscription (add new subscription)
   GET    /mq/:channel/:subscription (get messages published)
   DELETE /mq/:channel/:subscription (remove a subscription)
*/
import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/smartmq/smartmq"
	"github.com/smartmq/smartmq/cmd/smq-rest/restapi"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

var _redisURL string

func main() {
	// setup redis conenction
	url := flag.String("url", "redis://127.0.0.1:6379", "redis url")
	redisURL, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		redisURL = *url
	}
	_redisURL = redisURL

	// setup rest server
	router := NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"Publish",
		"GET",
		"/",
		Info,
	},
	Route{
		"Publish",
		"PUT",
		"/publish",
		Publish,
	},
	Route{ /*POST con url fissa e payload json con la 'subscription' */
		"AddSubscription",
		"POST",
		"/subscribe/{subscription}",
		AddSubscription,
	},
	Route{
		"GetMessages",
		"GET",
		"/subscribe/{subscription}",
		GetMessages,
	},
	Route{
		"RemoveSubscription",
		"DELETE",
		"/mq/{channel}/{subscription}",
		RemoveSubscription,
	},
}

func Info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status: %s", "OK")
}

func Publish(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//channel := vars["channel"]
	//channel := r.Header.Get("channel")
	channel := r.URL.Query().Get("channel")

	//msgContent := string(getBody(r))
	msgContent := getBody(r)

	mq := smartmq.New(_redisURL, false)
	mq.Channel(channel).Publish(msgContent)
	mq.Close()
}
func getBody(r *http.Request) []byte {
	msgContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	return msgContent
}

func AddSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//channel := vars["channel"]
	subscription := vars["subscription"]
	channel := r.URL.Query().Get("channel")

	mq := smartmq.New(_redisURL, false)
	mq.Channel(channel).AddSubscription(subscription)
	mq.Close()

	subs := restapi.Subscription{
		Channel: channel,
		Name:    subscription,
	}

	if err := json.NewEncoder(w).Encode(subs); err != nil {
		panic(err)
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//channel := vars["channel"]
	subscription := vars["subscription"]
	channel := r.URL.Query().Get("channel")

	mq := smartmq.New(_redisURL, false)

	mq.Channel(channel).AddSubscription(subscription)

	//msgContent := mq.Channel(channel).Subscription(subscription).GetMessage()
	msgContent := mq.Channel(channel).Subscription(subscription).GetMessageBlocking()
	mq.Close()
	//log.Printf("SUBSCRIBE: %v\n", msgContent)

	if msgContent != nil && len(msgContent) > 0 {
		contentType := http.DetectContentType(msgContent)
		w.Header().Add("Content-Type", contentType)
		w.Write(msgContent)
	}
}
func RemoveSubscription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channel := vars["channel"]
	subscription := vars["subscription"]

	mq := smartmq.New(_redisURL, false)
	mq.Channel(channel).Subscription(subscription).RemoveSubscription()
	mq.Close()

	ret := restapi.Operation{
		Status: "ok",
	}
	if err := json.NewEncoder(w).Encode(ret); err != nil {
		panic(err)
	}
}

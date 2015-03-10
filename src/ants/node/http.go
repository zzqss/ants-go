package node

import (
	Http "ants/http"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Router struct {
	Node *Node
	Mux  map[string]func(http.ResponseWriter, *http.Request)
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.Println("get request:" + url)
	if !this.Node.Cluster.IsReady() {
		w.Write([]byte("sorry,cluster not ready,please wait"))
		return
	}
	path := r.URL.Path
	if h, ok := this.Mux[path]; ok {
		h(w, r)
		return
	}
	this.Welcome(w, r)
}

type WelcomeInfo struct {
	Message  string
	Greeting string
	Time     string
}

func (this *Router) Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(Http.CONTENT_TYPE, Http.JSON_CONTENT_TYPE)
	now := time.Now().Format("2006-01-02 15:04:05")
	welcome := WelcomeInfo{
		"for crawl",
		"do not panic",
		now,
	}
	encoder, err := json.Marshal(welcome)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(encoder)
}
func (this *Router) Spiders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	spiderList := make([]string, 0, len(this.Node.Crawler.SpiderMap))
	for spider := range this.Node.Crawler.SpiderMap {
		spiderList = append(spiderList, spider)
	}
	encoder, err := json.Marshal(spiderList)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(encoder)
}

func (this *Router) Crawl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	spiderName := r.Form["spider"][0]
	now := time.Now().Format("2006-01-02 15:04:05")
	result := this.Node.StartSpider(spiderName)
	result.Time = now
	encoder, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(encoder)
}

func (this *Router) Cluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder, err := json.Marshal(this.Node.Cluster)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(encoder)
}

func NewRouter(node *Node) *Router {
	router := &Router{}
	router.Node = node
	router.Mux = make(map[string]func(http.ResponseWriter, *http.Request))
	router.Mux["/"] = router.Welcome
	router.Mux["/cluster"] = router.Cluster
	router.Mux["/spiders"] = router.Spiders
	router.Mux["/crawl"] = router.Crawl
	return router
}

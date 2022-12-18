package fastcache_lru

import (
	"fmt"
	"log"
	"net/http"
)

func runStatsServer(s *stats) {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		data := fmt.Sprintf(`<a href="/stats">/stats</a>
`)
		res.Write([]byte(data))
	})
	http.HandleFunc("/stats", func(res http.ResponseWriter, req *http.Request) {
		data := `<a href="/">/</a>
`
		data += fmt.Sprintf("%+v", s)
		res.Write([]byte(data))
	})
	log.Println("run stats server on 6377...")
	go http.ListenAndServe("localhost:6377", nil)
}

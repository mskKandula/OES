package websock

import (
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
)

var upgrader = ws.HTTPUpgrader{}

func Upgrade(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	conn, _, _, err := upgrader.Upgrade(r, w)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

package websock

import (
	"log"
	"net"
	"net/http"

	"github.com/gobwas/ws"
)

func Upgrade(w http.ResponseWriter, r *http.Request) (net.Conn, error) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

package control

import (
	"fmt"
	"net/http"
)

func (h DecoratedHandler) ApiShortenHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Routing check.Hi from JSON URL")
}

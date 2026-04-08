package controller

import (
	"github.com/christiandsol/main/errUtil"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("GOT HOME"))
	errUtil.CheckErr("Writing buf", nil, err)
}

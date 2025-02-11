package handler

import (
	"context"
	"net/http"
)

type Service interface {
}

func New(ctx context.Context, mux *http.ServeMux, service Service) {

}

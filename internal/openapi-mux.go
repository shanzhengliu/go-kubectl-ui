package internal

import (
	"context"
	"net/http"
)

type EnhancedMux struct {
	mux        *http.ServeMux
	globalPath string // 这里存储你的全局变量
}

func NewEnhancedMux(globalPath string) *EnhancedMux {
	return &EnhancedMux{
		mux:        http.NewServeMux(),
		globalPath: globalPath,
	}
}

func (e *EnhancedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), "globalPath", e.globalPath)
	e.mux.ServeHTTP(w, r.WithContext(ctx))
}

func (e *EnhancedMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	e.mux.HandleFunc(pattern, handler)
}

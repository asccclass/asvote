// router.go
package main

import(
   "os"
   "github.com/gorilla/mux"
   "github.com/asccclass/serverstatus"
   "github.com/asccclass/staticfileserver"
)

// Create your Router function
func NewRouter(srv *SherryServer.ShryServer, documentRoot string)(*mux.Router) {
   router := mux.NewRouter()

   srv.SSE.AddRouter(router) // add sse

   // Google Login
   vote, _ := NewVote(srv)
   vote.AddRouter(router)  // Google 認證後回傳的資訊

   //logger
   router.Use(SherryServer.ZapLogger(srv.Logger))
   // Crawer
   srv.Crawer.AddRouter(router)

   // health check
   systemName := os.Getenv("SystemName")
   m := serverstatus.NewServerStatus(systemName)
   router.HandleFunc("/healthz", m.Healthz).Methods("GET")

   // Static File server
   staticfileserver := SherryServer.StaticFileServer{documentRoot, "index.html"}
   router.PathPrefix("/").Handler(staticfileserver)

   return router
}

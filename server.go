package main
  
import (
   "os"
   "github.com/joho/godotenv"
   "github.com/asccclass/staticfileserver"
)

func main() {
   if err := godotenv.Load("envfile"); err != nil {
      panic(err)
   }
   port := os.Getenv("PORT")
   if port == "" {
      port = "10028"
   }
   documentRoot := os.Getenv("DocumentRoot")
   if documentRoot == "" {
      documentRoot = "www/html"
   }
   templateRoot := os.Getenv("TemplateRoot")
   if templateRoot == "" {
      templateRoot = "www/template"
   }

   server, err := SherryServer.NewServer(":" + port, documentRoot, templateRoot)
   if err != nil {
      panic(err)
   }

   // if you have your own router add this and implement router.go
   server.Server.Handler = NewRouter(server, documentRoot)

   server.Start()
}

package main

import(
   // "os"
   "fmt"
   // "context"
   "strings"
   "net/http"
   "database/sql"
   "encoding/json"
   "github.com/gorilla/mux"
   _ "github.com/mattn/go-sqlite3"
   "github.com/asccclass/sherrytime"
   "github.com/asccclass/staticfileserver" 
)

type UserProfile struct {
   ID		string		`json:"id"`
   Email	string		`json:"email"`
   VerifiedEmail	bool	`json:"verified_email"`
   Name		string		`json:"name"`
   GivenName	string		`json:"given_name"`
   FamiltName	string		`json:"family_name"`
   Picture	string		`json:"picture"`
   Locale	string		`json:"locale"`
}

type SryVote struct {
   Srv          *SherryServer.ShryServer
}

// 儲存投票資訊
func(app *SryVote) Save(voteNo string, profile *UserProfile)(int, error) {
   num := 0
   if voteNo == "" || profile.ID == "" {
      return num, fmt.Errorf("voteNo or id is empty")
   }

   db, err := sql.Open("sqlite3", "./data/vote.db")
   if err != nil {
      return num, err
   }
   defer db.Close()
   st := sherrytime.NewSherryTime("Asia/Taipei", "-")  // Initial
   today := st.Today()
   row := db.QueryRow("select count(*) from votez where voteNo=? and googleID=? and strftime('%Y-%m-%d', voreDate)=?", voteNo, profile.ID, today)
   if err := row.Scan(&num); err != nil {
      return num, err
   }
   if num > 0 {
      return num, fmt.Errorf("<script>alert('您已經投過票(You have already voted)');window.close();</script>")
   }

   sql := "insert into votez(googleID,email,voteNo,voteDate) values(?,?,?,?)"
   stmt, err := db.Prepare(sql)
   if err != nil {
      return num, err
   }
   _, err = stmt.Exec(profile.ID, profile.Email, voteNo, st.Now())
   if err != nil {
      return num, err
   }
   row = db.QueryRow("select count(*) from votez where voteNo=?", voteNo)
   if err := row.Scan(&num); err != nil {
      return num, err
   }
   return num, nil
}

// call back
func(app *SryVote) Callback(w http.ResponseWriter, r *http.Request) {
   state := r.FormValue("state")   // "Google-使用者投某作品之編號" || Facebook-使用者投某作品之編號
   code := r.FormValue("code")
   if state == "" || code == "" {
      app.Srv.Error.Error2Web(w, fmt.Errorf("State or Code error"))
      return
   }
   voteData := strings.Split(state, "-")
   voteData[0] = strings.ToLower(voteData[0])
   if len(voteData) != 2 || (voteData[0] != "google" && voteData[0] != "facebook") {
      app.Srv.Error.Error2Web(w, fmt.Errorf("Vote params error"))
      return
   }
   // 取得資料
   num := 0
   data := []byte("")
   if voteData[0] == "google" {
      g, err := app.Srv.LineLogin.NewGoogleLogin() 
      if err != nil {
         app.Srv.Error.Error2Web(w, err)
         return
      }
      data, err = g.GetUserProfile(code)
      if err != nil {
         app.Srv.Error.Error2Web(w, err)
         return
      }
   } else if voteData[0] == "facebook" {
      f, err := app.Srv.LineLogin.NewFacebookLogin()
      if err != nil {
         app.Srv.Error.Error2Web(w, err)
         return
      }
      data, err = f.GetUserProfile(code)
      if err != nil {
         app.Srv.Error.Error2Web(w, err)
         return
      }
   } else {
      app.Srv.Error.Error2Web(w, fmt.Errorf("不支援"))
      return
   }
   // 儲存投票資料
   var profile UserProfile
   var err error
   if err := json.Unmarshal(data, &profile); err != nil {
      app.Srv.Error.Error2Web(w, err)
      return
   }
   num, err = app.Save(voteData[1], &profile)
   if err != nil {
      fmt.Fprintf(w, err.Error())
      return
   }
   // 投票完成
   w.Header().Set("Content-Type", "text/html;charset=UTF-8")
   w.Header().Set("Cross-Origin-Opener-Policy", "same-origin-allow-popups")
   w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
   w.Header().Set("Cross-Origin-Opener-Policy", "unsafe-none")
   w.WriteHeader(http.StatusOK)
   fmt.Fprintf(w, "<html><body><script>alert('" + voteData[1] + "號投票完畢');window.close();</script><body></html>")

   sseStr := fmt.Sprintf("{\"stNo\":%s, \"count\":%d}", voteData[1], num)
fmt.Println(sseStr)
   go app.Srv.SSE.ReplyMessage("", sseStr)
}

// Google Login
func(app *SryVote) AddRouter(router *mux.Router) {
   app.Srv.LineLogin.AddRouter(router)		       // 加入 Line Login
   router.HandleFunc("/g/callback", app.Callback)      // Google 認證後回傳的資訊
   router.HandleFunc("/f/callback", app.Callback)      // Google 認證後回傳的資訊
   router.HandleFunc("/status/{voteNo}", app.GetStatusFromWeb).Methods("GET") // 取得目前全部的投票資訊
}

func NewVote(srv *SherryServer.ShryServer)(*SryVote, error) {
   return &SryVote {
      Srv: srv,
   }, nil
}

package main

import(
   "fmt"
   "net/http"
   "database/sql"
   "encoding/json"
   "github.com/gorilla/mux"
   _ "github.com/mattn/go-sqlite3"
   // "github.com/asccclass/sherrytime"
)
/*
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
*/

type Resultz struct {
   StNo		string		`json:"stNo"`
   Count	string		`json:"count"`
}

// 取得投票資訊
func(app *SryVote) Status(voteNo string)(string, error) {
   result := "" 
   if voteNo == "" {
      return result, fmt.Errorf("voteNo is empty")
   }

   db, err := sql.Open("sqlite3", "./data/vote.db")
   if err != nil {
      return result, err
   }
   defer db.Close()
   
   if voteNo == "all" {
      rows, err := db.Query("select voteNo,count(*) from votez group by voteNo")
      if err != nil {
         return result, err
      }
      defer rows.Close()
      data := []Resultz{}
      for rows.Next() {
         elem := Resultz{}
         if err := rows.Scan(&elem.StNo, &elem.Count); err != nil {
            return result, err
         }
         data = append(data, elem)
      }
      jsondata, _ := json.Marshal(data)
      result = string(jsondata)
   } else { // 個別投票資訊
      num := 0
      row := db.QueryRow("select count(*) from votez where voteNo=?", voteNo)
      if err := row.Scan(&num); err != nil {
         return result, err
      }
      result = fmt.Sprintf("[{\"voteNo\":%s, \"count\":%d}]", voteNo, num)
   }
   return result, nil
}

// 取得投票資訊
func(app *SryVote) GetStatusFromWeb(w http.ResponseWriter, r *http.Request) {
   w.Header().Set("Content-Type", "application/json;charset=UTF-8")
   w.WriteHeader(http.StatusOK)

   webVars := mux.Vars(r)
   if webVars["voteNo"] == "" {
      webVars["voteNo"] = "all"
   }

   result, err := app.Status(webVars["voteNo"])
   if err != nil {
      app.Srv.Error.Error2Web(w, err)
      return
   }
   fmt.Fprintf(w, result)
}

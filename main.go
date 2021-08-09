package main

import (
 "database/sql"
 "log"
 "net/http"
 "text/template"
 _ "github.com/go-sql-driver/mysql"
)

type Myemployee struct {
 Id 		int
 Full_Name	string
 Address	string
}

func dbConn() (db *sql.DB) {
 dbDriver := "mysql"
 dbUser := "root"
 dbPass := "password123"
 dbName := "gocrud"
 db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
 if err != nil {
  panic(err.Error())
 }
 return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

func Index(w http.ResponseWriter, r *http.Request) {
 db := dbConn()
 selDB, err := db.Query("select * from myemployee order by id desc")
 if err != nil {
  panic(err.Error())
 }
 emp := Myemployee{}
 res := []Myemployee{}
 for selDB.Next() {
  var id int
  var fullname, address string
  err = selDB.Scan(&id, &fullname, &address)
  if err != nil {
   panic(err.Error())
  }
  emp.Id = id
  emp.Full_Name = fullname
  emp.Address = address
  res = append(res, emp)
 }
 tmpl.ExecuteTemplate(w, "Index", res)
 defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
 db := dbConn()
 nId := r.URL.Query().Get("id")
 selDB, err := db.Query("select * from myemployee where id=?", nId)
 if err != nil {
  panic(err.Error())
 }
 emp := Myemployee{}
 for selDB.Next() {
  var id int
  var fullname, address string
  err = selDB.Scan(&id, &fullname, &address)
  if err != nil {
   panic(err.Error())
  }
  emp.Id = id
  emp.Full_Name = fullname
  emp.Address = address
 }
 tmpl.ExecuteTemplate(w, "Show", emp)
 defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
 tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
 db := dbConn()
 nId := r.URL.Query().Get("id")
 selDB, err := db.Query("select * from myemployee where id=?", nId)
 if err != nil {
  panic(err.Error())
 }
 emp := Myemployee{}
 for selDB.Next() {
  var id int
  var fullname, address string
  err = selDB.Scan(&id, &fullname, &address)
  if err != nil {
   panic(err.Error())
  }
  emp.Id = id
  emp.Full_Name = fullname
  emp.Address = address
 }
 tmpl.ExecuteTemplate(w, "Edit", emp)
 defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
 db := dbConn()
 if r.Method == "POST" {
  fullname := r.FormValue("fullname")
  address := r.FormValue("address")
  insForm, err := db.Prepare("insert into myemployee (fullname, address) values(?,?)")
  if err != nil {
   panic(err.Error())
  }
  insForm.Exec(fullname, address)
  log.Println("Insert: FullName: " + fullname + " | Address: " + address)
 }
 defer db.Close()
 http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
 db := dbConn()
 if r.Method == "POST" {
  fullname := r.FormValue("fullname")
  address := r.FormValue("address")
  id := r.FormValue("uid")
  insForm, err := db.Prepare("update myemployee set fullname=?, address=? where id=?")
  if err != nil {
   panic(err.Error())
  }
  insForm.Exec(fullname, address, id)
  log.Println("Update: Fullname: " + fullname + " | Address : " + address)
 }
 defer db.Close()
 http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
 db := dbConn()
 emp := r.URL.Query().Get("id")
 delForm, err := db.Prepare("delete from myemployee where id=?")
 if err != nil {
  panic(err.Error())
 }
 delForm.Exec(emp)
 log.Println("DELETE")
 defer db.Close()
 http.Redirect(w, r, "/", 301)
}

func main() {
 log.Println("Server started on: http://localhost:8080")
 http.HandleFunc("/", Index)
  http.HandleFunc("/show", Show) 
  http.HandleFunc("/new", New) 
  http.HandleFunc("/edit", Edit) 
  http.HandleFunc("/insert", Insert) 
  http.HandleFunc("/update", Update) 
  http.HandleFunc("/delete", Delete) 
  http.ListenAndServe(":8080", nil)
}

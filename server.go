package main

import (
    "log"
    "net/http"
    "encoding/json"

    "time"
    "github.com/googollee/go-socket.io"
)

type Profile struct {
  Name    string
  Hobbies []string
  Age int
}

type customServer struct {
  Server *socketio.Server
}

func respond(w http.ResponseWriter, r *http.Request) {
  profile := Profile{"Alex", []string{"snowboarding", "programming"}, 20}

  js, err := json.Marshal(profile)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
  defer timeTrack(time.Now(), "response")
}

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    log.Printf("%s took %s", name, elapsed)
}

func handleChatMessage(so socketio.Socket, msg string)  {
  log.Println("emit:", so.Emit("chat message", msg))
  so.BroadcastTo("chat", "chat message", msg)
}

func handleDisconnection()  {
  log.Println("on disconnect")
}

var MyServerName = "http://localhost"

func enableCors(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Access-Control-Allow-Origin", MyServerName)
    fn(w, r)
  }
}

func main() {
    server, err := socketio.NewServer(nil)


    if err != nil {
        log.Fatal(err)
    }

    server.On("connection", func(so socketio.Socket) {
        log.Println("on connection")
        so.Join("chat")

        so.On("chat message", handleChatMessage)
        so.On("disconnection", handleDisconnection)
    })

    server.On("error", func(so socketio.Socket, err error) {
        log.Println("Foobar")
        log.Println("error:", err)
    })

    http.Handle("/socket.io/", server)
    http.HandleFunc("/", enableCors(respond))

    log.Fatal(http.ListenAndServe(":3001", nil))
}

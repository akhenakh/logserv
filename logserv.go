package main

import (
	"fmt"
	"os"
	"flag"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"code.google.com/p/go.net/websocket"
	"github.com/howeyc/fsnotify"
	auth "github.com/abbot/go-http-auth"
	"io"
	"log"
)

var config Config

func rootHandler(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	io.WriteString(w, `<html>
<head>
    <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.4.2/jquery.min.js"></script>

  <script>
    window.onload = function() {
        var ws = new WebSocket("ws://websocket");

        ws.onopen = function() {
            ws.send('New participant joined');
			alert('toto');
        };

        ws.onmessage = function (evt)  {
            $("#chat").append("<div>" + evt.data + "</div>");
        };
    };
    </script>
</head>
<body><h3>Logs</h3>
<div id="chat" style="width: 60em; height: 20em; overflow:auto; border: 1px solid black">
</div>
</body></html>`)
}


func EchoServer(ws *websocket.Conn) {
	fmt.Print(ws.Request())
}

type Logfile struct {
      Path string
      Users []string
}

type Config struct {
	Port int
	Auth_file string
	Log_files []Logfile
	Channel_map map[string]chan *os.File
}

func main() {

	var config_path = flag.String("config",  "", "Path to the config file")
	flag.Parse()
	
	if *config_path == "" { log.Fatal("Give path to the config file with -path=") }

	b, err := ioutil.ReadFile(*config_path)
	if err != nil { log.Fatal(err) }

	err = json.Unmarshal(b, &config)
	if err != nil { log.Fatal(err) }
	
	if config.Port == 0 || config.Auth_file == "" || len(config.Log_files) == 0 { log.Fatal("Invalid configuration") }
	
	watcher, err := fsnotify.NewWatcher()
	if err != nil { log.Fatal(err) }
	
    go func() {
           for {
               select {
               case ev := <-watcher.Event:
                   log.Println("event:", ev)
               case err := <-watcher.Error:
                   log.Println("error:", err)
               }
           }
       }()
	   
	for _, v := range(config.Log_files) {
		log.Print(v.Path)
		err = watcher.Watch(v.Path)
		if err != nil { log.Fatal(err) }
	}

	secrets := auth.HtpasswdFileProvider(config.Auth_file)
	authenticator := auth.BasicAuthenticator("logserv", secrets)
	
	http.HandleFunc("/", authenticator(rootHandler))
	
	http.Handle("/websocket", websocket.Handler(EchoServer))
	listen_address := fmt.Sprintf(":%d", config.Port)
	http.ListenAndServe(listen_address, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
	
	watcher.Close()
}
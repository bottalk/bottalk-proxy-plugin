package main

import (
 "bytes"
 "encoding/json"
 "io/ioutil"
 "log"
 "net/http"

 bottalk "github.com/bottalk/go-plugin"
)

type btRequest struct {
 Token  string            `json:"token"`
 UserID string            `json:"user"`
 Vars   map[string]string `json:"vars"`
 Input  json.RawMessage   `json:"input"`
}

func errorResponse(message string) string {
 return "{\"result\": \"fail\",\"message\":\"" + message + "\"}"
}

func main() {

 plugin := bottalk.NewPlugin()
 plugin.Name = "Simple Proxy Plugin"
 plugin.Description = "This plugin proxies your alexa/google requests entirely to endpoints"

 plugin.Actions = map[string]bottalk.Action{"call": bottalk.Action{
  Name:        "call",
  Description: "This action calls remote endpoint",
  Endpoint:    "/call",
  Action: func(r *http.Request) string {

   var BTR btRequest
   decoder := json.NewDecoder(r.Body)

   err := decoder.Decode(&BTR)
   if err != nil {
    return errorResponse(err.Error())
   }

   if len(BTR.Vars["url"]) < 5 {
    return errorResponse("Call webhook is incorrect")
   }

   res, err := http.Post(BTR.Vars["url"], "application/json", bytes.NewBuffer([]byte(BTR.Input)))
   if err != nil {
    return errorResponse(err.Error())
   }
   output, _ := ioutil.ReadAll(res.Body)
   log.Println(string(output))

   return "{\"result\": \"ok\",\"output\":" + string(output) + "}"
  },
  Params: map[string]string{"url": "Endpoint url to call"},
 }}

 plugin.Run(":9080")
}

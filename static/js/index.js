"use strict";

if (!window["WebSocket"]) {
  document.body.innerHTML = "<h1>Sorry, your browser does not support this experiment.</h1>"
} else {
  let conn = new WebSocket("ws://" + document.location.host + "/ws");

  conn.onclose = function(e) {
    console.log("connection closed");
    document.body.innerHTML = "<h1>Connection Closed.</h1>";
  }

  conn.onmessage = function(e) {
    document.body.style.backgroundColor = e.data;
    console.log(e.data);
  }
}
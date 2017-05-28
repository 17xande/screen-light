"use strict";

const message = document.querySelector('#message')

if (!window["WebSocket"]) {
  message.textContent = "Sorry, your browser does not support this experiment."
} else {
  let socket = new WebSocket("ws://" + document.location.host + "/ws");

  socket.addEventListener('open', e => {
    message.textContent = "Connected!";
  });
  socket.addEventListener('error', e => {
    message.textContent = "Oops, something went wrong";
  });
  socket.addEventListener('close', e => {
    message.textContent = "Lost connection, trying to reconnect...";
  });
  socket.addEventListener('message', e => {
    document.body.style.backgroundColor = e.data;
  });
}
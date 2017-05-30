"use strict";

const mess = document.querySelector('#message')

if (!window["WebSocket"]) {
  mess.textContent = "Sorry, your browser does not support this experiment."
} else {
  let socket = new WebSocket("ws://" + document.location.host + "/ws");

  socket.addEventListener('open', e => {
    mess.textContent = "Connected!";
  });
  socket.addEventListener('error', e => {
    mess.textContent = "Oops, something went wrong";
  });
  socket.addEventListener('close', e => {
    mess.textContent = "Lost connection, trying to reconnect...";
  });
  socket.addEventListener('message', e => {
    document.body.style.backgroundColor = e.data;
  });
}
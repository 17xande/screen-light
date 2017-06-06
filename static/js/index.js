"use strict";

const mess = document.querySelector('#message');
const btnConn = document.querySelector('#btnConnect');
const divBG = document.querySelector('#background');
let socket;

if (!window["WebSocket"]) {
  mess.textContent = "Sorry, your browser does not support this experiment."
} else {
  btnConn.addEventListener('click', socketConnect);
}

function socketConnect(e) {
  const sLink = "ws://" + document.location.host + "/ws"
  socket = new WebSocket(sLink);

  socket.addEventListener('open', e => {
    btnConn.style.display = "none";
    divBG.style.filter = "none";
  });
  socket.addEventListener('error', e => {
    mess.textContent = "Oops, something went wrong";
  });
  socket.addEventListener('close', e => {
    btnConn.innerHTML = "Reconnect...";
    divBG.style.filter = "grayscale(1)";
    socket = null;
  });
  socket.addEventListener('message', e => {
    divBG.style.backgroundColor = e.data;
  });
}
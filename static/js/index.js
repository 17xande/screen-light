"use strict";

const mess = document.querySelector('#message');
const btnConn = document.querySelector('#btnConnect');
const divBG = document.querySelector('#background');

if (!window["WebSocket"]) {
  mess.textContent = "Sorry, your browser does not support this experiment."
} else {
  btnConn.addEventListener('click', socketConnect);
}

function socketConnect(e) {
  const sLink = "ws://" + document.location.host + "/ws"
  let socket = new WebSocket(sLink);

  socket.addEventListener('open', e => {
    btnConn.style.display = "none";
    divBG.style.filter = "none";
  });
  socket.addEventListener('error', e => {
    mess.textContent = "Oops, something went wrong";
  });
  socket.addEventListener('close', e => {
    mess.textContent = "Lost connection, trying to reconnect...";
    setTimeout(() => socket = new WebSocket(sLink), 2000);
  });
  socket.addEventListener('message', e => {
    divBG.style.backgroundColor = e.data;
  });
}
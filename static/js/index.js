"use strict";

const mess = document.querySelector('#message');
const btnConn = document.querySelector('#btnConnect');
const divBG = document.querySelector('#background');
let socket;
let presets;

loadPresets();

if (!window["WebSocket"]) {
  mess.textContent = "Sorry, your browser does not support this experiment."
} else {
  btnConn.addEventListener('click', socketConnect);
}

function loadPresets() {
  fetch("/static/js/presets.json")
    .then(res => res.json())
    .then(jsPre => presets = jsPre.presets)
    .catch(err => console.error(err));
}

function socketConnect(e) {
  const sLink = "ws://" + document.location.host + "/ws"
  socket = new WebSocket(sLink);

  socket.addEventListener('open', e => {
    btnConn.style.display = 'none';
    // divBG.style.backgroundImage = 'none';
    divBG.style.filter = 'none';
  });
  socket.addEventListener('error', e => {
    console.error(e);
  });
  socket.addEventListener('close', e => {
    btnConn.innerHTML = "Reconnect...";
    btnConn.style.display = "";
    divBG.style.filter = 'grayscale(1)';
    divBG.style.backgroundImage = '';
    socket = null;
  });
  socket.addEventListener('message', e => {
    console.log(e.data);
    processMessage(JSON.parse(e.data));
  });
}

function processMessage(message) {

  if (message.preset) {
    let p = presets[message.preset - 1];
    divBG.style.backgroundImage = p;
    divBG.style.backgroundColor = ''
    return;
  } 
  
  if (message.animation === "strobe") {

  }

  divBG.style.backgroundImage = 'none';
  divBG.style.backgroundColor = message.color;
}
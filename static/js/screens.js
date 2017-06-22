"use strict";

const divBG = document.querySelector('#background');
let socket;
let presets;

loadPresets();

if (!window["WebSocket"]) {
  console.warn("Websockets not supported.");
} else {
  socketConnect();
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
    divBG.backgroundColor = "black"
  });
  socket.addEventListener('error', e => {
    console.error(e);
  });
  socket.addEventListener('close', e => {
    console.warn("retrying connection. Try ", retries);
    setTimeout(socketConnect, 2000);
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
  } else {
    divBG.style.backgroundImage = 'none';
    divBG.style.backgroundColor = message.color;
  }
  
  if (message.animation) {
    divBG.classList.add(message.animation);
    if (message.frequency) {
      divBG.style.animationDuration = message.frequency + "ms";
    }
  } else {
    divBG.classList.remove('strobe');
    divBG.style.animationDuration = '';
  }
}
"use strict";

var mess = document.querySelector('#message');
var btnConn = document.querySelector('#btnConnect');
var divBG = document.querySelector('#background');
var retries = 0;
var socket;
var presets;

loadPresets();

if (!window["WebSocket"]) {
  mess.textContent = "Sorry, your browser does not support this experiment."
} else {
  btnConn.addEventListener('click', startConnect);
}

function startConnect(e) {
  var docEl = window.document.documentElement;
  var requestFullScreen = docEl.webkitRequestFullScreen || docEl.mozRequestFullScreen || docEl.requestFullScreen;

  if (requestFullScreen) {
    requestFullScreen.call(docEl);
  }
  socketConnect(e)
}

function loadPresets() {
  // fetch("/static/js/presets.json")
  //   .then(res => res.json())
  //   .then(jsPre => presets = jsPre.presets)
  //   .catch(err => console.error(err));

  var xhr = new XMLHttpRequest();
  xhr.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200){
      presets = JSON.parse(xhr.responseText).presets;
    } else {
      console.warn(xhr.status);
    }
  }  

  xhr.open('GET', '/static/js/presets.json', true);                  
  xhr.send(null); 
}

function socketConnect(e) {
  var sLink = "ws://" + document.location.host + "/ws"
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
    if (retries++ >= 10) {
      btnConn.innerHTML = "Reconnect";
      btnConn.style.display = "";
      divBG.style.filter = 'grayscale(1)';
      divBG.style.backgroundImage = '';
      socket = null;
      retries = 0;
      return;
    }

    console.log("retrying connection. Try ", retries);
    setTimeout(socketConnect, 2000);
  });
  socket.addEventListener('message', e => {
    console.log(e.data);
    processMessage(JSON.parse(e.data));
  });
}

function processMessage(message) {

  if (message.preset) {
    var p = presets[message.preset - 1];
    divBG.style.backgroundImage = p;
    divBG.style.backgroundColor = ''
  } else {
    divBG.style.backgroundImage = 'none';
    divBG.style.backgroundColor = message.color;
  }
  
  divBG.classList.remove('strobe', 'pulse');
  if (message.animation) {
    divBG.classList.add(message.animation);
    if (message.frequency) {
      divBG.style.animationDuration = message.frequency + "ms";
    }
  } else {
    divBG.style.animationDuration = '';
  }
}
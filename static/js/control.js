"use strict";

const container = document.querySelector('.container');
const tempColour = document.querySelector('#colour');
const divSettings = document.querySelector('#settings');
const divStatus = document.querySelector('#status');
const rngs = document.querySelectorAll('#setRanges input');
const txtColours = document.querySelectorAll('#setTexts input');
const settingsColour = document.querySelector('#setColour');
const btnSettings = document.querySelector('#btnSettings');
const btnDelete = document.querySelector('#btnDelete');
const btnLoad = document.querySelector('#btnLoad');
const btnSave = document.querySelector('#btnSave');
const sConn = document.querySelector('#sConn');

let conn;
let active;
let colours;
let setColourVal = {
  r: 255,
  g: 255,
  b: 255
};

// document.addEventListener('ontouchmove', e => e.preventDefault());

rngs.forEach(range => range.addEventListener('change', setChange));
rngs.forEach(range => range.addEventListener('mousemove', setChange));
rngs.forEach(range => range.addEventListener('touchmove', setChange));
txtColours.forEach(txtC => txtC.addEventListener('change', setChange));
settingsColour.addEventListener('click', e => addColour(settingsColour.style.backgroundColor));
btnSettings.addEventListener('click', toggleSettings);
btnDelete.addEventListener('click', removeElement);
btnLoad.addEventListener('click', loadColours);
btnSave.addEventListener('click', saveColours);

function loadColours(e) {
  fetch("/static/js/colours.json")
    .then(res => res.json())
    .then(jsCol => {
      container.innerHTML = '';
      colours = jsCol.colours;
      colours.forEach(colour => addColour(colour));
  });
}

if (!window["WebSocket"]) {
  container.innerHTML = "<h1>Sorry, your browser does not support this experiment.</h1>"
} else {
  conn = new WebSocket("ws://" + document.location.host + "/ws/control");

  conn.onclose = function(e) {
    console.log("connection closed");
    container.innerHTML = "<h1>Connection Closed.</h1>";
  }

  conn.onmessage = function(e) {
    setTotalConnections(e.data);
  }

  conn.onerror = function(e) {
    console.error(error);
  }

  loadColours();
}

function setChange(e) {
  let base = e.target.dataset.base;
  setColourVal[e.target.dataset.base] = e.target.value;
  changeColour();
}

function changeColour() {
  let colour = `rgb(${setColourVal.r}, ${setColourVal.g}, ${setColourVal.b})`;
  rngs[0].value = txtColours[0].value = setColourVal.r;
  rngs[1].value = txtColours[1].value = setColourVal.g;
  rngs[2].value = txtColours[2].value = setColourVal.b;

  settingsColour.style.backgroundColor = colour;
  if (active != null) active.style.backgroundColor = colour;
}

function addColour(colour) {
  let clone = document.importNode(tempColour.content, true);
  let div = clone.firstElementChild;
  div.style.backgroundColor = colour;
  div.addEventListener("click", sendColour);
  container.appendChild(clone);
}

function removeElement(e) {
  active.parentNode.removeChild(active);
}

function sendColour(e) {
  if (!conn) {
    return false;
  }

  let col = this.style.backgroundColor;
  conn.send(col);
  let [r, g, b] = col.substring(4, col.length-1).split(', ');
  setColourVal.r = r;
  setColourVal.g = g;
  setColourVal.b = b;

  if (active != null) active.classList.remove('active');
  this.classList.add('active');
  active = this;

  changeColour();
}

function saveColours(e) {
  let cols = Array.from(container.querySelectorAll('.colour'))
    .map(c => c.style.backgroundColor);
  let p = JSON.stringify({colours: cols});
  fetch("/api/colours/save", {
    method: "POST",
    body: p
  });

  toggleSettings(e);
}

function toggleSettings(e) {
  if (divSettings.classList.contains('hide')) {
    divSettings.classList.remove('hide');
    divStatus.classList.add('hide');
  } else {
    divSettings.classList.add('hide');
    divStatus.classList.remove('hide');
  }
}

function setTotalConnections(iConns) {
  sConn.innerText = `Connections: ${iConns}`;
}
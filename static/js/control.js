"use strict";

const container = document.querySelector('.container');
const tempColour = document.querySelector('#colour');
const logger = document.querySelector('#log');
let conn;

const colours = [
  "#f00",
  "#0f0",
  "#00f",
  "#000",
  "#fff"
]

colours.forEach(colour => addColour(colour));

if (!window["WebSocket"]) {
  container.innerHTML = "<h1>Sorry, your browser does not support this experiment.</h1>"
} else {
  conn = new WebSocket("ws://" + document.location.host + "/ws");

  conn.onclose = function(e) {
    console.log("connection closed");
    container.innerHTML = "<h1>Connection Closed.</h1>";
  }

  conn.onmessage = function(e) {
    console.log(`message received from host: ${e.data}`);
  }
}

// function appendLog(message) {
//   let p = document.createElement('p');
//   p.innerHTML = message;
//   logger.appendChild(p);
// }

function addColour(colour) {
  let clone = document.importNode(tempColour.content, true);
  let div = clone.firstElementChild;
  div.dataset.colour = colour;
  div.style.backgroundColor = colour;
  div.addEventListener("click", sendColour);
  container.appendChild(clone);
}

function sendColour(e) {
  if (!conn) {
    return false;
  }

  console.log(`sending colour: ${this.dataset.colour}`);
  conn.send(this.dataset.colour);
}
let ws;

function onload() {
  console.log("onload()");

  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);

  console.log(queryString);
  console.log(urlParams.get("roomId"));
  console.log(urlParams.get("roomName"));
}

function send(username, content) {
  console.log("send()");
  ws.send(
    JSON.stringify({
      username: username,
      content: content,
    }),
  );
}

function createAndConnect(roomName, messages) {
  console.log("createAndConnect()");

  try {
    ws?.close();
    ws = new WebSocket(`ws://localhost:3000/ws/create-and-connect/${roomName}`);
    ws.onmessage = onmessage(messages);
  } catch (err) {
    alert(err);
  }
}
function connect(roomId, messages) {
  console.log("connect()");

  try {
    ws?.close();
    ws = new WebSocket(`ws://localhost:3000/ws/connect/${roomId}`);
    ws.onmessage = onmessage(messages);
  } catch (err) {
    alert(err);
  }
}

function onmessage(messages) {
  return (event) => {
    data = JSON.parse(event.data);
    if (data.type == "error") {
      alert("error: " + data.content);
      ws.close();
    }

    messages.push(data);
  };
}

let ws;

function send(username, content) {
	console.log("send()");
	ws.send(
		JSON.stringify({
			username: username,
			content: content,
		}),
	);
}

function createAndConnect(messages) {
	console.log("createAndConnect()");
	const roomName = document.querySelector("#roomName").value;

	try {
		ws?.close();
		ws = new WebSocket(`ws://localhost:3000/ws/create-and-connect/${roomName}`);
		ws.onmessage = onmessage(messages);
	} catch (err) {
		alert(err);
	}
}
function connect(messages) {
	console.log("connect()");
	const roomId = document.querySelector("#roomId").value;

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

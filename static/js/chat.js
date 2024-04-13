window.onbeforeunload = (e) => {
	e.preventDefault();
};

window.onunload = () => {
	document.getElementById("disconnect").click();
};

document.addEventListener("htmx:wsClose", (_) => {
	window.location.assign("/");
});

document.addEventListener("htmx:wsAfterMessage", (_) => {
	const div = document.getElementById("messages");
	div.scrollTop = div.scrollHeight;
});

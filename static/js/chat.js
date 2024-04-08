document.addEventListener("htmx:wsClose", (_) => {
	console.log("after htmx:wsClose");
	window.location.assign("/");
});

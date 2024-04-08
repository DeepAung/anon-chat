window.onload = () => {
	console.log("add event beforeunload");
	window.addEventListener("beforeunload", (e) => {
		e.preventDefault();
		console.log("hello world beforeunload");
		e.returnValue = true;
	});
};

document.addEventListener("htmx:wsClose", (e) => {
	// window.location.assign("/");
});

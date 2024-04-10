document.addEventListener("htmx:wsOpen", (_) => {
	console.log("on htmx:wsOpen");
});

// document.onvisibilitychange = () => {
// 	document.getElementById("disconnect-form").submit();
// 	alert("onvisibilitychange");
// };

document.addEventListener("htmx:wsClose", (_) => {
	console.log("on htmx:wsClose");
	window.location.assign("/");
});

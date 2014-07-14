var websocket;
jQuery(document).ready(function() {
	websocket = new WebSocket("ws://localhost:8080/ws");

	console.log(websocket);
});
function generateUUID(){
	var d = new Date().getTime();
	var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
		var r = (d + Math.random()*16)%16 | 0;
		d = Math.floor(d/16);
		return (c=='x' ? r : (r&0x7|0x8)).toString(16);
	});
	return uuid;
}

var websocket;
var uuid;

function onopenHandler() {
	websocket.send(uuid);
}
function oncloseHandler() {

}
function onmessageHandler(event) {
	$('body').append('<p>' + event.data + '</p>');
}

jQuery(document).ready(function() {
	uuid = generateUUID();
	$('h1').html(uuid);
	websocket = new WebSocket("ws://localhost:8080/ws");
	websocket.onopen = onopenHandler;
	websocket.onclose = oncloseHandler;
	websocket.onmessage = onmessageHandler;
});
MessageTypeLogin = 1;
MessageTypeJoin  = 2;
MessageTypeLeave = 3;
MessageTypeText  = 4;

MessageTypeSynchMembers  = 101;

var members = {};

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
	var msg = {
		MessageType: MessageTypeLogin,
		Data: JSON.stringify({
			UUID: uuid
		})
	}
	websocket.send(JSON.stringify(msg))
	//websocket.send(uuid);
}
function oncloseHandler() {
	ShowMessage("server is down", "error");
}
function onmessageHandler(event) {
	var data = JSON.parse(event.data);
	if (data.Type == MessageTypeJoin) {
		ShowMessage(data.Uuid + " joined!", "join");
		NewMember(data.Uuid)
	}
	else if (data.Type == MessageTypeLeave) {
		if (data.Uuid in members) {
			members[data.Uuid].anchor.remove();	
			delete(members[data.Uuid]);
		}

		ShowMessage(data.Uuid + ' leaved!', 'leave');
	}
	else if (data.Type == MessageTypeSynchMembers) {
		ShowMessage('Synchronizing members...');
		for (var i = 0; i < data.Members.length; i++) {
			NewMember(data.Members[i]);
		}
	}
	else if (data.Type == MessageTypeText) {
		if (data.Uuid == uuid) {
			ShowMessage(" me: \"" + data.Text + "\"", "me");
		}
		else {
			ShowMessage(data.Uuid + " says: \"" + data.Text + "\"");
		}
	}
	else {
		ShowMessage(event.data);
	}
}

function NewMember(uuid) {
	if (!(uuid in members)) {
		var member = {
			uuid: uuid,
			anchor: $('<div class="member">'+uuid+'</div>'),
		};
		$('.chat_members').append(member.anchor);
		members[uuid] = member;

		return member;
	}
	else {
		return members[uuid];
	}
}

function ShowMessage(text, messageType) {
	if (typeof messageType == 'undefined' || messageType == null) {
		messageType = "";
	}

	$('.chat_window').append('<div class="message ' + messageType + '">' + text + '</div>');
}

function SendMessage(type, data) {
	var msg = {
		MessageType: type,
		Data: JSON.stringify(data)
	}
	websocket.send(JSON.stringify(msg))
}

jQuery(document).ready(function() {
	uuid = generateUUID();
	$('h1').html(uuid);
	websocket = new WebSocket("ws://" + window.location.host + "/ws");
	websocket.onopen = onopenHandler;
	websocket.onclose = oncloseHandler;
	websocket.onmessage = onmessageHandler;

	$('#chat_form').submit(function(event) {
		event.preventDefault();

		SendMessage(MessageTypeText, {Text: $('.chat_input').val()});
		$('.chat_input').val('');
	})
});

$(document).ready(function() {
	$('select').material_select();
});

$("#start").click(function() {
	var roomsSelectContainer = $(".rooms_select_container").addClass("hidden");
	var roomPreloader = $(".room_preloader").removeClass("hidden");
	var roomContainer = $(".room_container");
	
	var socketUrl = "ws://192.168.20.83:12321/ws?name=" + $("#name").val() 
		+ "&instrument=" + $("#instrument").val()
		+ "&band=" + $("#band").val();
	console.log(socketUrl);
	socket = new WebSocket(socketUrl);
	socket.onopen = function() { 
	    console.log("Соединение установлено."); 
		setTimeout(function(){
			roomPreloader.addClass("hidden");
			roomContainer.filter(".instrument_" + $("#instrument").val()).removeClass("hidden");
			initRoom(roomContainer, socket);
		}, 3000); //just for fun
	};
	
	socket.onclose = function(event) { 
		if (event.wasClean) {
		    console.log('Соединение закрыто чисто');
		} else {
		    console.log('Обрыв соединения'); // например, "убит" процесс сервера
		}
		console.log('Код: ' + event.code + ' причина: ' + event.reason);
	};
	 
	socket.onmessage = function(e) {
	    console.log("Получены данные " + e.data);
		var data = e.data;
		console.log(data);
		var event = $.parseJSON(data);
		console.log(event);
		if (event.event == "play")
		{
			var note = parseInt(event.action);
			var audioPath = "audio/" + event.audio;
			if (audioPath)
			{
				var audio = new Audio(audioPath);
				audio.play();
			}
		}
	};
	
	socket.onerror = function(error) { 
	    console.log("Ошибка " + error.message); 
	};
});

function initRoom (roomContainer, socket) {
	roomContainer.find(".btn-floating").click(function(){
		var action = $(this).attr("id").replace(/[^\d]+/, '');
		socket.send(action);
	});
	
	$("body").keydown(function(event){
		switch (event.keyCode) {
			case 49: 
			case 97: 
				socket.send(1);	
				break;
			case 50: 
			case 98: 
				socket.send(2);	
				break;
			case 51: 
			case 99: 
				socket.send(3);	
				break;	
			case 52: 
			case 100: 
				socket.send(4);	
				break;	
			case 53: 
			case 101: 
				socket.send(5);	
				break;
			case 54: 
			case 102: 
				socket.send(6);	
				break;					
		}
	})
}
function chat(roomName){

	var socket = new WebSocket('ws://localhost:8080/chat/global/ws');
	$("#messages").empty();

	var display = function(event) {
	    $('#messages').append(template({event: event}));
	    $('#messages').scrollTo('max')
	  }

	var template = function(event) {
		$('messages').innerText(event)
	}
	   // Message received on the socket
	  socket.onmessage = function(event) {
	    display(JSON.parse(event.data))
	  }

	  $('#send').click(function(e) {
	    var message = $('#message').val();
	    $('#message').val('');
	    socket.send(message)
	  });

	  $('#message').keypress(function(e) {
	    if(e.charCode == 13 || e.keyCode == 13) {
	      $('#send').click()
	      e.preventDefault()
	    }
	  });


}

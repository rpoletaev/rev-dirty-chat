$(function(){
    var ws = new WebSocket("ws://" + location.host + "/ws");
    ws.onmessage = function(event){
        var data = JSON.parse(event.data);
        if (data.event == "message"){
            // $('#chat_menu_item').append('<span class="new badge">&nbsp</span>');
            console.log(data);
        }
    }
});
// 
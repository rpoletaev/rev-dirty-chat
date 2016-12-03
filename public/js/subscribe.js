$(function(){
    var ws = new WebSocket("ws://" + location.host + "/ws");
    ws.onmessage = function(event){
        var data = JSON.parse(event.data);
        if (data.event == "message"){
            // $('#chat_menu_item').append('<span class="new badge">&nbsp</span>');
            console.log(data);
            var $toastTitile = '<a href="/user/' + data.data.User.Name +'">'+ data.data.User.Name + '</a>';
            var $toastImage = '<img class="circle responsive-img" src="'+data.data.User.Avatar+'" width=30px/>';
            var $toastText = '<p>'+ data.data.Text +'</p>';
            var $toastContent = '<div class="row valign-wrapper">';
            $toastContent = $toastContent + '<div class="col s2">' + $toastImage + '</div><div class="col s10">' + $toastTitile + $toastText + '</div></div>'
            console.log($toastContent); 
            Materialize.toast($toastContent, 3000);
        }
    }
});
// 
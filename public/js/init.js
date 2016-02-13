$('.datepicker').pickadate({
    selectMonths: true, 
    selectYears: 50, 
    clear: 'Очистить',
    close: 'Закрыть',
    today: '',
    weekdaysShort: ['Пн', 'Вт', 'Ср', 'Чет', 'Пт', 'Суб', 'Вс'],
  	formatSubmit: 'dd-mm-yyyy',
  	min: new Date(1940, 1, 1),
  	closeOnSelect: true
  });

$('input.file-path').on('change', function(){
	var file_data = $('input#avatar')[0].files[0];
	var form_data = new FormData();
	var token = $("input[name='csrf_token']")[0].value;
	form_data.append('avatar', file_data);
	form_data.append('csrf_token', token);
	$.ajax({
		url: 'mainimageupload',
		dataType: 'json',
		cache: false,
		contentType: false,
		processData: false,
		data: form_data,
		type: 'post',
		success: function(data) {
			$('img#portrait').attr('src', data.Big);
		}
	});
});

$('.tool-item').each(function(){
	var btn = $(this).find('a.btn-floating');
	btn.hover(
		function(){
			btn.removeClass('transparent').addClass('grey');
		},

		function(){
			btn.removeClass('grey').addClass('transparent');
		}		
	);
});

//Initialization
$(function(){
	$(".button-collapse").sideNav();
	$('.materialboxed').materialbox();
	// using croppie plugin https://github.com/Foliotek/Croppie
	$hid = $('#hidden-image').croppie({
		viewport: {
			width: 100,
			height: 100
		},
		boundary: {
			width: 300,
			height: 300
		}
	});

	$('.modal-trigger').leanModal({
		complete: function(){
			$hid.croppie('result', {
				type: 'canvas',
				size: 'viewport'
			}).then(function (resp) {
				var form_data = new FormData();
				var token = $("input[name='csrf_token']")[0].value;
				form_data.append("avatar", dataURItoBlob(resp));
				form_data.append('csrf_token', token);
				$.ajax({
						url: 'avatarupload',
						dataType: 'json',
						cache: false,
						contentType: false,
						processData: false,
						data: form_data,
						type: 'post',
						success: function(data) {
							alert("Аватар загружен");
						}
					});
			});
		}
	});

	$hid.croppie('bind', $('#portrait').attr('src'));

	$('#image-edit').on('click', function () {
					$hid.toggle();
					$hid.croppie('bind');
				});	
	//End of Using Croppie

	//DataBinding
	$("input.bindable").not(":text").each(function(){
		$(this).on('change', function(){
			var val='';
			if($(this).attr("type")=='checkbox') {
				val = $(this).prop("checked");
			}else{
				val = $(this).val();
			}

			$(this).prop('disabled', true);

			var token = $("input[name='csrf_token']")[0].value;
			var user = new FormData();
			user.append("name", $(this).attr("name"));
			user.append("val", val);
			user.append("csrf_token", token);
			$.ajax({
				url: 'update',
				dataType: 'json',
				cache: false,
				contentType: false,
				processData: false,
				data: user,
				type: 'post',
				success: function(data) {
					if(data.Error){
						alert(data.Error);
					}
				}
			});

			$(this).prop('disabled', false);
		});
	});

	$('textarea.bindable').each(function(){
		var btn = $('<a></a>').addClass('waves-effect waves-light btn-flat grey lighten-4').text("Сохранить");
		var name = $(this).attr('name');
		
		btn.click(function(){
			var area = $(this).parent().children('.bindable');
			var val = area[0].value;
			var token = $("input[name='csrf_token']")[0].value;
			var user = new FormData();
			user.append("name", name);
			user.append("val", val);
			user.append("csrf_token", token);
			$.ajax({
				url: 'update',
				dataType: 'json',
				cache: false,
				contentType: false,
				processData: false,
				data: user,
				type: 'post',
				success: function(data) {
					if(data.Error){
						alert(data.Error);
						Materialize.toast("OK!", 4000);
					}
				}
			});
			btn.remove();
		});

		$(this).focus(function(){
			$(this).parent().addClass('hoverable');
			$(this).parent().append(btn);
		});
		$(this).focusout(function(){
			$(this).parent().removeClass('hoverable');
		});
	});

	$(':text.bindable').each(function(){
		var btn = $('<a></a>').addClass('waves-effect waves-light btn-flat grey lighten-4').text("Сохранить");
		var name = $(this).attr('name');
		
		btn.click(function(){
			var area = $(this).parent().children('.bindable');
			var val = area[0].value;
			var token = $("input[name='csrf_token']")[0].value;
			var user = new FormData();
			user.append("name", name);
			user.append("val", val);
			user.append("csrf_token", token);
			$.ajax({
				url: 'update',
				dataType: 'json',
				cache: false,
				contentType: false,
				processData: false,
				data: user,
				type: 'post',
				success: function(data) {
					if(data.Error){
						alert(data.Error);
						Materialize.toast("OK!", 4000);
					}
				}
			});
			btn.remove();
		});

		$(this).focus(function(){
			$(this).parent().addClass('hoverable');
			$(this).parent().append(btn);
		});
		$(this).focusout(function(){
			$(this).parent().removeClass('hoverable');
		});
	});
})


function dataURItoBlob(dataURI) {
        var split = dataURI.split(','),
            dataTYPE = split[0].match(/:(.*?);/)[1],
            binary = atob(split[1]),
            array = []
        for(var i = 0; i < binary.length; i++) array.push(binary.charCodeAt(i))
        return new Blob([new Uint8Array(array)], {
            type: dataTYPE
        })
    }

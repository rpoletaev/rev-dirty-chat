$('.datepicker').pickadate({
    selectMonths: true, // Creates a dropdown to control month
    selectYears: 15 // Creates a dropdown of 15 years to control year
  });

// $('.drop-link').click(function(){
// 	var $btn = $(this);
// 	var strings = $(this).prev().attr('href').split("/");
// 	console.log(strings);

// 	var name = "/" + strings[1] + "/" + strings[2]
// 	$.ajax({
// 		url: name,
// 		type: 'DELETE',
// 		timeout: 30000,
// 		success: function(result){
// 			$btn.closest('li').remove();
// 		}
// 	});
// });

$('input.file-path').on('change', function(){
	var file_data = $('input#avatar')[0].files[0];
	var form_data = new FormData();
	var token = $("input[name='csrf_token']")[0].value;
	form_data.append('avatar', file_data);
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
			$('img#portrait').attr('src', data.Big);
		}
	});
});
// try using croppie plugin https://github.com/Foliotek/Croppie
$('input.file-path').on('change', function(){
	var $uploadCrop;

	var file_data = $('input#avatar')[0].files[0];

	function readFile(input) {
		if (input.files && input.files[0]) {
			var reader new FileReader();

			reader.onload = function(e){
				$uploadCrop.croppie('bind', {
					url: e.target.result
				});

				
			}
		}
	}
	var form_data = new FormData();
	var token = $("input[name='csrf_token']")[0].value;
	form_data.append('avatar', file_data);
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
			$('img#portrait').attr('src', data.Big);
		}
	});
});

$('#crop-container').croppie({
	viewport: {
		width: 200, 
		height: 200
	}
});
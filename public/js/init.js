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
	//var upload_path = '/user/' + $('#login')[0].text + '/avatarupload'
	$.ajax({
		url: 'avatarupload',
		dataType: 'text',
		cache: false,
		contentType: false,
		processData: false,
		data: form_data,
		type: 'post',
		success: function(data) {
			alert(data);
		}
	});
});
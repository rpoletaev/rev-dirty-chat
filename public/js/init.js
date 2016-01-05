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


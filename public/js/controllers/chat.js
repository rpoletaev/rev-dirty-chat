angular.module('chat', ['ngWebsocket'])
.config(function($interpolateProvider, $locationProvider) {
    $locationProvider.html5Mode({enabled: true, requireBase: false}).hashPrefix('!');
  })

.run(['$anchorScroll', function($anchorScroll){
	$anchorScroll.yOffset = 50;
}])

.controller('ChatMessages', ['$scope', '$filter', '$anchorScroll', '$location', '$http', '$websocket', function($scope, $filter, $anchorScroll, $location, $http, $websocket){
	$scope.messages = [];
	$scope.newMessage = "";
	$scope.msgCount = 0;

	$scope.ws = $websocket.$new('ws://' + location.host + location.pathname + '/ws');
	$scope.ws.$on('$message', function(data) {
			//console.log(data);
			if (data.event == 'message') {
				$scope.addMessage(data.data);		
			}
			
			setTimeout($scope.gotoBottom(), 0);
		});

	$scope.addMessage = function (message) {
		message.Datestr = $filter('date')(new Date(message.Timestamp*1000), 'dd.MM.yyyy');
		message.hash = $scope.msgCount;

		 if ($scope.msgCount > 0){
		 	var lastMessage = $scope.messages[$scope.msgCount - 1];

		 	if (lastMessage.User.OriginalID == message.User.OriginalID)
		 	{
		 		$scope.messages[$scope.msgCount - 1].Strings.push(message.Text);
		 		$scope.messages[$scope.msgCount - 1].Datestr = message.Datestr;
		 	}else{
				message.Strings = [message.Text];
		 		$scope.messages.push(message);
		 	}
		 }else{
			message.Strings = [message.Text];
		 	$scope.messages.push(message);
		 }
		 
		$scope.$apply();
	};

	$scope.$watchCollection('messages', function(newMsgs, oldMsges){
		$scope.msgCount = newMsgs.length;
	});

	$scope.send = function() {
		if ($scope.newMessage != "") {
			$scope.ws.$emit('message', $scope.newMessage);
			$scope.newMessage = "";
		}
	};

	$scope.gotoBottom = function() {
    	$location.hash('anchor' + $scope.msgCount);
    	$anchorScroll();
    };
}])

.directive('enterSubmit', function () {
    return {
      restrict: 'A',
      link: function (scope, elem, attrs) {
       
        elem.bind('keydown', function(event) {
          var code = event.keyCode || event.which;
                  
          if (code === 13) {
            if (event.ctrlKey) {
              event.preventDefault();
              scope.$apply(attrs.enterSubmit);
            }
          }
        });
      }
    }
  });

	
angular.module('chat', ['ngWebsocket'])
.config(function($interpolateProvider, $locationProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
    $locationProvider.html5Mode({enabled: true, requireBase: false}).hashPrefix('!');
  })

.run(['$anchorScroll', function($anchorScroll){
	$anchorScroll.yOffset = 50;
}])

.controller('ChatMessages', ['$scope', '$filter', '$anchorScroll', '$location', '$http', '$websocket', function($scope, $filter, $anchorScroll, $location, $http, $websocket){
	$scope.messages = [];
	$scope.newMessage = "";
	$scope.msgCount = 0;

	$scope.ws = $websocket.$new('ws://' + location.host() + '/' + $location.path() + '/ws');
	$scope.ws.$on('$message', function(data) {
			console.log(data);
			if (data.event == 'message') {
				$scope.addMessage(data.data);
				$scope.gotoBottom();		
			}
		});

	$scope.addMessage = function (message) {
		message.Datestr = $filter('date')(new Date(message.Timestamp*1000), 'dd.MM.yyyy');
		// if ($scope.messages.length() > 0 && $scope.messages[messages.length() - 1].User.OriginalID == message.User.OriginalID){
		// 	$scope.messages[messages.length() - 1].Text + '\n' + message.Text;
		// 	$scope.messages[messages.length() - 1].Datestr = message.Datestr;
		// }else{
			message.hash = $scope.msgCount;
			$scope.messages.push(message);
		// }
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
    	$location.hash('bottom');
    	$anchorScroll();
    };
}]);
	
	// .controller('Rooms', ['$scope', '$http', function($scope, $http){
// 	$scope.rooms = [];
// 	$http.get('/chat/myrooms').success(function(data){
// 		console.log(data);
// 		$scope.rooms = data;
// 	});


// }])
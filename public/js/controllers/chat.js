angular.module('chat', ['ngWebsocket'])
.config(function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
  })

.run(['$anchorScroll', function($anchorScroll){
	$anchorScroll.yOffset = 50;
}])

.controller('ChatMessages', ['$scope', '$websocket', '$filter', '$anchorScroll', '$location', function($scope, $websocket, $filter, $anchorScroll, $location){
	var ws = $websocket.$new('ws://localhost:9000/chat/socket?user=roma');
	ws.$on('$open', function () {
		console.log("connection established");
	});

	ws.$on('$close', function () {
    	console.log('Noooooooooou, I want to have more fun with ngWebsocket, damn it!');
    });

	$scope.messages = [];
	$scope.newMessage = "";
	$scope.msgCount = 0;

	$scope.addMessage = function (message) {
		message.Datestr = $filter('date')(new Date(message.Timestamp*1000), 'dd.MM.yyyy');
		message.hash = $scope.msgCount;
		$scope.messages.push(message);
		console.log("last hash" + message.hash);
		$scope.$apply();
	};

	$scope.$watchCollection('messages', function(newMsgs, oldMsges){
		$scope.msgCount = newMsgs.length;
	});

	$scope.send = function() {
		if ($scope.newMessage != "") {
			ws.$emit('message', $scope.newMessage);
			$scope.newMessage = "";
		}
	};

	ws.$on('$message', function(event) {
		if (event.event == 'message') {
			$scope.addMessage(event.data);
			$scope.gotoBottom();
		}
	});		

    $scope.gotoBottom = function() {
    	$location.hash('bottom');
    	$anchorScroll();
    }
}]);
	
var chat = angular.module('chat', []);
chat.factory('ChatService', function() {
var ws = new WebSocket("ws://localhost:9000/chat/socket?user=roma");
var service = {};
ws.onopen = function() {
	//console.log("Succeeded to open a connection");
	service.callback("Succeeded to open a connection");
};
 
 ws.onerror = function() {
 service.callback("Failed to open a connection");
 }
 
 ws.onmessage = function(message) {
 	service.callback(message.data);
 };
 
 service.ws = ws;
  
 service.send = function(message) {
 	//console.log(message);
 	service.ws.send(message);
 }
 
 service.subscribe = function(callback) {
 	service.callback = callback;
 }
 
 return service;
});
 
 
chat.controller('ChatMessages', ['$scope', 'ChatService', function($scope, ChatService) {
 $scope.messages = [];
 
 ChatService.subscribe(function(message) {
 	console.log("from Subscribe");
 	console.log(message.event);
 	console.log(message.data);
 	$scope.messages.push({User: "roma", Text: "xy"});
	$scope.$apply();
 });
 
 // $scope.connect = function() {
	// chatService.connect();
 // }
 
 $scope.send = function() {
	 ChatService.send($scope.message);
	 $scope.message = "";
 }
}]);
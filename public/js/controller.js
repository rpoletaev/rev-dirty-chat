var commonObjects = angular.module('commonObjects', []);

commonObjects.controller('CommonObjectList', ['$scope', '$http', function($scope, $http){
	$scope.path = window.location.pathname;
	$scope.typeName = "";
	$scope.objects = [];
	$scope.typeName = "Sexes"

	$http.get($scope.path + '.json').success(function(data) {
		$scope.objects = data.Collection; 
	});

	$scope.delete = function(name) {
		$http.delete($scope.path + '/' + name);
		for(i=0; i < $scope.objects.length; i++) {
			if ($scope.objects[i].Name == name){
				$scope.objects.splice(i, 1);
			}
		} 
	};

}]);

	

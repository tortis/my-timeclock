var mtcApp = angular.module('mtcApp', []);

mtcApp.controller('mtcCtrl', function($scope, $http) {
	$scope.status = false;
	$scope.loadWeek = function(sunday) {
		$http.get('/week', {params: {sunday: sunday}})
		.success(function(data, status, headers, config) {
			$scope.week = data;
		})
		.error(function(err, status, headers, config) {
			console.log("week request failed for: " + sunday);
			$scope.error = "Well this is embarrassing.";
		});
	}

	$scope.getStatus = function() {
		$http.get('/status')
		.success(function(data, status, headers, config) {
			if (data == "true") {
				$scope.status = true;
			} else {
				$scope.status = false;
			}
		})
		.error(function(err, status, headers, config) {
			console.log("status request failed.");
		});
	}

	$scope.clockIn = function() {
		$http.get('/in')
		.success(function(data, status, headers, config) {
			$scope.status = true;
			$scope.loadWeek(1418536800);
		})
		.error(function(err, status, headers, config) {
			console.log("clock in request failed.");
			console.log(status);
			console.log(err);
		});
	}

	$scope.clockOut = function() {
		$http.get('/out')
		.success(function(data, status, headers, config) {
			$scope.status = false;
			$scope.loadWeek($scope.getSunday());
		})
		.error(function(err, status, headers, config) {
			console.log("clock out request failed.");
		});
	}

	$scope.deleteShift = function(shift) {
		$http.get('/deleteshift', {params: {id: shift.Id}})
		.success(function(date, status, headers, config) {
			$scope.loadWeek($scope.getSunday());
			$scope.getStatus();
		})
		.error(function(err, status, headers, config) {
			console.log("Failed to delete shift: " + shift);
		});
	}

	$scope.createShift = function(unix_in, unix_out) {

	}

	$scope.startEdit = function(shift) {
		shift.non = $scope.dateToTime(shift.On*1000);
		shift.noff = $scope.dateToTime(shift.Off*1000);
		shift.Editing = true;
	}

	$scope.modifyShift = function(shift) {
		console.log("New on: " + shift.non);
		console.log("New off: " + shift.noff);
	}

	$scope.dateToTime = function(time) {
		var d = new Date(time);
		return d.toLocaleTimeString();
	}

	$scope.dateToDate = function(date) {
		var d = new Date(date);
		return d.toLocaleDateString();
	}

	$scope.getSunday = function() {
		return Date.today().previous().sunday().getTime()/1000;
	}

	$scope.shiftHours = function(shift) {
		if (shift.Active) {
			var now = new Date();
			return ( now.getTime()/1000 - shift.On)/3600.0;
		} else {
			return (shift.Off - shift.On)/3600.0;
		}
	}

	$scope.getStatus();
	$scope.loadWeek($scope.getSunday());
});

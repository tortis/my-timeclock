var mtcApp = angular.module('mtcApp', []);

mtcApp.directive('ngEnter', function () {
	return function (scope, element, attrs) {
		element.bind("keydown keypress", function (event) {
			if(event.which === 13) {
				scope.$apply(function (){
					scope.$eval(attrs.ngEnter);
				});

				event.preventDefault();
			}
		});
	};
});

mtcApp.controller('mtcCtrl', function($scope, $http) {
	$scope.clock = false;
	$scope.dateShowing = moment().startOf('week');
	$scope.current = true;
	$scope.day_names = [];
	var m = moment();
	for (var i = 0; i < 7; i++) {
		m.weekday(i)
		$scope.day_names.push(m.format("ddd"));
	}
	$scope.loadWeek = function() {
		$http.get('/week', {params: {sunday: $scope.dateShowing.unix()}})
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
				$scope.clock = true;
			} else {
				$scope.clock = false;
			}
		})
		.error(function(err, status, headers, config) {
			console.log("status request failed.");
		});
	}

	$scope.clockIn = function() {
		$http.get('/in')
		.success(function(data, status, headers, config) {
			$scope.clock = true;
			$scope.loadWeek();
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
			$scope.clock = false;
			$scope.loadWeek();
		})
		.error(function(err, status, headers, config) {
			console.log("clock out request failed.");
		});
	}

	$scope.deleteShift = function(shift) {
		$http.get('/deleteshift', {params: {id: shift.Id}})
		.success(function(data, status, headers, config) {
			$scope.loadWeek();
			$scope.getStatus();
		})
		.error(function(err, status, headers, config) {
			console.log("Failed to delete shift: " + shift);
		});
	}

	$scope.createShift = function(day, on, off) {
		$http.get('/createshift', {params: {on: on.getTime()/1000, off: off.getTime()/1000}})
		.success(function(data, status, headers, config) {
			$scope.loadWeek();
			day.Adding = false;
		})
		.error(function(err, status, headers, config) {
			console.log("Failed to create new shift: " + err);
			day.Adding = false;
		});
	}

	$scope.startAdding = function(day) {
		day.new_on = new Date(day.date);
		day.new_on.setHours(9);
		day.new_on.setMinutes(0);
		day.new_on.setSeconds(0, 0);
		day.new_off = new Date(day.date);
		day.new_off.setHours(17);
		day.new_off.setMinutes(0);
		day.new_off.setSeconds(0, 0);
		day.Adding = true;
	}

	$scope.startEdit = function(shift) {
		shift.non = new Date(shift.On * 1000);
		shift.noff = new Date(shift.Off * 1000);
		shift.Editing = true;
	}

	$scope.modifyShift = function(shift, non, noff) {
		if (!non || !noff) {
			shift.Editing = false;
			shift.Error = true;
			setTimeout(function(){$scope.$apply(function(){shift.Error = false;})}, 250);
			return;
		}

		var on = moment.unix(shift.On);
		var off = moment.unix(shift.Off);
		var new_on = moment(non);
		var new_off = moment(noff);
		on.hour(new_on.hour()).minute(new_on.minute());
		off.hour(new_off.hour()).minute(new_off.minute());
		
		if (on.unix() == shift.On && off.unix() == shift.Off) {
			shift.Editing = false;
			return;
		}
		var params = {params: {id: shift.Id}};
		if (off.unix() != shift.Off) {
			params.params.off = off.unix();
		}
		if (on.unix() != shift.On) {
			params.params.on = on.unix();
		}
		$http.get('/editshift', params)
		.success(function(data, status, headers, config) {
			$scope.loadWeek();
		})
		.error(function(err, status, headers, config) {
			console.log("Failed to modify shift: " + err);
			shift.Error = true;
			setTimeout(function(){$scope.$apply(function(){shift.Error = false;})}, 500);
			shift.Editing = false;
		});
	}

	$scope.previousWeek = function() {
		$scope.dateShowing = $scope.dateShowing.subtract(1, 'w');
		$scope.loadWeek();
		if (moment().startOf('week').unix() == $scope.dateShowing.unix()) {
			$scope.current = true;
		} else {
			$scope.current = false;
		}
	}
	
	$scope.nextWeek = function() {
		$scope.dateShowing = $scope.dateShowing.add(1, 'w');
		$scope.loadWeek();
		if (moment().startOf('week').unix() == $scope.dateShowing.unix()) {
			$scope.current = true;
		} else {
			$scope.current = false;
		}
	}

	$scope.currentWeek = function() {
		$scope.dateShowing = moment().startOf('week');
		$scope.loadWeek();
		$scope.current = true;
	}

	$scope.parseTime = function(timeString, dt) {
		if (!dt) {
			dt = moment();
		}

		var time = timeString.match(/(\d+)(?::(\d\d))?\s*(p?)/i);
		if (!time) {
			return null;
		}
		var hours = parseInt(time[1], 10);
		if (hours == 12 && !time[3]) {
			hours = 0;
		} else {
			hours += (hours < 12 && time[3]) ? 12 : 0;
		}
		dt.hour(hours);
		dt.minutes(parseInt(time[2], 10) || 0);
		return dt;
	}

	$scope.dateToTime = function(time) {
		return	moment.unix(time).format("h:mm A");
	}

	$scope.dateToDate = function(date) {
		return moment(date).format("MM/D/YYYY");
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
	$scope.loadWeek();
});

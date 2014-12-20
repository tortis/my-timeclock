var mtcApp = angular.module('mtcApp', []);

mtcApp.controller('mtcCtrl', function($scope, $http) {
    $http.get('/week', {sunday: 1418515200})
    .success(function(data, status, headers, config) {
        console.log(data);
    })
    .error(function(err, status, headers, config) {
        
    });
    
    $scope.week = { days: [
       {Name: "Sunday", Hours: 0.0},
       {Name: "Monday", Date: "12/15/2014", Hours: 0.0, Shifts: [
           {In: "9:05 AM", Out: "5:30 PM"},
           {In: "", Out: "asd"}]},
       {Name: "Tuesday", Hours: 1.8},
       {Name: "Wednesday", Hours: 5.6},
       {Name: "Thursday", Hours: 1.6},
       {Name: "Friday", Hours: 5.7},
       {Name: "Saturday", Hours: 1.8}
   ]};
});
var IN_URL = "/in"
var IN_URL = "/status"
var STATUS_URL = "/status"
var SHIFTS_URL = "/shifts"
var WEEK_URL = "/week" 

var make_week_view = doT.template('<div class="panel panel-default"><div class="panel-heading"><span class="weekname"><strong>{{=it.name}}</strong></span><span class="caps">{{=it.date}}</span><span class="caps pull-right">{{=it.total_hours.toFixed(1)}}</span></div><table class="table"><tr><td>Sun</td><td>Mon</td><td>Tues</td><td>Wed</td><td>Thur</td><td>Fri</td><td>Sat</td></tr><tr><td>{{=it.days[0].T.toFixed(1)}}</td><td>{{=it.days[1].T.toFixed(1)}}</td><td>{{=it.days[2].T.toFixed(1)}}</td><td>{{=it.days[3].T.toFixed(1)}}</td><td>{{=it.days[4].T.toFixed(1)}}</td><td>{{=it.days[5].T.toFixed(1)}}</td><td>{{=it.days[6].T.toFixed(1)}}</td></tr></table></div>');

var make_detail = doT.template('{{~it.days :day:index}}<div class="panel panel-default detail"><div class="panel-heading">{{=day_names[index]}}<span class="caps pull-right">day.Date</span></div><div class="panel-content list-group">{{~day.S :shift:index2}}<div class="list-group-item"><span class="left"><span class="caps">IN</span><span class="">{{=shift.On}}<span></span><span class="pull-right"><span class="caps">OUT</span><span>{{=shift.Off}}</span></span></div>{{~}}</div></div>{{~}}');

function displayWeek(sunday) {
	$.ajax({
		url: WEEK_URL,
		cache: false,
		data: {sunday: sunday}
	}).done(function (response) {
		week = JSON.parse(response);
	});
}

function addWeek(sunday) {
	$.ajax({
		url: WEEK_URL,
		cache: false,
		data: {sunday: sunday}
	}).done(function(response) {
		week = JSON.parse(response);
		var total = 0;
		for (var i = 0; i < week.length; i++) {
			total += week[i].T
		}
		var result = make_week_view({days: week, name: "This Week", date: "12/15/2014", total_hours: total});
		$("#weeks").append(result);
		console.log(result);
		console.log(response);
	});
}

function addShifts(start, end) {

}

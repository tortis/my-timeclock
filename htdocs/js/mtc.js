var onClock = false;
var timeOn = 0;
var timer;

var weekday = new Array(7);
weekday[0]=  "Sunday";
weekday[1] = "Monday";
weekday[2] = "Tuesday";
weekday[3] = "Wednesday";
weekday[4] = "Thursday";
weekday[5] = "Friday";
weekday[6] = "Saturday";

function getStatus(cb) {
	var req = new XMLHttpRequest();
	req.open("GET", "/status/", true);
	req.send();
	req.onload = function () {
		var response = JSON.parse(req.responseText);
		if (response.OnClock) {
			cb(true, response.TimeOn);
		} else {
			cb(false, 0);
		}
	};
}

function clockIn(cb) {
	var req = new XMLHttpRequest();
	req.open("GET", "/in/", true);
	req.send();
	req.onload = function () {
		if (req.responseText == "true") {
			cb(true);
		} else {
			cb(false);
		}
	};
}

function clockOut(cb) {
	var req = new XMLHttpRequest();
	req.open("GET", "/out/", true);
	req.send();
	req.onload = function () {
		if (req.responseText == "true") {
			cb(true);
		} else {
			cb(false);
		}
	};

}

function getWeek(d, cb) {
	var req = new XMLHttpRequest();
	req.open("GET", "/week/?year="+ d.getFullYear()  +"&month="+ (d.getMonth() + 1) +"&day="+ d.getDate(), true);
	req.send();
	req.onload = function () {
		var w = JSON.parse(req.responseText);
		w.Monday = new Date(w.MondayDate);
		cb(w);
	};
}

function weekToHTML(week) {
	var html = '<tr><td>Sun</td><td>Mon</td><td>Tues</td><td>Wed</td><td>Thur</td><td>Fri</td><td>Sat</td></tr><tr><td>'+week.Days[0].Hours.toFixed(1)+'</td><td>'+week.Days[1].Hours.toFixed(1)+'</td><td>'+week.Days[2].Hours.toFixed(1)+'</td><td>'+week.Days[3].Hours.toFixed(1)+'</td><td>'+week.Days[4].Hours.toFixed(1)+'</td><td>'+week.Days[5].Hours.toFixed(1)+'</td><td>'+week.Days[6].Hours.toFixed(1)+'</td></tr>';
	return html;
}

function getBlocks(day) {
	var r = $("<div>").addClass("panel-content").addClass("list-group");
	for (var i = 0; i < day.Blocks.length; i++) {
		var tin = new Date(day.Blocks[i].In);
		var tout = new Date(day.Blocks[i].Out);
		r.append($("<div>").addClass("list-group-item").append(
			$("<span>").addClass("left").append(
				$("<span>").addClass("caps").text("IN ")).append(
				$("<span>").text(tin.toLocaleTimeString()))).append(
			$("<span>").addClass("right").append(
				$("<span>").addClass("caps").text("OUT ")).append(
				$("<span>").text(tout.toLocaleTimeString()))));
	}
	return r;
}

function updateDetails(week) {
	$("#details").empty();
	var dotw = (new Date()).getDay();
	for (;dotw >= 0; dotw--) {
		if (week.Days[dotw].Blocks == null) {
			continue;
		}
		$("#details").append(
			$("<div>").addClass("panel").addClass("panel-default").addClass("detail").append(
				$("<div>").addClass("panel-heading").text(weekday[dotw])).append(getBlocks(week.Days[dotw])));
	}
	
}

function updateTime() {
	timeOn += 0.1;
	$("#timeon").text(timeOn.toFixed(1) + " hrs");
}

function buttonUpdate(s, timeon) {
	if (s) {
		timeOn = timeon;
		timer = setInterval(updateTime, 1000*60*6);
		$("#clk").text("Clock Out").removeClass("btn-success").addClass("btn-danger");
		$("#timeon").text(timeOn.toFixed(1) + " hrs");
	} else {
		timeOn = 0;
		clearInterval(timer);
		$("#clk").text("Clock In").removeClass("btn-danger").addClass("btn-success");
		$("#timeon").empty();
	}
}

function updateStatus() {
	getStatus(function (s, timeon) {
		onClock = s;
		buttonUpdate(s, timeon);
	});
}

function updateWeeks() {
	getWeek(new Date(), function (week) {
		$("#tw").html(weekToHTML(week));
		$("#twhrs").text(week.Hours.toFixed(1) + " Hrs");
		$("#twdate").text(week.Monday.toLocaleDateString());
		updateDetails(week);
	});
	var oneWeekAgo = new Date();
	oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);
	getWeek(oneWeekAgo, function (week) {
		$("#lw").html(weekToHTML(week));
		$("#lwhrs").text(week.Hours.toFixed(1) + " Hrs");
		$("#lwdate").text(week.Monday.toLocaleDateString());
	});

}

window.onload = function () {
	$("#clk").click(function () {
		if (onClock) {
			clockOut(function (s) {
				if (s) {
					onClock = false;
					buttonUpdate(onClock, 0);
					updateWeeks();
				}
			});
		} else {
			clockIn(function (s) {
				if (s) {
					onClock = true;
					buttonUpdate(onClock, 0);
				}

			});
		}
	});
	updateStatus();
	updateWeeks();
}

var onClock = false;
var timeOn = 0;
var timer;

function getWeekNumber(d) {
	d = new Date(+d);
	d.setHours(0,0,0);
	d.setDate(d.getDate() + 4 - (d.getDay()||7));
	var yearStart = new Date(d.getFullYear(),0,1);
	var weekNo = Math.ceil(( ( (d - yearStart) / 86400000) + 1)/7)
	return [d.getFullYear(), weekNo];
}

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
	var ISOWeek = getWeekNumber(d);
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

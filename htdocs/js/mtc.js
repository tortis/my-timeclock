var state = false;

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
		if (req.responseText == "true") {
			cb(true);
		} else {
			cb(false);
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
	req.open("GET", "/week/?year="+ISOWeek[0]+"&week="+ISOWeek[1], true);
	req.send();
	req.onload = function () {
		var w = JSON.parse(req.responseText);
		if (w.Year == ISOWeek[0] && w.WeekNum == ISOWeek[1]) {
			cb(w);
		} else {
			console.log("JSON Parse failed.");
			cb(null);
		}
	};
}

function weekToHTML(week) {
	var html = '<table><tr><td>Monday</td><td>Tuesday</td><td>Wednesday</td><td>Thursday</td><td>Friday</td></tr><tr><td>'+week.Days[1].Hours.toFixed(1)+'</td><td>'+week.Days[2].Hours.toFixed(1)+'</td><td>'+week.Days[3].Hours.toFixed(1)+'</td><td>'+week.Days[4].Hours.toFixed(1)+'</td><td>'+week.Days[5].Hours.toFixed(1)+'</td></tr></table>';
	return html;
}

window.onload = function () {
	$("#clk").click(function () {
		if (state) {
			clockOut(function (s) {
				if (s) {
					state = false;
					$("#clk").text("Clock In");
				}
			});
		} else {
			clockIn(function (s) {
				if (s) {
					state = true;
					$("#clk").text("Clock Out");
				}

			});
		}
	});
	getStatus(function (s) {
		state = s;
		if (s) {
			$("#clk").text("Clock Out");
		} else {
			$("#clk").text("Clock In");
		}
	});
	getWeek(Date.now(), function (week) {
		$("#tt").html(weekToHTML(week));
	});
}

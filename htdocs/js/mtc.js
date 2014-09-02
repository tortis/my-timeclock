var onClock = false;
var timeOn = 0;
var timer;
var showingPreviousWeek = false;

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

function addWeek(week, wname, hrsid) {
	$("#weeks").append( $("<div>").addClass("panel").addClass("panel-default").append(
		$("<div>").addClass("panel-heading").append(
			$("<span>").addClass("weekname").append($("<strong>").text(wname))).append(
			$("<span>").addClass("caps").text(week.Monday.toLocaleDateString())).append(
			$("<span>").attr("id", hrsid).addClass(
				"caps").addClass(
				"pull-right").text(week.Hours.toFixed(1)+" hrs"))).append(
		$("<table>").addClass("table").append(
			$("<tr>").append(
				$("<td>").text("Sun")).append(
				$("<td>").text("Mon")).append(
				$("<td>").text("Tues")).append(
				$("<td>").text("Wed")).append(
				$("<td>").text("Thur")).append(
				$("<td>").text("Fri")).append(
				$("<td>").text("Sat"))).append(
			$("<tr>").append(
				$("<td>").text(week.Days[0].Hours.toFixed(1))).append(
				$("<td>").text(week.Days[1].Hours.toFixed(1))).append(
				$("<td>").text(week.Days[2].Hours.toFixed(1))).append(
				$("<td>").text(week.Days[3].Hours.toFixed(1))).append(
				$("<td>").text(week.Days[4].Hours.toFixed(1))).append(
				$("<td>").text(week.Days[5].Hours.toFixed(1))).append(
				$("<td>").text(week.Days[6].Hours.toFixed(1))))));

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
			$("<span>").addClass("pull-right").append(
				$("<span>").addClass("caps").text("OUT ")).append(
				$("<span>").text(tout.toLocaleTimeString()))));
	}
	return r;
}

function updateDetails(week, full) {
	$("#details").empty();
	var dotw = 6;
	if (!full) {
		dotw = (new Date()).getDay();
	}
	for (;dotw >= 0; dotw--) {
		if (week.Days[dotw].Blocks == null) {
			continue;
		}
		$("#details").append(
			$("<div>").addClass("panel").addClass("panel-default").addClass("detail").append(
				$("<div>").addClass("panel-heading").text(weekday[dotw]).append(
						$("<span>").addClass("caps").addClass("pull-right").text("date"))).append(
					getBlocks(week.Days[dotw])));
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

function showCurrentWeeks() {
	$("#weeks").empty();
	var oneWeekAgo = new Date();
	oneWeekAgo.setDate(oneWeekAgo.getDate() - 7);
	getWeek(oneWeekAgo, function (week) {
		addWeek(week, "Last Week", "lwhrs");
		getWeek(new Date(), function (week) {
			addWeek(week, "This Week", "twhrs");
			updateDetails(week, false);
		});
	});
}

function showWeek(ev) {
	if (!showingPreviousWeek) {
		showingPreviousWeek = true;
		$("#clk").text("Current Week").removeClass("btn-danger").removeClass("btn-success").addClass("btn-primary");
		$("#timeon").empty();
	}
	$("#weeks").empty();
	getWeek(ev.date, function (week) {
		addWeek(week, ev.date.toLocaleDateString(), "sw");
		updateDetails(week, true);
	});
}

window.onload = function () {
	$("#clk").click(function () {
		if (showingPreviousWeek) {
			updateStatus();
			showCurrentWeeks();
			showingPreviousWeek = false;
		} else {
			if (onClock) {
				clockOut(function (s) {
					if (s) {
						onClock = false;
						buttonUpdate(onClock, 0);
						showCurrentWeeks();
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
		}
	});
	$('#dp1').attr("value", (new Date()).toLocaleDateString());
	$('#dp1').datepicker().on('changeDate', showWeek);
	updateStatus();
	showCurrentWeeks();
}

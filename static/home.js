(function() {
	var $team = $('[name="team"]');
	var $board = $('[name="board"]');
	var $code = $('code');

	$team.keyup(updateUrl);
	$board.keyup(updateUrl);
	$team.change(updateUrl);
	$board.change(updateUrl);

	function updateUrl() {
		var team = $team.val();
		var board = $board.val();

		var url = "http://scores.zduck.com/";
		if (team) {
			url += team + "/";
		}
		else {
			url += "{team}/";
		}
		if (board) {
			url += board;
		}
		else {
			url += "{board}";
		}

		if (team && board) {
			console.log('both');
			$code.html("<a href='" + url + "'>" + url + "</a>");
		}
		else {
			$code.text(url);
		}
	}

	// Do it on page load to handle the back button scenario.
	$(updateUrl);
}())
(function() {
	var bdesc = document.getElementById('board-desc');

	var previousDesc;
	bdesc.addEventListener('blur', function() {
		// Check whether the description was actually changed.
		if (bdesc.innerText !== previousDesc) {

			var url = 'desc';

			// Save the description using an AJAX POST.
			$.ajax({
				data: {
					desc: bdesc.innerText
				},
				url: window.location.pathname + '/desc',
				type: 'POST'
			});
		}
	});

	// Select the description text on focus.
	bdesc.addEventListener('focus', function() {
		// Save the initial value so we can detect whether it changes.
		previousDesc = bdesc.innerText;

		// Stolen from http://stackoverflow.com/a/3806004/31308
	    window.setTimeout(function() {
	        var sel, range;
	        if (window.getSelection && document.createRange) {
	            range = document.createRange();
	            range.selectNodeContents(bdesc);
	            sel = window.getSelection();
	            sel.removeAllRanges();
	            sel.addRange(range);
	        } else if (document.body.createTextRange) {
	            range = document.body.createTextRange();
	            range.moveToElementText(bdesc);
	            range.select();
	        }
	    }, 1);
	});
}())
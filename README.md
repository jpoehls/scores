# Go Scoreboard

This a small project I'm using to learn Go. It is a scoreboard web app
for tracking personal records between friends.

Scoreboards are organize under teams, with each team able to have
any number of boards.

Teams and boards are controlled by a simple URL scheme.

Example: http://scores.zduck.com/{team}/{board}

There isn't any privacy or security for your board, the URL
is your only protection.

Only the top 10 scores are kept.

## Wish List

### Code-wise

* Remember (and pre-populate on subsequent loads) the entered Name and Email using a cookie.
* Show Gravatar next to each record.
* Show a menu of other boards owned by the current team.
* Save the board description when it is edited.
* Create a landing page for the app with a friendly form for creating a new board.

### Server-wise

* Use a Github hook to automatically rebuild and restart the app when updates are pushed on the master branch.
* Create an `upstart` service script to use to manage the process on the server.
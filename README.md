# Go Scoreboard

This a small project I'm using to learn Go. It is a scoreboard web app
for tracking personal records between friends.

Scoreboards are organize under teams, with each team able to have
any number of boards.

Teams and boards are controlled by a simple URL scheme.

Example: `http://scores.zduck.com/{team}/{board}`

There isn't any privacy or security for your board, the URL
is your only protection.

Only the top 10 scores are kept.

## Wish List

### Code-wise

#### Features

* Reload templates on change in development.
* Add a login system. If signed in, boards created are linked to your user. You are made the "owner". You can delete and reset boards you own. (Store all users in a single JSON file. Owned boards are listed in `users.json`)
* Use correct verbiage in the footer when the board has only been updated once. (Currently says "updated 1 times")
* Implement a config file for things like the /data directory location
* Add an admin page for
	* listing all teams and boards that exist, sorted by last activity date
	* deleting boards
	* clearing scores on a board (without deleting the board itself)
* Show human friendly relative timestamps. i.e. "yesterday", "last monday", etc (create a go package for this?)

#### Refactoring

* Cache boards in memory and use a channel queue to persist changes to disk.
	* Use a lock to sync changes to the board instance in memory.

### Server-wise

* Use a Github hook to automatically rebuild and restart the app when updates are pushed on the master branch.
	* Or, setup pushing to a git repo on the server and use a hook to rebuild/restart the app that way.
* Create an `upstart` service script to use to manage the process on the server.
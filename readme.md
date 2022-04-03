<h1>NBA-Simulation</h1>
This project was developed for the İnsider Go case. Project the first week of an NBA fixture was be simulated.

<h2>Project operation:</h2>
When the project is run for the first time, it checks the groups, teams, players, games to be played, and players of the
games from the database. If there are missing data, it completes the missing data. Simulates games after checking for
missing data. Records events from all games in the "Events" table in seconds. It just shows the events in order on the
site. In fact, all events were determined at the beginning of the project.

<h2>For run:</h2>
Start docker from local machine.

````shell
chmod +x build.sh
./build.sh
````

<h3>Then you visit <a href="http://localhost:3003">localhost:3003</a></h3>

Folder Structure
----------------
<ul>
<li>cmd:<br>
This folder contains main function. <br>
The project starts with this function.<br>
This function calls database connection, route setup and http server starts functions.
</li>
<li>data:<br>
This folder contains data creation functions. The functions in this folder are called when this project starts running. Groups, teams, players, games, players of games, and activities done in games are created by functions in this folder.
</li>
<li>model:<br>
This folder holds files containing the models required for the project.
</li>
<li>router:<br>
This folder contains route definition function.
</li>
<li>storage:<br>
This folder contains database clients, and connect functions.
<ul>
<li>postgre:<br>
This folder contains postgre database client, connect and migration functions.
</li>
</ul>
</li>
<li>utils:<br>
This folder contains utility functions.
</li>
<li>views:<br>
This folder contains html files.
</li>
<li>ws:<br>
This folder contains web socket functions.
</li>
</ul>
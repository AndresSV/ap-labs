# Architecture of Pacman Challenge

Multithreaded Pacman Game
=========================

Implement a multithreaded version of the arcade video game [Pacman](https://en.wikipedia.org/wiki/Pac-Man). This version will be a
Computer vs Human game. Each enemy will be independent and the number of enemies is configurable. Below you can see the general
requirements for the enemies and human player (pacman) interaction.

![Pacman](pacman.png)

Technical Requirements
----------------------
- The game's maze layout can be static.
- The `pacman` gamer must be controlled by the user.
- Enemies are autonomous entities that will move in a random way.
- Enemies and pacman should respect the layout limits and walls.
- Enemies number can be configured on game's start.
- Each enemy's behaviour will be implemented as a separate thread.
- Enemies and pacman threads must use the same map or game layout data structure resource.
- Display obtained pacman's scores.
- Pacman loses when an enemy touches it.
- Pacman wins the game when he has taken all the coins in the map.

#Requirements

– Go installed
– Terminal where we can run emojis and emoticons.


## Usage 

– make test

Or

– go run pacman.go

The development of this project was as simple as making the core components as small as possible:

### Files as map
– We get the core design of our system designed on txt files, there we can make fast changes and not worrying about the code
– The map is divided into three core components the upper and lower level of our map and the middle part which is generated with a function ghostGen which main purpose is to initialize the position of our ghosts

 ### Pacman
– Then we have our controllable character which we are taking advantage of modern terminals which can execute emojis on the cli.
– We control our character with the common direction arrows so we can tell to move up, left, down and right

### Ghosts
– Ghost are independent entities that are continuously moving unfortunately we couldn't make an AI that could make use of algorithms like A* or use the current position of the character like the original one

### Code Functionality

![TASK FLOW](pacman.png)

Structs:
sprite: For coordinates information



| Attibutes                  |   DataType   | Description                                                                      |
|----------------------------|--------------|----------------------------------------------------------------------------------|
|     row                |     int   | Actual row on the map                                             |
|     col                |     int   | Actual column on the map                                                 |
|     startRow           |  int      | Initial row on the map                                   |
|     startCol                |     int   | Initial column on the map                                               |
| Attibutes                  |   DataType   | Description                                                                      |
|----------------------------|--------------|----------------------------------------------------------------------------------|
|     position                |     sprite   | struct of type sprite                                           |
|     status                |     GhostStatus   | struct of type sprite GhostStatus                                             |


| Attibutes                  |   DataType   | Description                                                                      |
|----------------------------|--------------|----------------------------------------------------------------------------------|
|     GhostStatusNormal      |     string   |   "Normal"                                      |
|     GhostStatusBlue       |     string   | "Blue"                                            |




| Attibutes                  |   DataType   | Description                                                                      |
|----------------------------|--------------|----------------------------------------------------------------------------------|
|     Player                  |      string     |   Player emoticon                                 |
|     Ghost                  |      string     | Ghost emoticon                                           |
|     GhostBlue                  |      string     |   Vulnerable ghost emoticon                                 |
|     Wall                  |      string     | Wall emoticon                                           |
|     Dot                   |      string     |   Dot emoticon                                 |
|     Pill                  |      string     | Pill emoticon                                           |
|     Space                     |      string     |   Blank Space emoticon                                 |
|     PillDurationSec                  |      time.Duration     | Time available for pill                                 |



`ghost`: To represent each of the ghosts with its own parameters.
`ghostStatus`: To check if the player can ‘eat’ one of the ghosts or not.
`Config`: to represent the information of each of the emojis used in the game
`init()`: Golang native method that runs before the main method with the specifications for printing our game.
`getIcons()`: Where we can declare which emoticons our system uses.
`createMap()`: Here we specify the files to read to build the top and bottom part of the game as well as the instance of the ghosts
`ghostGen()`: String generator with the number of ghosts for the createMap function.
`makeMove()`: Whether or not our character can move where the input was pressed.

# Consideration

– We opted out for the use of flags, because most of the game's configuration is adjusted once it's running.
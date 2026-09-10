*This activity has been created as part of the 42 curriculum by rfeghali, mmoskale.*

Description
section that clearly presents the activity, including its goal and a brief overview.

Instructions
section containing any relevant information about compilation, installation, and/or execution.

Resources 
section listing classic references related to the topic (documentation, articles, tutorials, etc.), as well as a description of how AI was used — specifying for which tasks and which parts of the activity.

Any required additions will be explicitly listed below.
Architecture
section explaining your server design choices (dispatcher/router vs inline handling, concurrency model, etc.).

Protocol Implementation
section documenting any deviations from RFC 42TAP and justifying your choices.

Combat System
section describing your turn-based combat mechanics, damage formulas, initiative order, and additional combat commands (DEFEND, FLEE, etc.).


World Design
section describing your world layout, room connections, NPC roles, and item distribution.

Server Logging
section documenting your logging implementation, including log format, event types, output destinations, and how to monitor server behavior and detect abuse patterns.

Group Contributions
section clearly indicating each team member’s responsibilities and contributions to different components (server, CLI client, GUI client, world design, etc.).

Building and Running
section with detailed instructions for your chosen building tool and how to run each component (server, CLI client, GUI client).

Testing
section explaining how to test the multiplayer functionality, combat system, and quest mechanics.



## Architecture

The server is implemented in Go and uses TCP connections for communication between clients and the game server.

The main architecture consists of:

* A TCP server listening on port `8090`.
* One goroutine per client connection.
* A central `CommandRegistry` containing all supported commands.
* Shared world state protected by server-level mutexes.
* Per-player state protected by a player-level mutex.
* World data loaded from `data.json`.
* A shared logging system for server events, errors, and world-state changes.

The main player state includes the player's current room, inventory, combat state, statistics, group membership, and quest states.

The server starts by loading the world and then accepts TCP connections. Each connection is handled independently in a goroutine.

## Protocol Implementation

Communication uses a line-based TCP protocol. Clients send commands as text lines and the server responds with either:

```text
OK ...
```

or:

```text
ERR <code> <message>
```

The initial connection response is:

```text
OK hello proto=1
```

The server uses a central command registry to validate commands, their arguments, and authentication requirements before executing handlers.

Structured information such as rooms, inventory, combat results, NPC dialogue, and quests is returned using JSON structures defined in `common/api.go`.

## Error Codes

The server uses numeric error codes to allow clients to distinguish between different categories of failures. The numeric code is followed by a descriptive error message.

| Code    | Error                        | Why it is used                                                                                                                                                        |
| ------- | ---------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **201** | `NAME_IN_USE`                | The requested username is already being used by another connected player.                                                                                             |
| **301** | `NO_EXIT`                    | The player attempted to move in a direction for which the current room has no exit.                                                                                   |
| **400** | General client/command error | Used when the command is invalid, arguments are missing/incorrect, the player is in an invalid state, or an action is not currently allowed.                          |
| **401** | `NOT_AUTHENTICATED`          | Used when a command requiring authentication is sent before the player has successfully connected.                                                                    |
| **403** | `CANNOT_MOVE_IN_COMBAT`      | Prevents a player from changing rooms while they are currently fighting an NPC.                                                                                       |
| **404** | Resource not found           | Used when the requested NPC or item cannot be found in the current game context.                                                                                      |
| **405** | `NPC_NOT_HOSTILE`            | Prevents the player from attacking an NPC that is not marked as hostile.                                                                                              |
| **406** | Quest/action not available   | Used when a quest cannot be accepted, for example when an NPC has no quest, the quest has already been accepted, or it has already been completed.                    |
| **409** | Conflict                     | Used when the requested action conflicts with the player's current state, such as trying to attack another NPC while already fighting a different NPC.                |
| **500** | Server/internal error        | Used for unexpected command-handler errors or missing quest definitions, indicating that the problem is on the server side rather than caused by normal player input. |


## Combat System

Combat is turn-based at the command level. Each `ATTACK` command represents one combat exchange:

1. The player attacks the NPC.
2. Damage is calculated from attack minus defense.
3. If the NPC survives, it counter-attacks.
4. Both resulting HP values are sent to the client.
5. Combat ends when either side is defeated.

Damage uses:

```text
damage = max(1, attack - defense)
```

The player starts with 100 HP, 12 attack, and 10 defense.

The player can also use `FLEE`. Fleeing has a 50% success chance. A failed attempt costs 5 HP.

The current implementation does not contain an initiative system or a `DEFEND` command. Therefore, the combat model is based on the order of the `ATTACK` command rather than initiative.

## Quest System

Quests are given by NPCs marked as quest givers. The quest system stores each player's quest state as either `active`, `completed` or `abandoned`.

A quest cannot be taken twice:

* If the player has never taken it, the quest becomes `active`.
* If it is already `active`, the server returns `QUEST_ALREADY_ACCEPTED`.
* If it is already `completed`, the server returns `QUEST_ALREADY_COMPLETED`.

A quest can be abandoned and once abandoned it cannot be taken again.

There are two quest types:

1. **Fetch item** – the player must obtain a specific item.
2. **Defeat NPC** – the player must defeat a specific NPC.

Quest completion is automatic. When the required objective is detected, the server changes the quest state to `completed`, gives the reward immediately, and sends a quest completion event to the player.

### Quest Logging

Quest state changes are logged so that quest progression can be audited.

The server should record:

```text
QUEST_ACCEPT
QUEST_COMPLETE
QUEST_ABANDON
QUEST_REWARD
```

Example:

```text
INFO QUEST_ACCEPT player=Alice quest=lost_amulet npc=Blacksmith Torin room=blacksmith
INFO QUEST_COMPLETE player=Alice quest=lost_amulet type=fetch_item target=silver_amulet reward=iron_sword
INFO QUEST_REWARD player=Alice quest=lost_amulet item=iron_sword
```

This makes it possible to determine when a player accepted a quest, when the objective was completed, and which reward was distributed and if the quest was abandoned.

## World Design

The world is defined in `data.json`.

The map contains several interconnected locations, including:

* Village Square
* The Prancing Pony
* Tavern Cellar
* Ancient Catacombs
* Forgotten Crypt
* Whispering Woods
* General Store
* Blacksmith & Armory

NPCs have different roles, including:

* Quest giver
* Dialogue NPC
* Trader
* Hostile enemy

For example, the Village Guard and Blacksmith Torin are quest-giving NPCs, while the Giant Cellar Rat is a hostile enemy.

Item distribution is defined per location in the world data. Items can be marked as obtainable, allowing players to take them from the world.

The map layout was designed using Claude and then represented in the game's world configuration.

## Server Logging

The server provides a centralized logger with three levels:

* `INFO`
* `WARN`
* `ERROR`

The logger uses a mutex so that simultaneous client goroutines do not write conflicting log messages.

The server already logs:

* Client connections and disconnections.
* Commands received.
* Authentication.
* Error responses.
* Combat activity.
* World state change.
* Quest progression.
* Potential abuse pattern.

Example:

```text
INFO WORLD_MOVE player=Alice from=start to=tavern direction=north
INFO WORLD_ITEM_TAKE player=Alice item=silver_amulet room=crypt
INFO WORLD_ITEM_DROP player=Alice item=silver_amulet room=start
INFO NPC_TALK player=Alice npc=Blacksmith Torin room=blacksmith
INFO QUEST_ACCEPT player=Alice quest=lost_amulet npc=Blacksmith Torin
INFO QUEST_COMPLETE player=Alice quest=lost_amulet type=fetch_item target=silver_amulet reward=iron_sword
INFO QUEST_REWARD player=Alice quest=lost_amulet item=iron_sword
INFO COMBAT_ROUND player=Alice npc=giant_rat player_damage=7 npc_damage=2 npc_hp=13 player_hp=98
INFO NPC_DEFEATED player=Alice npc=giant_rat
```

These logs provide an audit trail for all important changes to the game world and player progression.

## Group Contributions

* `rfeghali` – `CLI client and handlers implementation.`

* `mmoskale` – `<contribution>`


## Building and Running

### Server

From the `server` directory:

```bash
go run .
```

The server listens on:

```text
localhost:8090
```

### CLI Client

From the `cli-client` directory:

```bash
go run . <username>
```

The client connects to the TCP server and allows the player to enter commands interactively.

### GUI


```text
<GUI startup command>
```

## Testing

The server should be tested using multiple clients to verify multiplayer state and concurrency.

### Connection Tests

* Connect a player successfully.
* Attempt to connect with a username that is already in use.
* Attempt authenticated commands before authentication.
* Connect and disconnect multiple players.

### Movement Tests

* Move through valid exits.
* Attempt to move through an invalid exit.
* Attempt to move while in combat.
* Verify that room presence events are sent correctly.

### Item Tests

* Take an obtainable item.
* Attempt to take an item that is not present.
* Drop an item.
* Attempt to drop an item that is not in the inventory.
* Verify item movement is logged.

### NPC Tests

* Talk to an NPC.
* Attempt to talk to an NPC that is not present.
* Attempt to attack a non-hostile NPC.
* Verify NPC interactions are logged.

### Combat Tests

* Attack a hostile NPC.
* Verify player and NPC damage.
* Verify NPC counter-attacks.
* Defeat an NPC.
* Verify player defeat and respawn.
* Test successful and failed fleeing.
* Verify combat results are logged.

### Quest Tests

* Take a quest from a quest-giving NPC.
* Attempt to take the same quest again.
* Complete a fetch quest.
* Complete a defeat-NPC quest.
* Verify the quest automatically changes to `completed`.
* Verify the reward is immediately added to the inventory.
* Attempt to take a completed quest again.
* Verify `QUEST_ACCEPT`, `QUEST_COMPLETE`, and `QUEST_REWARD` logs are generated.

### Logging Tests

For every world-state-changing action, verify that a corresponding log entry is produced.

At minimum:

```text
MOVE       -> WORLD_MOVE
TAKE       -> WORLD_ITEM_TAKE
DROP       -> WORLD_ITEM_DROP
TALK       -> NPC_TALK
QUEST      -> QUEST_ACCEPT
completion -> QUEST_COMPLETE
reward     -> QUEST_REWARD
ATTACK     -> COMBAT_ROUND
defeat     -> NPC_DEFEATED / PLAYER_DEFEATED
```

This ensures the server provides an auditable record of player actions and important changes to the shared game state.


bufio.Reader -> low level bufferred byte/stream reading. No token size limit. Keeps the delimiter in  the string
bufio.Scanner -> tokenization (splitting by line , word ..), default 64 kb max token size. Strips the delimiter from the token.

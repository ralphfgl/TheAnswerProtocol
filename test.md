.4.7 QUEST Command

Syntax: QUEST

Purpose: Request quest from quest-giver NPC

Responses:

    OK - Quest available, returns quest information
    ERR 404 NPC_NOT_FOUND - NPC not present
    ERR 406 NO_QUEST_AVAILABLE - NPC has no quests or quest already completed

Example:

Client → Server: QUEST merchant
Server → Client: OK {"quest_id": "fetch_herbs", "description": "Bring me 3 healing herbs", "reward": "gold_coin", "status": "available"}

5.4.8 QUESTS Command

Syntax: QUESTS

Purpose: List all player's active and completed quests

Response: OK

Example:

Client → Server: QUESTS
Server → Client: OK [{"quest_id": "fetch_herbs", "status": "active", "progress": "1/3"}, {"quest_id": "defeat_goblin", "status": "completed"}]

IMPORTANT: This RFC intentionally leaves certain aspects underspecified to allow for creative implementation by development teams. The following areas require design decisions and justification:
6.1.1 Combat System Implementation

Basic QUEST and QUESTS commands are provided, but implementers must design:

    Quest progression: How quest objectives are tracked and validated
    Quest completion: Automatic vs manual completion mechanisms
    Quest rewards: Distribution and inventory management
    Quest dependencies: Prerequisites and quest chains
    Additional commands: COMPLETE_QUEST, ABANDON_QUEST, or similar

Implementation teams MUST document their design choices and provide justification for their combat and quest system implementations in their project README.


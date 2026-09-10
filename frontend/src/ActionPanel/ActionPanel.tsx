import { useEffect } from 'react'
import './ActionPanel.css'

export interface Inventory {
    type: string;
    items: number[];
}

export interface Group {
    type: string;
    group_list: string[] | null;
    my_group: string;
}

export interface QuestStatus {
    type: string;
    quest_map: Record<string, string>;
    count: number;
}

export interface QuestObject {
    id: string;
    title: string;
    giver: string;
    type: string;
    target: string;
    description: string;
    reward_item: string;
    dialogue_start: string;
    dialogue_complete: string;
}

export interface QuestData {
    type: string;
    quest: QuestObject
    status: string;
}


interface ActionPanelProps {
    onCommand: (command: string) => void;
    inventory: Inventory | null
    group: Group | null
    questData: QuestData | Record<string, QuestData>
}

function ActionPanel(props: ActionPanelProps) {
    const { onCommand, inventory, group, questData } = props;
    const actions = ["LOOK", "WHO", "INVENTORY", "STATUS", "QUESTS"]

    const createGroup = () => {
        onCommand("GROUP CREATE")
        onCommand("GROUP DISPLAY")
    }

    const joinGroup = (value: string) => {
        onCommand("GROUP JOIN " + value)
        onCommand("GROUP DISPLAY")
    }

    const leaveGroup = () => {
        onCommand("GROUP LEAVE")
        onCommand("GROUP DISPLAY")
    }

    useEffect(() => {
        const timer = setInterval(() => {
            onCommand("GROUP DISPLAY");
        }, 15000);
        return () => clearInterval(timer);
    }, [onCommand]);

    return (
        <section className='action_panel'>
            <h2 className='action_title'>Inventory and action panel</h2>
            <p>My Inventory:</p>
            {inventory?.items && inventory.items.length > 0
                ?
                (inventory.items.map((value, index) => (
                    <div key={index} className='inventory_item'>
                        <p>{value}</p>
                        <button className='inventory_button' onClick={() => onCommand("DROP " + value)}>DROP</button>
                    </div>
                )))
                :
                <p>The inventory is empty</p>
            }
            <hr className='line' />
            <p>Quests:</p>
            {Object.values(questData).length > 0 ? (
                Object.values(questData).map((quest) => (
                    <div
                        key={quest.quest.id}
                        className={`quest_item ${quest.status === "active" ? 'quest_active' : quest.status === "abandoned" ? 'quest_abondoned' : 'quest_done'}`}
                    >
                        <p>{quest.quest.description}</p>
                        {quest.status === "active" && <button className='action_button' onClick={() => onCommand("ABANDON_QUEST " + quest.quest.id)}>Abandon</button>}
                    </div>
                ))
            ) : (
                <p>No active quests</p>
            )}
            <hr className='line' />
            <p>Actions:</p>
            <div className='action_list'>
                {actions.map((value, index) => (
                    <button className='action_button' onClick={() => onCommand(value)} key={index}>{value}</button>
                ))}
            </div>
            <hr className='line' />
            <p>Groups:</p>
            <button className='action_button' onClick={() => createGroup()}>GROUP CREATE</button>
            {
                group?.group_list && group.group_list.length > 0
                    ?
                    (group?.group_list?.sort().map((value, index) => {
                        const isMyGroup = group.my_group === value;
                        const hasNoGroup = !group.my_group;
                        return (
                            <div className='group_item' key={index}>
                                <p>{value.toUpperCase()}</p>
                                {isMyGroup && <button className='action_button' onClick={() => leaveGroup()}>Leave</button>}
                                {hasNoGroup && <button className='action_button' onClick={() => joinGroup(value)}>Join</button>}
                            </div>
                        )
                    }))
                    :
                    <p>No groups created</p>
            }
        </section >
    )
}

export default ActionPanel

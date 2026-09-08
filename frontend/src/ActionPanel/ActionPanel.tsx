import { useEffect } from 'react'
import './ActionPanel.css'

function ActionPanel({ onCommand, inventory, group }) {
    const actions = ["LOOK", "WHO", "INVENTORY", "STATUS"]

    const createGroup = () => {
        onCommand("GROUP CREATE")
        onCommand("GROUP DISPLAY")
    }

    const joinGroup = (value) => {
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
        }, 10000);
        return () => clearInterval(timer);
    }, [onCommand]);

    return (
        <section className='action_panel'>
            <h2 className='action_title'>Inventory and action panel</h2>
            <p>My Inventory:</p>
            {inventory?.items?.length > 0
                ?
                <ul className='inventory_list'>
                    {(inventory.items.map((value, index) => (
                        <li key={index}>
                            <div className='inventory_item'>
                                <p>{value}</p>
                                <button className='inventory_button' onClick={() => onCommand("DROP " + value)}>DROP</button>
                            </div>
                        </li>
                    )))}
                </ul>
                :
                <p>The inventory is empty</p>
            }
            <p>Quests:</p>
            <p>Actions:</p>
            <div className='action_list'>
                {actions.map((value, index) => (
                    <button className='action_button' onClick={() => onCommand(value)} key={index}>{value}</button>
                ))}
            </div>
            <p>Groups:</p>
            <button className='action_button' onClick={() => createGroup()}>GROUP CREATE</button>
            {
                group?.group_list?.length > 0
                    ?
                    (group.group_list.map((value, index) => {
                        group.group_list = group.group_list.sort()
                        const isMyGroup = group.my_group === value;
                        const hasNoGroup = !group.my_group;
                        return (
                            <div className='group_list' key={index}>
                                <p>{value}</p>
                                {isMyGroup && <button className='action_button' onClick={() => leaveGroup()}>Leave</button>}
                                {hasNoGroup && <button className='action_button' onClick={() => joinGroup(value)}>Join</button>}
                            </div>
                        )
                    }))
                    :
                    <p>No groups cread</p>
            }
        </section >
    )
}

export default ActionPanel

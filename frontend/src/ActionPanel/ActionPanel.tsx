import './ActionPanel.css'

function ActionPanel({ onCommand, inventory }) {
    const actions = ["LOOK", "WHO", "INVENTORY", "STATUS"]
    return (
        <section className='action_panel'>
            <h2 className='action_title'>Inventory and action panel</h2>
            <p>My Inventory:</p>
            {inventory.length > 0
                ?
                <ul className='inventory_list'>
                    {(inventory.map((value, index) => (
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
        </section >
    )
}

export default ActionPanel

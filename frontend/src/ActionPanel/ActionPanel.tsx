import './ActionPanel.css'

function ActionPanel({ onCommand }) {
    const actions = ["LOOK"]
    return (
        <section className='action_panel'>
            <h2 className='action_title'>Inventory and action panel</h2>
            <p>My Inventory:</p>
            <p>Quests:</p>
            <p>Actions:</p>
            {actions.map((value, index) => (
                <button className='action_button' onClick={() => onCommand(value)} key={index}>{value}</button>
            ))}
        </section>
    )
}

export default ActionPanel

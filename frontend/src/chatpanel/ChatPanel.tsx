import './ChatPanel.css'

function ChatPanel() {
    const tabs = ["Global", "Room", "Group", "System log"]
    return (
        <section className='chat_panel'>
            <h2>Communication and logs</h2>
            <div className='tabs'>
                {tabs.map((value, index) => (
                    <button className='tab_button' key={index}>{value}</button>
                ))}
            </div>
            <hr />
            <p>Logs</p>
            <p>Logs</p>
            <p>Logs</p>
            <p>Logs</p>
            <p>Logs</p>
            <hr />
            <form>
                <input type="text" placeholder='Chat...' />
                <button>Send</button>
            </form>
        </section>
    )
}

export default ChatPanel

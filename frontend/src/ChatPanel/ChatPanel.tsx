import { useState } from 'react'
import './ChatPanel.css'

function ChatPanel({ onCommand }) {
    const tabs = ["Global", "Room", "Group", "System log"]
    const [message, setMessage] = useState("")
    const [activeTab, setActiveTab] = useState("GLOBAL")

    const submitMessage = (event) => {
        event.preventDefault()
        onCommand(`CHAT ${activeTab} ${message}`)
        setMessage("")
    }
    return (
        <section className='chat_panel'>
            <h2>Communication and logs</h2>
            <div className='tabs'>
                {tabs.map((value, index) => (
                    <button className='tab_button' key={index} onClick={() => setActiveTab(value.toUpperCase())}>{value}</button>
                ))}
            </div>
            <hr />
            <p>Logs</p>
            <p>Logs</p>
            <p>Logs</p>
            <p>Logs</p>
            <p>Logs</p>
            <hr />
            <form onSubmit={(event) => submitMessage(event)}>
                <input value={message} onChange={(e) => setMessage(e.target.value)} type="text" placeholder='Chat...' />
                <button>Send</button>
            </form>
        </section>
    )
}

export default ChatPanel

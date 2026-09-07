import { act, useState } from 'react'
import './ChatPanel.css'

function ChatPanel({ onCommand, messages }) {
    const tabs = ["Global", "Room", "Group", "System"]
    const [message, setMessage] = useState("")
    const [activeTab, setActiveTab] = useState("Global")

    const submitMessage = (event) => {
        event.preventDefault()
        onCommand(`CHAT ${activeTab} ${message}`)
        setMessage("")
    }
    return (
        <section className='chat_panel'>
            <h2 className='chat_title'>Communication and logs</h2>
            <div className='tabs'>
                {tabs.map((value, index) => (
                    <button className={activeTab == value ? 'tab_button--active' : 'tab_button'} key={index} onClick={() => setActiveTab(value)}>{value}</button>
                ))}
            </div>
            <hr />
            <p>Logs:</p>
            <div className='chat_logs'>
                {messages.map((value, index) => (
                    <p key={index}>{value}</p>
                ))}
            </div>
            <hr />
            <form className='chat_form' onSubmit={(event) => submitMessage(event)}>
                <input className='chat_input' value={message} onChange={(e) => setMessage(e.target.value)} type="text" placeholder='Chat...' />
                <button className='tab_button'>Send</button>
            </form>
        </section>
    )
}

export default ChatPanel

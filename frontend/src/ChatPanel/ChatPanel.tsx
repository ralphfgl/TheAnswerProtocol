import { useState, type SubmitEvent } from 'react'
import './ChatPanel.css'

interface ChatPanelProps {
    onCommand: (command: string) => void;
    messages: string[]
}

function ChatPanel(props: ChatPanelProps) {
    const { onCommand, messages } = props
    const tabs = ["Global", "Room", "Group", "System"]
    const [message, setMessage] = useState("")
    const [activeTab, setActiveTab] = useState("GLOBAL")

    const submitMessage = (event: SubmitEvent) => {
        event.preventDefault()
        onCommand(`CHAT ${activeTab} ${message}`)
        setMessage("")
    }

    const filteredMessages = messages.filter((message: string) => {
        if (activeTab === "GLOBAL") {
            return message.includes("EVT GLOBAL CHAT")
        }
        if (activeTab === "ROOM") {
            return message.includes("EVT ROOM CHAT")
        }
        if (activeTab === "GROUP") {
            return message.includes("EVT GROUP CHAT")
        }
        if (activeTab === "SYSTEM") {
            return !message.includes("CHAT")
        }
        return true
    })

    return (
        <section className='chat_panel'>
            <h2 className='chat_title'>Communication and logs</h2>
            <div className='tabs'>
                {tabs.map((value, index) => (
                    <button className={activeTab == value.toUpperCase() ? 'tab_button--active' : 'tab_button'} key={index} onClick={() => setActiveTab(value.toUpperCase())}>{value}</button>
                ))}
            </div>
            <hr className='line' />
            <div className='chat_logs'>
                {filteredMessages.map((value, index) => (
                    <p key={index}>{value}</p>
                ))}
            </div>
            <hr className='line' />
            {activeTab !== "SYSTEM" &&
                <form className='chat_form' onSubmit={(event) => submitMessage(event)}>
                    <input className='chat_input' value={message} onChange={(e) => setMessage(e.target.value)} type="text" placeholder='Chat...' />
                    <button className='tab_button'>Send</button>
                </form>
            }

        </section>
    )
}

export default ChatPanel

import './App.css'
import { useEffect, useRef, useState } from "react"
import Header from './Header/Header'
import ChatPanel from './ChatPanel/ChatPanel'
import RoomView from './RoomView/RoomView'
import ActionPanel from './ActionPanel/ActionPanel'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [nickname, setNickname] = useState("")
  const [messages, setMessages] = useState([])
  const [inputValue, setInputValue] = useState("")
  const wsRef = useRef(null)
  const [roomData, setRoomData] = useState([])

  const headerdata =
  {
    "name": nickname,
    "health": 100,
    "players": 10,
    "players_room": 2
  }

  const parseMessage = (message) => {
    if (message.startsWith("OK connected")) {
      setIsAuthenticated(true)
    }
    else if (message.startsWith("{")) {
      const data = JSON.parse(message)
      if (data.type == "room") {
        setRoomData(data)
      }
    }
  }

  useEffect(() => {
    let isMounted = true;
    let reconnectTimeout;
    function connect() {
      const websocket = new WebSocket("ws://localhost:8080/ws");
      wsRef.current = websocket;

      websocket.onopen = function () {
        console.log("Connected to WebSocket server");
      };

      websocket.onmessage = function (event) {
        parseMessage(event.data)
        setMessages(prevMessage => [...prevMessage, event.data])
      };

      websocket.onclose = function () {
        console.log("WebSocket connection closed, retrying...");
        if (isMounted) {
          console.log("retrying in 1 second...");
          reconnectTimeout = setTimeout(connect, 1000);
        }
      };

      websocket.onerror = function (error) {
        console.error("WebSocket error:", error);
      };
    }
    connect();
    return () => {
      isMounted = false;
      clearTimeout(reconnectTimeout);
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [])

  const sendMessage = (e) => {
    e.preventDefault()
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(inputValue)
    }
    setInputValue("")
  }

  const sendCommand = (command) => {
    console.log(command)
    wsRef.current.send(command)
  }

  const submitLogin = (e) => {
    e.preventDefault()
    wsRef.current.send("CONNECT " + nickname)
  }

  const handleLogout = () => {
    setIsAuthenticated(false)
    setNickname("")
  }

  if (!isAuthenticated) {
    return (
      <form className='login_form' onSubmit={(event) => submitLogin(event)}>
        <input className='login_input' required type="text" placeholder='Enter your name' value={nickname} onChange={(e) => setNickname(e.target.value)} />
        <button className='login_button'>Apply</button>
      </form>
    )
  }

  return (
    <>
      <main className='main'>
        <h1 className='title'>The answer Protocol</h1>
        <Header data={headerdata} onLogout={handleLogout} />
        <div className='panel_list'>
          <ChatPanel />
          <RoomView data={roomData} onCommand={sendCommand} />
          <ActionPanel />
        </div>
        <form onSubmit={sendMessage}>
          <input value={inputValue} onChange={(e) => setInputValue(e.target.value)} />
          <button>Send</button>
        </form>
        {messages.map((value, index) => (
          <p key={index}>{value}</p>
        ))}
      </main>
    </>
  )
}

export default App

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
  const [roomData, setRoomData] = useState({})
  const [inventoryData, setInventoryData] = useState({})
  const [headerData, setHeaderData] = useState({})
  const [talkData, setTalkData] = useState({})
  const [groupData, setGroupData] = useState({})
  const [playersServer, setPlayersServer] = useState(0)
  const [playersRoom, setPlayersRoom] = useState(0)
  const [attackData, setAttackData] = useState("")


  const parseMessage = (message: string) => {
    if (message.startsWith("OK connected")) {
      setIsAuthenticated(true)
      wsRef.current.send("LOOK")
      wsRef.current.send("STATUS")
      wsRef.current.send("GROUP DISPLAY")
      wsRef.current.send("WHO")
    }
    else if (message.startsWith("OK {")) {
      let data = JSON.parse(message.substring(3))
      console.log(data)
      if (data.room) {
        setRoomData(data)
        console.log(data.players.length)
        if (data?.players?.length < 1) {
          setPlayersRoom(1)
        }
        else {
          setPlayersRoom(data.players?.length + 1)
        }
      }
      if (data.type == "status") {
        setHeaderData(data)
      }
      if (data.type == "inventory") {
        setInventoryData(data)
      }
      if (data.type == "talk") {
        setTalkData(data)
      }
      if (data.type == "group") {
        setGroupData(data)
      }
      if (data.type == "combat") {
        setHeaderData((prev) => ({ ...prev, hp: data.attacker_hp, status: data.status }))
      }
    }
    else if (message.startsWith("OK players=")) {
      setPlayersServer(message.substring(11))
    }
    // else if (message.startsWith("OK group=")) {
    //   setGroupData((prevGroup) => [...prevGroup, message.substring(9)])
    // }
    else if (message.startsWith("ERR")) {
      console.error(message.substring(3))
    }
    else if (message.startsWith("EVT ROOM COMBAT")) {
      setAttackData(message.substring(15))
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

  const sendMessage = (e: Event) => {
    e.preventDefault()
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(inputValue)
    }
    setInputValue("")
  }

  const sendCommand = (command: string) => {
    console.log(command)
    wsRef.current.send(command)
  }

  const submitLogin = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()
    wsRef.current.send("CONNECT " + nickname)
  }

  const handleLogout = () => {
    wsRef.current.send("QUIT")
    setIsAuthenticated(false)
    setNickname("")
  }

  if (!isAuthenticated) {
    return (
      <form className='login_form' onSubmit={(event) => submitLogin(event)}>
        <input className='login_input' required type="text" maxLength={15} placeholder='Enter your name' value={nickname} onChange={(e) => setNickname(e.target.value)} />
        <button className='login_button'>Apply</button>
      </form>
    )
  }
  return (
    <>
      <main className='main'>
        <h1 className='title'>The Answer Protocol</h1>
        <Header data={headerData} playersRoom={playersRoom} playersServer={playersServer} nickname={nickname} onLogout={handleLogout} />
        <div className='panel_list'>
          <ChatPanel onCommand={sendCommand} messages={messages} />
          <RoomView data={roomData} talk={talkData} attack={attackData} onCommand={sendCommand} />
          <ActionPanel onCommand={sendCommand} inventory={inventoryData} group={groupData} />
        </div>
        {/* <form onSubmit={sendMessage}>
          <input value={inputValue} onChange={(e) => setInputValue(e.target.value)} />
          <button>Send</button>
        </form>
        {messages.map((value, index) => (
          <p key={index}>{value}</p>
        ))} */}
      </main>
    </>
  )
}

export default App

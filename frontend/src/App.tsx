import './App.css'
import { useEffect, useRef, useState } from "react"
import Header, { type HeaderData } from './Header/Header'
import ChatPanel from './ChatPanel/ChatPanel'
import RoomView, { type Data, type Talk } from './RoomView/RoomView'
import ActionPanel, { type Group, type Inventory } from './ActionPanel/ActionPanel'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [nickname, setNickname] = useState<string>("")
  const [messages, setMessages] = useState<string[]>([])
  const wsRef = useRef<WebSocket | null>(null)
  const [roomData, setRoomData] = useState<Data | null>(null)
  const [inventoryData, setInventoryData] = useState<Inventory | null>(null)
  const [headerData, setHeaderData] = useState<HeaderData | null>(null)
  const [talkData, setTalkData] = useState<Talk | null>(null)
  const [groupData, setGroupData] = useState<Group | null>(null)
  const [playersServer, setPlayersServer] = useState<number>(0)
  const [playersRoom, setPlayersRoom] = useState<number>(0)
  const [attackData, setAttackData] = useState("")
  const [questData, setQuestData] = useState({})


  const parseMessage = (message: string) => {
    if (message.startsWith("OK connected")) {
      setIsAuthenticated(true)
      wsRef.current?.send("LOOK")
      wsRef.current?.send("STATUS")
      wsRef.current?.send("GROUP DISPLAY")
      wsRef.current?.send("WHO")
    }
    else if (message.startsWith("OK {")) {
      const data = JSON.parse(message.substring(3))
      if (data.type == "room") {
        setRoomData(data)
        setPlayersRoom((data.players?.length ?? 0) + 1)
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
        setHeaderData((prev) => (prev ? { ...prev, hp: data.attacker_hp, status: data.status } : prev));
      }
      if (data.type === "quest") {
        setQuestData((prev) => ({ ...prev, [data.quest.id]: data }))
      }
    }
    else if (message.startsWith("EVT STATS players=")) {
      setPlayersServer(parseInt(message.substring(18)))
      wsRef.current?.send("LOOK")
    }
    else if (message.startsWith("EVT ROOM ITEM_TAKEN") || message.startsWith("EVT ROOM ITEM_DROP")) {
      wsRef.current?.send("LOOK")
    }
    else if (message.startsWith("ERR")) {
      if (message.startsWith("ERR 406 NO_QUEST_AVAILABLE")) {
        alert("The quest is unavailable")
      }
      else if (message.startsWith("ERR 406 QUEST_ALREADY_ABANDONED")) {
        alert("The quest is abandoned")
      }
      else if (message.startsWith("ERR 406 QUEST_PREREQUISITE_NOT_MET")) {
        alert("You need to complete other quest first")
      }
      else {
        console.error(message.substring(3))
      }
    }
    else if (message.startsWith("EVT ROOM COMBAT")) {
      setAttackData(message.substring(15))
    }
    else if (message.startsWith("EVT ROOM KILL")) {
      setAttackData(message.substring(13))
      wsRef.current?.send("LOOK")
    }
  }

  useEffect(() => {
    let isMounted = true;
    let reconnectTimeout: number;
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

  const sendCommand = (command: string) => {
    console.log(command)
    wsRef.current?.send(command)
  }

  const submitLogin = (e: React.SubmitEvent<HTMLFormElement>) => {
    e.preventDefault()
    wsRef.current?.send("CONNECT " + nickname)
  }

  const handleLogout = () => {
    wsRef.current?.send("QUIT")
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
          <ActionPanel onCommand={sendCommand} inventory={inventoryData} group={groupData} questData={questData} />
        </div>
      </main>
    </>
  )
}

export default App

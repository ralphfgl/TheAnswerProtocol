import './App.css'
import { useEffect, useRef, useState } from "react"
import Header from './header/Header'
import ChatPanel from './chatpanel/ChatPanel'

function App() {
  const [messages, setMessages] = useState([])
  const [inputValue, setInputValue] = useState("")
  const wsRef = useRef(null)

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
  return (
    <>
      <main className='main'>
        <h1 className='title'>The answer Protocol</h1>
        <Header />
        <ChatPanel />
        <br></br>
        <br></br>
        <br></br>
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

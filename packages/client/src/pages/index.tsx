// pages/admin.tsx
import { useEffect, useState } from 'react'
import '@/styles/globals.css'

export default function AdminPage() {
  const [socket, setSocket] = useState<WebSocket | null>(null)
  const [gameID, setGameID] = useState('')
  const [clients, setClients] = useState<string[]>([])
  const [score, setScore] = useState(0)
  const [lives, setLives] = useState(3)
  const [isActive, setIsActive] = useState(false)
  
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws?type=admin`)
    
    ws.onopen = () => {
      console.log('Connected to server')
    }
    
    // TODO: Maybe try to bind the message types with the types in the Go project
    ws.onmessage = (event) => {
      let messages = (event.data as string).split("\n")
      messages.forEach((message => {
        const data = JSON.parse(message)
        
        switch (data.type) {
          case 'game_created':
            setGameID(data.payload.gameID)
            break
            
          case 'client_joined':
            setClients(prev => [...prev, data.payload.playerID])
            break
            
          case 'update_score':
            setScore(data.payload.score)
            break
            
          case 'update_lives':
            setLives(data.payload.lives)
            break
            
          case 'game_state':
            setIsActive(data.payload.isActive)
            setScore(data.payload.score)
            setLives(data.payload.lives)
            break
            
          case 'game_over':
            setIsActive(false)
            setLives(0)
            break
            
          case 'client_disconnected':
            setClients(prev => prev.filter(id => id !== data.payload.playerID))
            break
        }
      }))
    }
    
    setSocket(ws)
    
    return () => {
      ws.close()
    }
  }, [])
  
  const startGame = () => {
    if (socket && clients.length > 0) {
      socket.send(JSON.stringify({
        type: 'start_game',
        payload: {}
      }))
    }
  }
  
  const endGame = () => {
    if (socket) {
      socket.send(JSON.stringify({
        type: 'end_game',
        payload: {}
      }))
    }
  }
  
  // TODO: Add a countdown of the time left
  return (
    <div className="flex flex-col min-h-screen bg-slate-900 text-gray-100 font-sans">

      <header className="bg-gradient-to-r from-blue-800 to-indigo-900 p-6 shadow-lg">
        <h1 className="text-3xl font-bold text-white text-center">Key-bored</h1>
      </header>

      <main className="flex-grow p-6 md:p-8 max-w-4xl mx-auto w-full">

        <div className="bg-slate-800 rounded-lg p-5 mb-6 shadow-md border border-blue-700">
          <p className="text-gray-400 mb-2 text-sm">Click to open other windows to play in</p>
          <a href={'http://localhost:3000/play/' + gameID} target='_blank' className="text-2xl font-mono tracking-wide text-blue-400 font-bold">Click here</a>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
          <div className="bg-slate-800 rounded-lg p-5 shadow-md border border-indigo-700">
            <h2 className="text-xl font-semibold text-indigo-400 mb-4">Game Stats</h2>
            <div className="flex justify-between items-center mb-2">
              <p className="text-gray-300">Score:</p>
              <p className="text-2xl font-bold text-white">{score}</p>
            </div>
            <div className="flex justify-between items-center">
              <p className="text-gray-300">Lives:</p>
              <div className="flex">
                {[...Array(lives)].map((_, i) => (
                  <span key={i} className="text-xl text-red-500 mx-1">❤</span>
                ))}
              </div>
            </div>
          </div>
          
          <div className="bg-slate-800 rounded-lg p-5 shadow-md border border-indigo-700">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold text-indigo-400">Active Sessions</h2>
              <span className="bg-blue-700 text-white px-2 py-1 rounded-full text-xs font-medium">
                {clients.length}
              </span>
            </div>
            <ul className="space-y-2">
              {clients.map(player => (
                <li key={player} className="bg-slate-700 p-2 rounded flex items-center">
                  <span className="w-2 h-2 bg-green-500 rounded-full mr-2"></span>
                  <span className="font-mono text-gray-300">Client {player.substring(0, 5)}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
        
        <div className="flex justify-center mt-4">
          {!isActive ? (
            <button 
              onClick={startGame} 
              disabled={clients.length === 0}
              className={`px-8 py-3 text-lg font-medium rounded-lg shadow-md transition-all duration-200 
                ${clients.length === 0 
                  ? 'bg-gray-600 cursor-not-allowed' 
                  : 'bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-700 hover:to-indigo-700 text-white'}`}
            >
              Start Game
            </button>
          ) : (
            <button 
              onClick={endGame} 
              className="px-8 py-3 text-lg font-medium rounded-lg shadow-md
                bg-gradient-to-r from-red-500 to-red-700 hover:from-red-600 hover:to-red-800 text-white
                transition-all duration-200"
            >
              End Game
            </button>
          )}
        </div>
      </main>
      
      <footer className="bg-slate-800 p-4 border-t border-blue-900">
        <p className="text-center text-gray-500 text-sm">
          &copy; 2025 Key-bored Gaming
        </p>
      </footer>
    </div>
  )
}
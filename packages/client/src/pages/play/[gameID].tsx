// pages/play/[gameID].tsx
import { useEffect, useState } from 'react'
import { useRouter } from 'next/router'
import '@/styles/globals.css'

export default function PlayerPage() {
  const router = useRouter()
  const { gameID } = router.query
  
  const [socket, setSocket] = useState<WebSocket | null>(null)
  const [currentKey, setCurrentKey] = useState('')
  const [isActive, setIsActive] = useState(false)
  const [exploding, setExploding] = useState(false)
  
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws?type=player&gameID=${gameID}`)
    
    ws.onopen = () => {
      console.log('Connected to server')
    }
    
    // TODO: Maybe try to bind the message types with the types in the Go project
    ws.onmessage = (event) => {
      let messages = (event.data as string).split("\n")
      messages.forEach((message) => {
          const data = JSON.parse(message)
          switch (data.type) {
            case 'new_key':
              setExploding(true)
              setCurrentKey(getRandomLowerCaseLetter())
              setTimeout(() => {
                setExploding(false)
              }, 300)
              // TODO: The below imposes a random nuke of the current key
              // These random nukes stack as clients recieve more messages
              // If clients are oversaturated they will get nuked too often
              // I'd like to figure out something a little different
              // Maybe moving the timeout server side 
              setTimeout(() => {
                setCurrentKey('')
              }, 5000 * Math.random() + 0.5)
              break
              
            case 'game_state':
              setIsActive(data.payload.isActive)
              break
              
            case 'game_over':
              setIsActive(false)
              break
          }
      })
    }

    setSocket(ws)

    return () => {
        ws.close()
    }
  }, [gameID])

  useEffect(() => {

    const handleKeyDown = (event: KeyboardEvent) => {
      if (!isActive || !socket) {
        return
      }

      if (currentKey === event.key) {

        socket.send(JSON.stringify({
          type: 'add_score',
          payload: {}
        }))
      } else {
        socket.send(JSON.stringify({
            type: 'lose_life',
            payload: {}
        }))
      }

    }
    
    window.addEventListener('keydown', handleKeyDown)
    
    return () => {
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [isActive, currentKey])
  
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-slate-900 text-gray-100 font-sans">

      <header className="fixed top-0 left-0 right-0 bg-gradient-to-r from-blue-800 to-indigo-900 p-4 shadow-lg">
        <h1 className="text-2xl font-bold text-white text-center">Key-bored</h1>
      </header>
      
      <main className="flex-grow flex items-center justify-center w-full p-6">
        {!isActive ? (
          <div className="bg-slate-800 rounded-lg p-8 shadow-lg border border-blue-700 max-w-md w-full text-center">
            <div className="animate-pulse mb-4">
              <div className="w-4 h-4 bg-blue-500 rounded-full mx-auto mb-4"></div>
              <div className="w-6 h-6 bg-blue-400 rounded-full mx-auto mb-4"></div>
              <div className="w-8 h-8 bg-blue-300 rounded-full mx-auto"></div>
            </div>
            <h2 className="text-xl font-semibold text-indigo-400 mb-2">Waiting for game to start...</h2>
            <p className="text-gray-400 text-sm">The host will begin the game shortly</p>
          </div>
        ) : (
          <div className={`
            relative
            flex items-center justify-center
            w-48 h-48 md:w-64 md:h-64 lg:w-80 lg:h-80
            bg-gradient-to-br from-blue-700 to-indigo-800
            rounded-2xl
            shadow-lg
            border-4 border-blue-500
            cursor-pointer
            transition-all duration-300
          `}>
            {exploding && (
              <div className="absolute inset-0 flex items-center justify-center">
                <div className="absolute w-full h-full bg-yellow-500 rounded-2xl opacity-70 animate-ping"></div>
                <div className="absolute w-10 h-10 bg-white rounded-full animate-ping"></div>
                <div className="absolute w-20 h-20 border-4 border-white rounded-full animate-ping"></div>
                <div className="absolute w-40 h-40 border-2 border-white rounded-full animate-ping"></div>
              </div>
            )}
            
            <div className={`
              text-6xl md:text-8xl lg:text-9xl font-bold text-white
              transform hover:scale-105 transition-transform
            `}>
              {currentKey}
            </div>
          </div>
        )}
      </main>
      
      <footer className="fixed bottom-0 left-0 right-0 bg-slate-800 p-4 border-t border-blue-900">
        <p className="text-center text-gray-400 text-sm">
          {isActive 
            ? "Type the letter as quickly as possible!" 
            : "Get ready to type the letters that appear"}
        </p>
      </footer>
    </div>
  )
}

function getRandomLowerCaseLetter(): string {
    // ASCII code for 'a' is 97
    const randomCharCode = Math.floor(Math.random() * 26) + 97;
    return String.fromCharCode(randomCharCode);
  }
import React, { useState, useEffect } from 'react'

interface Keymap {
  id: string;
  name: string;
}

function App() {
  const [inputText, setInputText] = useState('')
  const [outputText, setOutputText] = useState('')
  const [selectedKeymap, setSelectedKeymap] = useState('')
  const [keymaps, setKeymaps] = useState<Keymap[]>([])
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    // Fetch available keymaps when component mounts
    fetchKeymaps()
  }, [])

  const fetchKeymaps = async () => {
    try {
      const response = await fetch('/api/keymaps')
      const data = await response.json()
      setKeymaps(data)
    } catch (error) {
      console.error('Error fetching keymaps:', error)
    }
  }

  const handleTransliterate = async () => {
    if (!selectedKeymap || !inputText) return

    setIsLoading(true)
    try {
      const response = await fetch('/api/transliterate', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          text: inputText,
          keymapId: selectedKeymap,
        }),
      })
      const data = await response.json()
      setOutputText(data.result)
    } catch (error) {
      console.error('Error during transliteration:', error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto px-4 py-8">
        <header className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">Aksharamala</h1>
          <p className="text-gray-600">Indic Script Transliteration</p>
        </header>

        <div className="bg-white rounded-lg shadow-lg p-6">
          <div className="mb-6">
            <label htmlFor="keymap" className="block text-sm font-medium text-gray-700 mb-2">
              Select Keymap
            </label>
            <select
              id="keymap"
              className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              value={selectedKeymap}
              onChange={(e) => setSelectedKeymap(e.target.value)}
            >
              <option value="">Select a keymap...</option>
              {keymaps.map((keymap) => (
                <option key={keymap.id} value={keymap.id}>
                  {keymap.name}
                </option>
              ))}
            </select>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label htmlFor="input" className="block text-sm font-medium text-gray-700 mb-2">
                Input Text
              </label>
              <textarea
                id="input"
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                rows={6}
                value={inputText}
                onChange={(e) => setInputText(e.target.value)}
                placeholder="Enter text to transliterate..."
              />
            </div>

            <div>
              <label htmlFor="output" className="block text-sm font-medium text-gray-700 mb-2">
                Output Text
              </label>
              <textarea
                id="output"
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm bg-gray-50"
                rows={6}
                value={outputText}
                readOnly
                placeholder="Transliterated text will appear here..."
              />
            </div>
          </div>

          <div className="mt-6 text-center">
            <button
              onClick={handleTransliterate}
              disabled={!selectedKeymap || !inputText || isLoading}
              className={`px-6 py-2 rounded-md text-white font-medium
                ${!selectedKeymap || !inputText || isLoading
                  ? 'bg-gray-400 cursor-not-allowed'
                  : 'bg-blue-600 hover:bg-blue-700'
                }`}
            >
              {isLoading ? 'Transliterating...' : 'Transliterate'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

export default App

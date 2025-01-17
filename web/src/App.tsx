import { useState, useEffect } from 'react'
import Footer from './Footer'

const Navbar = () => {
  return (
    <nav className="bg-blue-600 text-white shadow-md fixed top-0 w-full z-10">
      <div className="container mx-auto flex items-center justify-between px-4 py-2">
        <a href="/" className="text-lg font-bold">Aksharamala</a>
        <div className="space-x-4">
          <a href="/" className="hover:underline">Home</a>
          <a href="/history" className="hover:underline">History</a>
          <a href="/docs" className="hover:underline">Documentation</a>
          <a href="https://github.com/your-repo" target="_blank" rel="noopener noreferrer" className="hover:underline">
            GitHub
          </a>
          <a href="/about" className="hover:underline">About</a>
        </div>
      </div>
    </nav>
  );
};

export { Navbar };

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
      const response = await fetch('/api/m', {
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
      setOutputText(data.result || '')
    } catch (error) {
      console.error('Error during transliteration:', error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen font-['Plus_Jakarta_Sans']">
      <div className="max-w-7xl mx-auto px-4 py-12">
        {/* Header */}
        <header className="text-center mb-16">
          <h1 className="text-6xl font-bold text-primary-900 mb-4 font-['Plus_Jakarta_Sans']">
            Aksharamala
          </h1>
          <p className="text-xl text-primary-700 mb-12 font-['Plus_Jakarta_Sans']">
            Indic Script Transliteration
          </p>

          {/* Help Section */}
          <div className="bg-white/80 backdrop-blur-md border border-primary-200 rounded-2xl p-8 max-w-3xl mx-auto shadow-glass">
            <h2 className="text-2xl font-semibold text-primary-900 mb-4 font-['Plus_Jakarta_Sans']">
              Transliteration Help
            </h2>
            <p className="text-lg text-primary-700 mb-6 font-['Plus_Jakarta_Sans']">
              Type Roman letters to get the corresponding Indic script.
            </p>
            <div className="text-primary-700 p-4 bg-primary-50 rounded-xl">
              <span className="font-medium">Example:</span> Type{' '}
              <code className="px-2 py-1 bg-white rounded-md font-mono text-primary-700">
                namaste
              </code>{' '}
              with ITRANS scheme to get नमस्ते
            </div>
          </div>
        </header>

        {/* Main Content */}
        <div className="bg-white/90 backdrop-blur-md rounded-2xl shadow-xl p-8 mb-8">
          {/* Script Selection */}
          <div className="mb-8">
            <label
              htmlFor="keymap"
              className="block text-xl font-semibold text-primary-900 mb-3 font-['Plus_Jakarta_Sans']"
            >
              Select Script
            </label>
            <select
              id="keymap"
              className="w-full px-4 py-3 text-lg border border-primary-200 rounded-xl 
                       shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 
                       bg-white text-primary-900 appearance-none font-['Plus_Jakarta_Sans']"
              value={selectedKeymap}
              onChange={(e) => setSelectedKeymap(e.target.value)}
            >
              <option value="">Select a script...</option>
              {keymaps.map((keymap) => (
                <option key={keymap.id} value={keymap.id}>
                  {keymap.name}
                </option>
              ))}
            </select>
          </div>

          {/* Transliteration Interface */}
          <div className="grid grid-cols-[1fr,auto,1fr] gap-8 items-stretch relative">
            {/* Input */}
            <div className="space-y-3">
              <label
                htmlFor="input"
                className="block text-xl font-semibold text-primary-900 font-['Plus_Jakarta_Sans']"
              >
                Input Text
              </label>
              <textarea
                id="input"
                className="w-full px-4 py-3 text-lg border border-primary-200 
                         rounded-xl shadow-sm focus:outline-none focus:ring-2 
                         focus:ring-primary-500 bg-white min-h-[200px]
                         font-['Plus_Jakarta_Sans'] leading-relaxed"
                value={inputText}
                onChange={(e) => setInputText(e.target.value)}
                placeholder="Enter text to transliterate..."
              />
            </div>

            {/* Arrow Button */}
            <div className="flex items-center justify-center">
              <button
                onClick={handleTransliterate}
                disabled={!selectedKeymap || !inputText || isLoading}
                className="p-6 rounded-full bg-primary-600 hover:bg-primary-700 
                         disabled:bg-primary-200 transition-all duration-200 
                         transform hover:scale-105 disabled:hover:scale-100
                         shadow-lg disabled:shadow-none"
                title="Transliterate"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-8 w-8 text-white"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2.5}
                    d="M14 5l7 7m0 0l-7 7m7-7H3"
                  />
                </svg>
              </button>
            </div>

            {/* Output */}
            <div className="space-y-3">
              <label
                htmlFor="output"
                className="block text-xl font-semibold text-primary-900 font-['Plus_Jakarta_Sans']"
              >
                Output Text
              </label>
              <textarea
                id="output"
                className="w-full px-4 py-3 text-lg border border-primary-200 
                         rounded-xl shadow-sm bg-primary-50 font-devanagari
                         min-h-[200px] leading-relaxed"
                value={outputText}
                readOnly
                placeholder="Transliterated text will appear here..."
              />
            </div>
          </div>
        </div>
        <Footer />
      </div>
    </div>
  )
}

export default App

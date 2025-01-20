import { useState, useEffect } from 'react';
import { useThemeContext } from './ThemeProvider';
import Footer from './Footer';
import Navbar from './Navbar';
import { IconButton } from '@mui/material';
import { ArrowForward } from '@mui/icons-material';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081';

interface Keymap {
  id: string;
  name: string;
}

function App() {
  const [inputText, setInputText] = useState('');
  const [outputText, setOutputText] = useState('');
  const [selectedKeymap, setSelectedKeymap] = useState('');
  const [keymaps, setKeymaps] = useState<Keymap[]>([]);
  const { darkMode } = useThemeContext();

  useEffect(() => {
    fetchKeymaps();
  }, []);

  const fetchKeymaps = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/keymaps`);
      const data: Keymap[] = await response.json();
      setKeymaps(data);
    } catch (error) {
      console.error('Error fetching keymaps:', error);
    }
  };

  const handleTransliterate = async () => {
    if (!selectedKeymap || !inputText) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/m`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text: inputText, keymapId: selectedKeymap }),
      });

      const data = await response.json();
      setOutputText(data.result || '');
    } catch (error) {
      console.error('Error during transliteration:', error);
    }
  };

  return (
    <div className={`min-h-screen flex flex-col ${darkMode ? 'bg-[#121212] text-gray-100' : 'bg-gray-50 text-gray-900'}`}>
      <Navbar />
      <main className="flex-grow flex items-center justify-center px-4 py-8">
        <div className="w-full max-w-3xl space-y-8">
          {/* Help Section */}
          <section className={`rounded-2xl p-6 ${darkMode ? 'bg-[#1e1e1e]' : 'bg-white'} shadow-lg`}>
            <h2 className={`text-2xl font-semibold mb-4 ${darkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              Transliteration Help
            </h2>
            <p className={`text-lg mb-6 ${darkMode ? 'text-gray-300' : 'text-gray-700'}`}>
              Type Roman letters to get the corresponding Indic script.
            </p>
            <div className={`p-4 rounded-xl ${darkMode ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-700'}`}>
              <span className="font-medium">Example:</span>{' '}
              Type{' '}
              <code className={`px-2 py-1 rounded-md font-mono ${darkMode ? 'bg-gray-800' : 'bg-white'}`}>
                namaste
              </code>{' '}
              with ITRANS scheme to get नमस्ते
            </div>
          </section>

          {/* Script Selection */}
          <section className={`rounded-2xl p-6 ${darkMode ? 'bg-[#1e1e1e]' : 'bg-white'} shadow-lg`}>
            <label htmlFor="keymap" className={`block text-xl font-semibold mb-3 ${darkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              Select Script
            </label>
            <select
              id="keymap"
              className={`w-full px-4 py-3 text-lg rounded-xl shadow-sm focus:ring-2 focus:ring-primary-500 focus:outline-none ${
                darkMode ? 'bg-gray-700 text-gray-100 border-gray-600' : 'bg-white text-gray-900 border-gray-300'
              }`}
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
          </section>

          {/* Transliteration Interface */}
          <section className={`rounded-2xl p-6 ${darkMode ? 'bg-[#1e1e1e]' : 'bg-white'} shadow-lg`}>
            <label htmlFor="inputText" className={`block text-xl font-semibold mb-3 ${darkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              Input Text
            </label>
            <div className="flex items-center space-x-4">
              <textarea
                id="inputText"
                className={`w-full px-4 py-3 text-lg rounded-xl shadow-sm focus:ring-2 focus:ring-primary-500 focus:outline-none ${
                  darkMode ? 'bg-gray-700 text-gray-100 border-gray-600' : 'bg-white text-gray-900 border-gray-300'
                }`}
                value={inputText}
                onChange={(e) => setInputText(e.target.value)}
                placeholder="Enter text to transliterate..."
                rows={4}
              />
              <IconButton
                color="primary"
                onClick={handleTransliterate}
                disabled={!selectedKeymap || !inputText}
                aria-label="Transliterate"
                className={`p-3 text-2xl shadow-md ${
                  darkMode ? 'bg-primary-600 hover:bg-primary-700' : 'bg-primary-500 hover:bg-primary-600'
                }`}
              >
                <ArrowForward fontSize="large" />
              </IconButton>
            </div>

            <label htmlFor="outputText" className={`block text-xl font-semibold mt-6 mb-3 ${darkMode ? 'text-gray-100' : 'text-gray-900'}`}>
              Output Text
            </label>
            <textarea
              id="outputText"
              className={`w-full px-4 py-3 text-lg rounded-xl shadow-sm focus:outline-none ${
                darkMode ? 'bg-gray-700 text-gray-100 border-gray-600' : 'bg-gray-50 text-gray-900 border-gray-300'
              }`}
              value={outputText}
              readOnly
              placeholder="Transliterated text will appear here..."
              rows={4}
            />
          </section>
        </div>
      </main>
      <Footer />
    </div>
  );
}

export default App;

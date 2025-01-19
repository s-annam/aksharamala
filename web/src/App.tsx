import { useState, useEffect } from 'react';
import Footer from './Footer';
import Navbar from './Navbar';
import { IconButton } from '@mui/material';
import { ArrowForward } from '@mui/icons-material';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081';

// Define Keymap type
interface Keymap {
  id: string;
  name: string;
}

function App() {
  const [inputText, setInputText] = useState('');
  const [outputText, setOutputText] = useState('');
  const [selectedKeymap, setSelectedKeymap] = useState('');
  const [keymaps, setKeymaps] = useState<Keymap[]>([]);

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
    <div className="min-h-screen font-sans flex flex-col">
      {/* Navbar */}
      <Navbar />

      {/* Separator Below Navbar */}
      <div className="w-full border-b border-gray-300 dark:border-gray-700"></div>

      {/* Main Content */}
      <div className="flex-grow max-w-7xl mx-auto px-4 py-12">

        {/* Transliteration UI */}
        <div className="max-w-7xl mx-auto px-4 py-12">
          {/* Transliteration Help Section */}
          <div className="transliteration-help mx-auto max-w-3xl">
            <h2 className="text-2xl font-semibold text-primary-900 mb-4">Transliteration Help</h2>
            <p className="text-lg text-primary-700 mb-6">
              Type Roman letters to get the corresponding Indic script.
          </p>
          <div className="text-primary-700 p-4 rounded-xl">
            <span className="font-medium">Example:</span> Type{' '}
            <code className="px-2 py-1 bg-white rounded-md font-mono text-primary-700">
              namaste
            </code>{' '}
            with ITRANS scheme to get नमस्ते
          </div>
        </div>

        {/* Select Script Section */}
        <div className="transliteration-help mx-auto mt-6 max-w-3xl">
          <label htmlFor="keymap" className="block text-xl font-semibold text-primary-900 mb-3">
            Select Script
          </label>
          <select
            id="keymap"
            className="w-full px-4 py-3 text-lg border rounded-xl shadow-sm focus:outline-none"
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
        <div className="transliteration-help mx-auto mt-6 max-w-3xl">
          <label htmlFor="inputText" className="block text-xl font-semibold text-primary-900 mb-3">
            Input Text
          </label>
          <div className="flex items-center space-x-4">
            <textarea
              id="inputText"
              className="w-full px-4 py-3 text-lg border rounded-xl shadow-sm focus:outline-none"
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder="Enter text to transliterate..."
            />
          {/* Transliteration Button */}
          <IconButton
            color="primary"
            onClick={handleTransliterate}
            disabled={!selectedKeymap || !inputText}
            aria-label="Transliterate"
            className="p-3 text-2xl shadow-md"
          >
            <ArrowForward fontSize="large" />
          </IconButton>
          </div>
          <label htmlFor="outputText" className="block text-xl font-semibold text-primary-900 mt-6 mb-3">
            Output Text
          </label>
          <textarea
            id="outputText"
            className="w-full px-4 py-3 text-lg border rounded-xl shadow-sm focus:outline-none bg-gray-100"
            value={outputText}
            readOnly
            placeholder="Transliterated text will appear here..."
          />
        </div>
      </div>
    </div>

    {/* Separator Above Footer */}
    <div className="w-full border-t border-gray-300 dark:border-gray-700"></div>

    {/* Footer */}
    <Footer />
  </div>
  );
}

export default App;

import { useState } from 'react';
import { useThemeContext } from './ThemeProvider';
import Footer from './Footer';
import Navbar from './Navbar';
import { 
  Box,
  Container,
  TextField,
  Button,
  Paper,
} from '@mui/material';
import { languages } from './config/transliterationConfig';
import LanguageSelector from './components/LanguageSelector';
import type { TransliterationScheme } from './config/transliterationConfig';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081';

function App() {
  const [inputText, setInputText] = useState('');
  const [outputText, setOutputText] = useState('');
  const [selectedScheme, setSelectedScheme] = useState<TransliterationScheme | null>(null);
  const { darkMode } = useThemeContext();

  const handleTransliterate = async () => {
    if (!selectedScheme || !inputText) return;

    try {
      const response = await fetch(`${API_BASE_URL}/api/m`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          text: inputText, 
          keymapId: selectedScheme.id 
        }),
      });

      const data = await response.json();
      setOutputText(data.result || '');
    } catch (error) {
      console.error('Error during transliteration:', error);
    }
  };

  return (
    <Box sx={{ 
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
      bgcolor: darkMode ? 'grey.900' : 'grey.50'
    }}>
      <Navbar />
      
      <Container component="main" sx={{ flex: 1, py: 4 }}>
        <LanguageSelector 
          languages={languages}
          onSchemeSelect={setSelectedScheme}
        />

        <Paper elevation={2} sx={{ p: 3 }}>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
            <TextField
              fullWidth
              multiline
              rows={4}
              label="Input Text"
              value={inputText}
              onChange={(e) => setInputText(e.target.value)}
              placeholder={selectedScheme ? `Type using ${selectedScheme.name} scheme...` : 'Select a transliteration scheme...'}
              disabled={!selectedScheme}
            />

            <Button
              variant="contained"
              onClick={handleTransliterate}
              disabled={!selectedScheme || !inputText}
              sx={{ alignSelf: 'flex-end' }}
            >
              Transliterate
            </Button>

            <TextField
              fullWidth
              multiline
              rows={4}
              label="Output Text"
              value={outputText}
              InputProps={{ readOnly: true }}
            />
          </Box>
        </Paper>
      </Container>

      <Footer />
    </Box>
  );
}

export default App;

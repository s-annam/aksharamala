import { useState } from 'react';
import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Tooltip,
  Typography,
  FormControlLabel,
  Switch,
  Paper,
} from '@mui/material';
import InfoIcon from '@mui/icons-material/Info';
import { Language, TransliterationScheme } from '../config/transliterationConfig';

interface LanguageSelectorProps {
  languages: Language[];
  onSchemeSelect: (scheme: TransliterationScheme) => void;
}

export default function LanguageSelector({ languages, onSchemeSelect }: LanguageSelectorProps) {
  const [selectedLang, setSelectedLang] = useState<string>('');
  const [selectedScheme, setSelectedScheme] = useState<string>('');
  const [showAdvanced, setShowAdvanced] = useState(false);

  const currentLanguage = languages.find(lang => lang.id === selectedLang);
  const availableSchemes = currentLanguage?.schemes.filter(
    scheme => !scheme.isAdvanced || showAdvanced
  ) || [];

  const handleLanguageChange = (langId: string) => {
    setSelectedLang(langId);
    setSelectedScheme('');
  };

  const handleSchemeChange = (schemeId: string) => {
    setSelectedScheme(schemeId);
    const scheme = currentLanguage?.schemes.find(s => s.id === schemeId);
    if (scheme) {
      onSchemeSelect(scheme);
    }
  };

  return (
    <Paper elevation={2} sx={{ p: 3, mb: 3 }}>
      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
        <FormControl fullWidth>
          <InputLabel>Select Language</InputLabel>
          <Select
            value={selectedLang}
            label="Select Language"
            onChange={(e) => handleLanguageChange(e.target.value)}
          >
            {languages.map((lang) => (
              <MenuItem key={lang.id} value={lang.id}>
                {lang.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        {currentLanguage && (
          <>
            <FormControl fullWidth>
              <InputLabel>Select Transliteration Scheme</InputLabel>
              <Select
                value={selectedScheme}
                label="Select Transliteration Scheme"
                onChange={(e) => handleSchemeChange(e.target.value)}
              >
                {availableSchemes.map((scheme) => (
                  <MenuItem key={scheme.id} value={scheme.id}>
                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                      {scheme.name}
                      <Tooltip title={scheme.description} arrow>
                        <InfoIcon sx={{ ml: 1, fontSize: '1rem', opacity: 0.7 }} />
                      </Tooltip>
                    </Box>
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            {currentLanguage.schemes.some(s => s.isAdvanced) && (
              <FormControlLabel
                control={
                  <Switch
                    checked={showAdvanced}
                    onChange={(e) => setShowAdvanced(e.target.checked)}
                  />
                }
                label="Show Advanced Options"
              />
            )}

            {selectedScheme && (
              <Box sx={{ mt: 1 }}>
                {currentLanguage.schemes
                  .find(s => s.id === selectedScheme)
                  ?.helpTips?.map((tip, index) => (
                    <Typography
                      key={index}
                      variant="body2"
                      color="text.secondary"
                      sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}
                    >
                      <InfoIcon sx={{ fontSize: '1rem' }} />
                      {tip}
                    </Typography>
                  ))}
              </Box>
            )}
          </>
        )}
      </Box>
    </Paper>
  );
}

import { AppBar, Toolbar, Typography, IconButton, Tooltip, Button } from '@mui/material';
import { Brightness4, Brightness7, GitHub } from '@mui/icons-material';
import { useThemeContext } from './ThemeProvider';

const Navbar = () => {
  const { darkMode, toggleTheme } = useThemeContext();

  return (
    <AppBar 
      position="static" 
      elevation={0} 
      sx={{ 
        backgroundColor: darkMode ? '#121212' : '#ffffff',
        color: darkMode ? '#f0f0f0' : '#222',
        borderBottom: darkMode ? '1px solid rgba(255, 255, 255, 0.12)' : '1px solid rgba(0, 0, 0, 0.12)'
      }}
    >
      <Toolbar className="max-w-3xl mx-auto w-full flex justify-between items-center px-4">
        <div className="flex items-center space-x-3">
          <img src="/og-image.png" alt="Aksharamala Logo" className="h-8 w-auto" />
          <div>
            <Typography variant="h6" className="font-semibold text-lg">
              Aksharamala
            </Typography>
            <Typography variant="caption" className={`hidden sm:block ${darkMode ? 'text-gray-400' : 'text-gray-600'}`}>
              Open Source Indic Script Tools
            </Typography>
          </div>
        </div>

        <div className="flex items-center space-x-2">
          <Button
            href="https://github.com/s-annam/aksharamala"
            target="_blank"
            rel="noopener noreferrer"
            startIcon={<GitHub />}
            sx={{ 
              color: darkMode ? '#f0f0f0' : '#222',
              '&:hover': {
                backgroundColor: darkMode ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.04)'
              }
            }}
            className="hidden sm:flex"
          >
            Contribute
          </Button>
          
          <Tooltip title="Toggle Theme">
            <IconButton 
              onClick={toggleTheme} 
              sx={{ 
                color: darkMode ? '#f0f0f0' : '#222'
              }}
            >
              {darkMode ? <Brightness7 /> : <Brightness4 />}
            </IconButton>
          </Tooltip>
        </div>
      </Toolbar>
    </AppBar>
  );
};

export default Navbar;

import { AppBar, Toolbar, Typography, IconButton, Tooltip } from '@mui/material';
import { Brightness4, Brightness7 } from '@mui/icons-material';
import { useThemeContext } from './ThemeProvider';

const Navbar = () => {
  const { darkMode, toggleTheme } = useThemeContext();

  return (
    <AppBar position="static" elevation={0} 
      sx={{ 
        backgroundColor: darkMode ? '#1A1A1A' : '#ffffff', 
        color: darkMode ? '#f0f0f0' : '#222',
        borderBottom: darkMode ? '1px solid #444' : '1px solid #ddd'
        }}
    >
      <Toolbar className="max-w-7xl mx-auto flex justify-between items-center">
        <div className="flex items-center space-x-3">
          <img src="/og-image.png" alt="Aksharamala Logo" className="h-8 w-auto" />
          <Typography variant="h6" className="font-bold text-lg">
            Aksharamala
          </Typography>
        </div>

      <Tooltip title="Toggle Theme">
        <IconButton onClick={toggleTheme} color="inherit">
          {darkMode ? <Brightness7 /> : <Brightness4 />}
        </IconButton>
      </Tooltip>
      </Toolbar>
    </AppBar>
  );
};

export default Navbar;

import { Typography, IconButton } from '@mui/material';
import { Settings } from '@mui/icons-material';

const Navbar = () => {
  return (
    <nav className="w-full border-b border-gray-300 dark:border-gray-700 px-4 py-3 flex justify-between items-center">
      {/* Logo & Title (Left-Aligned) */}
      <div className="flex items-center space-x-3">
        <img src="/og-image.png" alt="Aksharamala Logo" className="h-8 w-auto" />
        <Typography variant="h6" className="font-bold text-lg">
          Aksharamala
        </Typography>
      </div>

      {/* Future User Links (Right-Aligned) */}
      <div className="flex items-center space-x-3">
        <IconButton color="inherit" aria-label="Settings">
          <Settings />
        </IconButton>
      </div>
    </nav>
  );
};

export default Navbar;

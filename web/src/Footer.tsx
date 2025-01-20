import { useThemeContext } from './ThemeProvider';

const Footer = () => {
  const { darkMode } = useThemeContext();

  return (
    <footer className={`w-full py-6 ${
      darkMode ? 'bg-[#121212] text-gray-400 border-gray-800' : 'bg-white text-gray-600 border-gray-200'
    } border-t`}>
      <div className="max-w-3xl mx-auto px-6 flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <p className="text-sm">
          © {new Date().getFullYear()} Aksharamala. All Rights Reserved.
        </p>

        <div className="flex items-center gap-6 text-sm">
          <a 
            href="https://github.com/s-annam/aksharamala/blob/main/HISTORY.md" 
            target="_blank" 
            rel="noopener noreferrer" 
            className={`hover:text-primary-500 transition-colors ${
              darkMode ? 'hover:text-primary-400' : 'hover:text-primary-600'
            }`}
          >
            Project History
          </a>
          <a 
            href="https://github.com/s-annam/aksharamala/blob/main/docs/CONTRIBUTING.md" 
            target="_blank" 
            rel="noopener noreferrer" 
            className={`hover:text-primary-500 transition-colors flex items-center gap-1 ${
              darkMode ? 'hover:text-primary-400' : 'hover:text-primary-600'
            }`}
          >
            <span>Join Us</span>
            <span className="inline-block animate-pulse">→</span>
          </a>
        </div>
      </div>
    </footer>
  );
};

export default Footer;

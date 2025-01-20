import { useThemeContext } from './ThemeProvider';

const Footer = () => {
  const { darkMode } = useThemeContext();

  return (
    <footer className={`w-full px-6 py-4 flex justify-center items-center text-sm border-t ${
      darkMode ? 'bg-[#1A1A1A] text-gray-400 border-gray-700' : 'bg-gray-100 text-gray-700 border-gray-300'
    }`}>
      <p>© {new Date().getFullYear()} Aksharamala. All Rights Reserved.</p>

      {/* External Links */}
      <div className="flex space-x-3">
        <a 
          href="https://github.com/s-annam/aksharamala/blob/main/HISTORY.md" 
          target="_blank" 
          rel="noopener noreferrer" 
          className="hover:underline"
        >
          History
        </a>
        <span>|</span>
        <a 
          href="https://github.com/s-annam/aksharamala/blob/main/docs/CONTRIBUTING.md" 
          target="_blank" 
          rel="noopener noreferrer" 
          className="hover:underline"
        >
          Contribute
        </a>
      </div>
    </footer>
  );
};

export default Footer;

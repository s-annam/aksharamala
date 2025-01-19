const Footer = () => {
  return (
    <footer className="w-full border-t border-gray-300 dark:border-gray-700 px-6 py-4 flex justify-between items-center text-sm text-gray-700 dark:text-gray-300">
      {/* Logo */}
      <img src="/og-image.png" alt="Aksharamala Logo" className="h-5 w-auto opacity-75" />

      {/* Copyright & Links (Centered) */}
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

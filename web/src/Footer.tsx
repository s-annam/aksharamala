const Footer = () => {
  return (
<footer className="bg-gradient-to-br from-primary-50 to-primary-100 text-primary-900 py-4 text-center">
  <div className="max-w-7xl mx-auto px-4">
    {/* Footer Links Section */}
    <div className="mb-3 text-lg font-medium">
      <a 
        href="https://github.com/s-annam/aksharamala/blob/main/HISTORY.md" 
        target="_blank" 
        rel="noopener noreferrer" 
        className="text-primary-700 hover:text-primary-900 underline decoration-2 underline-offset-2 mx-2"
      >
        History
      </a>
      {' | '}
      <a 
        href="https://github.com/s-annam/aksharamala/blob/main/CONTRIBUTING.md" 
        target="_blank" 
        rel="noopener noreferrer" 
        className="text-primary-700 hover:text-primary-900 underline decoration-2 underline-offset-2 mx-2"
      >
        How to Contribute?
      </a>
    </div>

    {/* Copyright */}
    <p className="text-primary-800 text-sm">
      © {new Date().getFullYear()} Aksharamala. All Rights Reserved.
    </p>
  </div>
</footer>
  );
};

export default Footer;

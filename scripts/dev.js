const { spawn, exec } = require('child_process');
const process = require('process');
const path = require('path');

// Helper to kill processes by port
function killProcessOnPort(port) {
    const cmd = process.platform === 'win32' 
        ? `netstat -ano | findstr :${port}`
        : `lsof -i :${port} -t`;
    
    exec(cmd, (error, stdout) => {
        if (error) return;
        
        stdout.split('\n').forEach(line => {
            const pid = process.platform === 'win32'
                ? line.split(/\s+/).pop()
                : line.trim();
            
            if (pid) {
                try {
                    process.platform === 'win32'
                        ? exec(`taskkill /F /PID ${pid}`)
                        : exec(`kill -9 ${pid}`);
                } catch (e) {
                    // Ignore errors if process is already gone
                }
            }
        });
    });
}

// Kill existing processes
killProcessOnPort(8081); // Go server
killProcessOnPort(5173); // Vite dev server

// Wait a bit for processes to be killed
setTimeout(() => {
    // Start Go server
    const goServer = spawn('go', ['run', 'cmd/webserver/main.go'], {
        cwd: path.join(__dirname, '..'),
        stdio: 'inherit',
        shell: true
    });

    goServer.on('error', (err) => {
        console.error('Failed to start Go server:', err);
    });

    // Start frontend
    const frontend = spawn('npm', ['start'], {
        cwd: path.join(__dirname, '..', 'web'),
        stdio: 'inherit',
        shell: true
    });

    frontend.on('error', (err) => {
        console.error('Failed to start frontend:', err);
    });

    // Handle cleanup on exit
    const cleanup = () => {
        goServer.kill();
        frontend.kill();
        process.exit();
    };

    process.on('SIGINT', cleanup);
    process.on('SIGTERM', cleanup);

    console.log('\nDevelopment servers started!');
    console.log('Frontend: http://localhost:5173');
    console.log('Backend: http://localhost:8081\n');
    console.log('Press Ctrl+C to stop both servers.\n');
}, 1000);

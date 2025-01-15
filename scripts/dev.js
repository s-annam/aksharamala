const { spawn, exec } = require('child_process');
const process = require('process');
const path = require('path');

// Helper to check if port is in use
async function checkPort(port) {
    return new Promise((resolve) => {
        const cmd = process.platform === 'win32'
            ? `netstat -ano | findstr :${port}`
            : `lsof -i :${port} -t`;

        exec(cmd, (error, stdout) => {
            if (error || !stdout) {
                resolve(false);
                return;
            }
            resolve(true);
        });
    });
}

// Helper to kill processes by port
async function killProcessOnPort(port) {
    const cmd = process.platform === 'win32'
        ? `netstat -ano | findstr :${port}`
        : `lsof -i :${port} -t`;

    return new Promise((resolve) => {
        exec(cmd, async (error, stdout) => {
            if (error || !stdout) {
                resolve();
                return;
            }

            const pids = stdout.split('\n')
                .map(line => process.platform === 'win32'
                    ? line.split(/\s+/).pop()
                    : line.trim())
                .filter(Boolean);

            for (const pid of pids) {
                try {
                    process.platform === 'win32'
                        ? exec(`taskkill /F /PID ${pid}`)
                        : exec(`kill -9 ${pid}`);
                } catch (e) {
                    // Ignore errors if process is already gone
                }
            }

            // Wait and verify
            await new Promise(resolve => setTimeout(resolve, 1000));
            
            // Check if port is still in use
            const stillInUse = await checkPort(port);
            if (stillInUse) {
                console.log(`Port ${port} is still in use. Attempting to kill again...`);
                await killProcessOnPort(port); // Recursive call
            }

            resolve();
        });
    });
}

async function main() {
    console.log('Stopping any existing processes...');

    // Kill existing processes
    await killProcessOnPort(8081); // Go server
    await killProcessOnPort(5173); // Vite dev server
    await killProcessOnPort(5174); // Backup Vite port

    // Verify ports are free
    const port5173InUse = await checkPort(5173);
    const port5174InUse = await checkPort(5174);
    const port8081InUse = await checkPort(8081);

    if (port5173InUse || port5174InUse || port8081InUse) {
        console.error('Unable to free required ports. Please check running processes manually.');
        process.exit(1);
    }

    console.log('Starting development servers...');

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
}

main().catch(console.error);

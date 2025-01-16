const { spawn, exec } = require('child_process');
const process = require('process');
const path = require('path');
const tcpPortUsed = require(path.resolve(__dirname, '../web/node_modules/tcp-port-used'));

// Helper to check if port is in use
async function checkPort(port) {
    return new Promise((resolve) => {
        const cmd = process.platform === 'win32'
            ? `powershell.exe -NoProfile -Command "Get-NetTCPConnection -LocalPort ${port} -ErrorAction SilentlyContinue"`
            : `lsof -i :${port} -t`;

        exec(cmd, (error, stdout, stderr) => {
            if (error) {
                console.error(`Error executing ${cmd}: ${error.message}`);
                console.error(stderr);
                resolve(false);
                return;
            }
            resolve(!!stdout.trim());
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
            if (error) {
                console.error(`Error executing ${cmd}:`, error.message);
                resolve(false);
                return;
            }

            const pids = stdout.split('\n')
                .map(line => line.trim().split(/\s+/).pop()) // Extract the PID
                .filter(pid => !isNaN(pid));

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
            const stillInUse = await checkPort(port);
            if (stillInUse) {
                console.error(`Failed to free port ${port}.`);
            }

            resolve();
        });
    });
}

async function main() {
    console.log('Stopping any existing processes...');

    // Kill existing processes
    await killProcessOnPort(8080);
    await killProcessOnPort(5173);
    await killProcessOnPort(5174);

    // Verify ports are free
    const port5173InUse = await checkPort(5173);
    const port5174InUse = await checkPort(5174);
    const port8080InUse = await checkPort(8080);

    if (port5173InUse || port5174InUse || port8080InUse) {
        console.error('Unable to free required ports. Please check running processes manually.');
        process.exit(1);
    }

    console.log('Starting development servers...');

    const goServer = spawn('go', ['run', 'cmd/webserver/main.go'], {
        cwd: path.join(__dirname, '..'),
        stdio: 'inherit',
        shell: true,
    });

    goServer.on('error', (err) => {
        console.error('Failed to start Go server:', err);
    });

    const frontend = spawn('npm', ['start'], {
        cwd: path.join(__dirname, '..', 'web'),
        stdio: 'inherit',
        shell: true,
    });

    frontend.on('error', (err) => {
        console.error('Failed to start frontend:', err);
    });

    const cleanup = () => {
        if (goServer && !goServer.killed) goServer.kill();
        if (frontend && !frontend.killed) frontend.kill();
        process.exit();
    };

    process.on('exit', cleanup);
    process.on('SIGINT', cleanup);
    process.on('SIGTERM', cleanup);
    process.on('uncaughtException', (err) => {
        console.error('Uncaught exception:', err);
        cleanup();
    });

    console.log('\nDevelopment servers started!');
    console.log('Frontend: http://localhost:5173');
    console.log('Backend: http://localhost:8080\n');
    console.log('Press Ctrl+C to stop both servers.\n');
}

main().catch(console.error);
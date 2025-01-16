const { spawn, exec } = require('child_process');
const process = require('process');
const path = require('path');
const tcpPortUsed = require(path.resolve(__dirname, '../web/node_modules/tcp-port-used'));

// Helper to check if port is in use
async function checkPort(port) {
    return new Promise((resolve) => {
        const cmd = process.platform === 'win32'
            ? `powershell.exe -NoProfile -Command "Get-NetTCPConnection -LocalPort ${port} | Select-Object -ExpandProperty OwningProcess"`
            : `lsof -i :${port} -t`;

        exec(cmd, (error, stdout, stderr) => {
            if (error) {
                console.error(`Error executing ${cmd}: ${error.message}`);
                console.error(stderr);
                resolve(false);
                return;
            }

            const pids = stdout
                .split('\n')
                .map((line) => line.trim())
                .filter((line) => line && !isNaN(line)); // Ensure we only keep numeric values

            if (pids.length > 0) {
                console.log(`Port ${port} is occupied by PID(s): ${pids.join(', ')}`);
                resolve(true);
            } else {
                console.log(`Port ${port} is free.`);
                resolve(false);
            }
        });
    });
}

// Helper to kill processes by port
async function killProcessOnPort(port) {
    return new Promise((resolve) => {
        const cmd = process.platform === 'win32'
            ? `powershell.exe -NoProfile -Command "Get-NetTCPConnection -LocalPort ${port} | ForEach-Object {Stop-Process -Id $_.OwningProcess -Force}"`
            : `lsof -i :${port} -t | xargs kill -9`;

        exec(cmd, async (error, stdout, stderr) => {
            if (error) {
                console.error(`Error executing ${cmd}: ${error.message}`);
                console.error(stderr);
                resolve(false);
                return;
            }

            console.log(`Successfully freed port ${port}.`);

            // Double-check if the port is still in use
            setTimeout(async () => {
                const stillInUse = await checkPort(port);
                if (stillInUse) {
                    console.error(`Port ${port} is still occupied! Attempting to kill again...`);
                    await killProcessOnPort(port);
                } else {
                    console.log(`Port ${port} is now free.`);
                }
                resolve();
            }, 1000);
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

    const cleanup = async () => {
        console.log('Cleaning up...');
    
        await killProcessOnPort(8080);
        await killProcessOnPort(5173);
        await killProcessOnPort(5174);
    
        console.log('Cleanup complete. Exiting...');
        process.exit();
    };
    
    process.on('exit', cleanup);
    process.on('SIGINT', async () => {
        console.log('Received SIGINT (Ctrl+C), performing cleanup...');
        await cleanup();
    });
    process.on('SIGTERM', async () => {
        console.log('Received SIGTERM, performing cleanup...');
        await cleanup();
    });
    process.on('uncaughtException', async (err) => {
        console.error('Uncaught exception:', err);
        await cleanup();
    });

    console.log('\nDevelopment servers started!');
    console.log('Frontend: http://localhost:5173');
    console.log('Backend: http://localhost:8080\n');
    console.log('Press Ctrl+C to stop both servers.\n');
}

main().catch(console.error);
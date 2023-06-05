/**
 * SPDX-License-Identifier: Apache-2.0
 * Copyright 2023 Authors of KubeArmor
 */

const childProcess = require('child_process')
const os = require('os')
const process = require('process')

// Get current git commit hash
const VERSION = childProcess.execSync('git rev-parse HEAD').toString().trim()

// Output directory
const OUTPUT_DIR = `${__dirname}/_output/bin`

// Choose binary based on platform and architecture
function chooseBinary() {
    const platform = os.platform()
    const arch = os.arch()

    if (platform === 'linux' && arch === 'amd64') {
        return `linux-amd64-${VERSION}`
    }
    if (platform === 'linux' && arch === 'arm64') {
        return `linux-arm64-${VERSION}`
    }
    if (platform === 'windows' && arch === 'amd64') {
        return `windows-amd64-${VERSION}`
    }
    if (platform === 'windows' && arch === 'arm64') {
        return `windows-arm64-${VERSION}`
    }

    console.error(`Unsupported platform (${platform}) and architecture (${arch})`)
    process.exit(1)
}

function main() {
    const binary = chooseBinary()
    const mainScript = `${OUTPUT_DIR}/${binary}`
    const spawnSyncReturns = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
    const status = spawnSyncReturns.status
    if (typeof status === 'number') {
        process.exit(status)
    }
    process.exit(1)
}

if (require.main === module) {
    main()
}
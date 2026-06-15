const CONFIG = {
    API_URL: window.location.hostname === '127.0.0.1' || window.location.hostname === 'localhost' 
        ? 'http://localhost:8080/api' 
        : 'https://your-live-backend-domain.com/api',
    WS_URL: window.location.hostname === '127.0.0.1' || window.location.hostname === 'localhost' 
        ? 'ws://localhost:8080/ws' 
        : 'wss://your-live-backend-domain.com/ws'
};

let ws = null;
let username = "";
let score = 0;
let hits = 0;
let misses = 0;
let targetSpawnTime = 0;
let totalReactionTime = 0;
let gameInterval = null;
let targetTimeout = null;
let gameLength = 15000; // 15 seconds match execution

// UI Elements
const setupScreen = document.getElementById('setup-screen');
const countdownScreen = document.getElementById('countdown-screen');
const arenaScreen = document.getElementById('arena-screen');
const usernameInput = document.getElementById('username-input');
const startBtn = document.getElementById('start-btn');
const countdownNumber = document.getElementById('countdown-number');
const scoreVal = document.getElementById('score-val');
const hitsVal = document.getElementById('hits-val');
const missesVal = document.getElementById('misses-val');
const targetContainer = document.getElementById('target-container');

// Core Initialization
document.addEventListener('DOMContentLoaded', () => {
    connectWebSocket();
    fetchInitialLeaderboard();
    startBtn.addEventListener('click', prepareGame);
    targetContainer.addEventListener('click', handleArenaClick);
});

async function fetchInitialLeaderboard() {
    try {
        const response = await fetch(`${CONFIG.API_URL}/leaderboard/top?limit=10`);
        if (!response.ok) throw new Error('Network fault');
        const scores = await response.json();
        updateLeaderboardUI(scores);
    } catch (err) {
        console.error('Failed to sync initial rankings:', err);
    }
}

function connectWebSocket() {
    ws = new WebSocket(CONFIG.WS_URL);
    ws.onmessage = (event) => {
        try {
            const message = JSON.parse(event.data);
            if (message.event === 'leaderboard_update' && Array.isArray(message.data)) {
                updateLeaderboardUI(message.data);
            }
        } catch (err) {
            console.error('Error parsing live update packet:', err);
        }
    };
    ws.onclose = () => setTimeout(connectWebSocket, 3000);
}

function updateLeaderboardUI(scores) {
    const container = document.getElementById('leaderboard-display');
    if (!container) return;
    container.innerHTML = '';
    if (!scores || scores.length === 0) {
        container.innerHTML = '<div class="no-scores">No records registered</div>';
        return;
    }
    scores.forEach((entry, idx) => {
        const row = document.createElement('div');
        row.className = 'leaderboard-row';
        row.innerHTML = `
            <span>#${idx + 1}</span>
            <span>${escapeHTML(entry.player_name)}</span>
            <span>${entry.score}</span>
            <span>${entry.accuracy}%</span>
            <span>${entry.avg_reaction_time}ms</span>
        `;
        container.appendChild(row);
    });
}

function prepareGame() {
    username = usernameInput.value.trim();
    if (!username) return alert("Please enter a valid handle");
    setupScreen.classList.add('hidden');
    countdownScreen.classList.remove('hidden');
    
    let count = 3;
    countdownNumber.innerText = count;
    const countInterval = setInterval(() => {
        count--;
        if (count === 0) {
            clearInterval(countInterval);
            countdownScreen.classList.add('hidden');
            arenaScreen.classList.remove('hidden');
            startGame();
        } else {
            countdownNumber.innerText = count;
        }
    }, 1000);
}

function startGame() {
    score = 0; hits = 0; misses = 0; totalReactionTime = 0;
    scoreVal.innerText = 0; hitsVal.innerText = 0; missesVal.innerText = 0;
    spawnTarget();
    gameInterval = setTimeout(endGame, gameLength);
}

function spawnTarget() {
    targetContainer.innerHTML = '';
    const target = document.createElement('div');
    target.className = 'target';
    
    const maxX = targetContainer.clientWidth - 40;
    const maxY = targetContainer.clientHeight - 60;
    target.style.left = `${Math.max(10, Math.floor(Math.random() * maxX))}px`;
    target.style.top = `${Math.max(50, Math.floor(Math.random() * maxY))}px`;
    
    targetSpawnTime = Date.now();
    targetContainer.appendChild(target);

    targetTimeout = setTimeout(() => {
        misses++;
        missesVal.innerText = misses;
        spawnTarget();
    }, 1500);
}

function handleArenaClick(e) {
    if (e.target.classList.contains('target')) {
        clearTimeout(targetTimeout);
        hits++;
        score += 100;
        totalReactionTime += (Date.now() - targetSpawnTime);
        scoreVal.innerText = score;
        hitsVal.innerText = hits;
        spawnTarget();
    } else if (e.target.id === 'target-container') {
        misses++;
        missesVal.innerText = misses;
    }
}

async function endGame() {
    clearTimeout(targetTimeout);
    targetContainer.innerHTML = '';
    arenaScreen.classList.add('hidden');
    setupScreen.classList.remove('hidden');

    const totalShots = hits + misses;
    const accuracy = totalShots > 0 ? ((hits / totalShots) * 100).toFixed(2) : 0;
    const avgReaction = hits > 0 ? Math.floor(totalReactionTime / hits) : 0;

    await submitScoreToBackend(username, score, accuracy, avgReaction);
}

async function submitScoreToBackend(name, finalScore, finalAcc, finalReact) {
    try {
        await fetch(`${CONFIG.API_URL}/scores`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                player_name: name,
                score: parseInt(finalScore),
                accuracy: parseFloat(finalAcc),
                avg_reaction_time: parseInt(finalReact)
            })
        });
    } catch (err) {
        console.error('Failed pushing score telemetry:', err);
    }
}

function escapeHTML(str) {
    return str.replace(/[&<>'"]/g, t => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', "'": '&#39;', '"': '&quot;' }[t] || t));
}
let history = [];
let currentCommandType = 'text';

// Load history from localStorage
try {
    const storedHistory = localStorage.getItem('kowalskiHistory');
    if (storedHistory) {
        history = JSON.parse(storedHistory);
    }
} catch (e) {
    console.error('Failed to load history:', e);
    localStorage.removeItem('kowalskiHistory');
}

document.addEventListener('DOMContentLoaded', () => {
    renderHistory();
    
    // Check if FST commands should be shown
    checkFSTAvailability();
    
    // Attach event listeners
    document.querySelectorAll('button[data-command]').forEach(button => {
        button.addEventListener('click', () => {
            const command = button.dataset.command;
            const type = button.dataset.type;
            const special = button.dataset.special;
            
            executeCommand(command, type, special);
        });
    });
    
    document.getElementById('clearHistory').addEventListener('click', clearHistory);
});

async function checkFSTAvailability() {
    try {
        const response = await fetch('/api/command', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ command: 'fstanagram', input: 'test' })
        });
        const data = await response.json();
        if (data.success || (data.error && !data.error.includes('FST model not loaded'))) {
            document.getElementById('fstCommands').style.display = 'block';
        }
    } catch (error) {
        console.log('FST not available');
    }
}

async function executeCommand(command, type, special) {
    const input = document.getElementById('input').value.trim();
    
    if (type === 'text' && !input) {
        alert('Please enter some text first');
        return;
    }
    
    if (type === 'image') {
        const fileInput = document.getElementById('imageFile');
        if (!fileInput.files[0]) {
            alert('Please select an image file');
            return;
        }
        await executeImageCommand(command, fileInput.files[0]);
    } else {
        let finalInput = input;
        
        // Handle chunk command specially
        if (special === 'chunk') {
            const chunkSizes = document.getElementById('chunkSizes').value.trim();
            if (!chunkSizes) {
                alert('Please enter chunk sizes (e.g., "3" or "2 3 4")');
                return;
            }
            finalInput = chunkSizes + ' ' + input;
        }
        
        await executeTextCommand(command, finalInput);
    }
}

async function executeTextCommand(command, input) {
    const historyItem = {
        command,
        input,
        time: new Date().toISOString(),
        type: 'text'
    };
    
    try {
        addLoadingToHistory(historyItem);
        
        const response = await fetch('/api/command', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ command, input })
        });
        
        const data = await response.json();
        
        if (data.success) {
            historyItem.result = data.result;
        } else {
            historyItem.error = data.error;
        }
    } catch (error) {
        historyItem.error = error.message;
    }
    
    removeLoadingFromHistory();
    addToHistory(historyItem);
}

async function executeImageCommand(command, file) {
    const historyItem = {
        command,
        input: file.name,
        time: new Date().toISOString(),
        type: 'image'
    };
    
    try {
        addLoadingToHistory(historyItem);
        
        const formData = new FormData();
        formData.append('command', command);
        formData.append('image', file);
        
        const response = await fetch('/api/image', {
            method: 'POST',
            body: formData
        });
        
        const data = await response.json();
        
        if (data.success) {
            historyItem.result = data.result;
        } else {
            historyItem.error = data.error;
        }
    } catch (error) {
        historyItem.error = error.message;
    }
    
    removeLoadingFromHistory();
    addToHistory(historyItem);
}

function addToHistory(item) {
    // Create a copy for storage without image data
    const storageItem = JSON.parse(JSON.stringify(item));
    
    // Remove image data from storage
    if (storageItem.result) {
        if (storageItem.command === 'hidden' && storageItem.result.image) {
            storageItem.result.image = '[IMAGE_DATA_REMOVED]';
        }
        if (storageItem.command === 'rgb' && storageItem.result) {
            if (storageItem.result.red) storageItem.result.red = '[IMAGE_DATA_REMOVED]';
            if (storageItem.result.green) storageItem.result.green = '[IMAGE_DATA_REMOVED]';
            if (storageItem.result.blue) storageItem.result.blue = '[IMAGE_DATA_REMOVED]';
        }
    }
    
    history.unshift(item);
    if (history.length > 50) {
        history = history.slice(0, 50);
    }
    
    // Store history without image data
    const storageHistory = history.map(h => {
        const copy = JSON.parse(JSON.stringify(h));
        if (copy.result) {
            if (copy.command === 'hidden' && copy.result.image) {
                copy.result.image = '[IMAGE_DATA_REMOVED]';
            }
            if (copy.command === 'rgb' && copy.result) {
                if (copy.result.red) copy.result.red = '[IMAGE_DATA_REMOVED]';
                if (copy.result.green) copy.result.green = '[IMAGE_DATA_REMOVED]';
                if (copy.result.blue) copy.result.blue = '[IMAGE_DATA_REMOVED]';
            }
        }
        return copy;
    });
    
    try {
        localStorage.setItem('kowalskiHistory', JSON.stringify(storageHistory));
    } catch (e) {
        console.error('Failed to save history:', e);
        // If still failing, clear old history
        if (e.name === 'QuotaExceededError') {
            localStorage.removeItem('kowalskiHistory');
        }
    }
    
    renderHistory();
}

function clearHistory() {
    if (confirm('Are you sure you want to clear the history?')) {
        history = [];
        localStorage.removeItem('kowalskiHistory');
        renderHistory();
    }
}

function addLoadingToHistory(item) {
    const historyDiv = document.getElementById('history');
    const loadingDiv = document.createElement('div');
    loadingDiv.className = 'history-item loading';
    loadingDiv.id = 'loading-item';
    loadingDiv.innerHTML = `
        <div class="history-header">
            <span class="history-command">${item.command}</span>
            <span class="history-time">Processing...</span>
        </div>
        <div class="loading">Executing command...</div>
    `;
    historyDiv.insertBefore(loadingDiv, historyDiv.firstChild);
}

function removeLoadingFromHistory() {
    const loadingItem = document.getElementById('loading-item');
    if (loadingItem) {
        loadingItem.remove();
    }
}

function renderHistory() {
    const historyDiv = document.getElementById('history');
    historyDiv.innerHTML = '';
    
    history.forEach((item, index) => {
        const itemDiv = document.createElement('div');
        itemDiv.className = 'history-item';
        
        const time = new Date(item.time).toLocaleString();
        
        let html = `
            <div class="history-header">
                <span class="history-command">${item.command}</span>
                <span class="history-time">${time}</span>
            </div>
        `;
        
        if (item.type === 'text') {
            html += `<div class="history-input">${escapeHtml(item.input)}</div>`;
        } else {
            html += `<div class="history-input">Image: ${escapeHtml(item.input)}</div>`;
        }
        
        html += '<div class="history-result">';
        
        if (item.error) {
            html += `<div class="error">Error: ${escapeHtml(item.error)}</div>`;
        } else if (item.result) {
            html += renderResult(item.command, item.result);
        }
        
        html += '</div>';
        
        itemDiv.innerHTML = html;
        historyDiv.appendChild(itemDiv);
    });
}

function renderResult(command, result) {
    switch (command) {
        case 'anagram':
        case 'match':
        case 'morse':
        case 'multianagram':
        case 'multimatch':
        case 'offbyone':
        case 't9':
            return renderWordList(result.result);
            
        case 'analysis':
            return renderAnalysis(result.result);
            
        case 'chunk':
            return renderChunks(result.result);
            
        case 'letters':
            return renderLetterDistribution(result.distribution);
            
        case 'shift':
            return renderShifts(result.shifts);
            
        case 'transpose':
            return `<pre>${escapeHtml(result.result)}</pre>`;
            
        case 'wordsearch':
            return renderWordSearch(result);
            
        case 'colours':
        case 'colors':
            return renderColours(result);
            
        case 'hidden':
            return renderImage(result.image, 'Hidden pixels result');
            
        case 'rgb':
            return renderRGBImages(result);
            
        case 'fstanagram':
        case 'fstregex':
        case 'fstmorse':
            return renderFSTMatches(result.matches);
            
        case 'wordlink':
            return renderWordLink(result);
            
        default:
            return `<pre>${JSON.stringify(result, null, 2)}</pre>`;
    }
}

function renderWordList(words) {
    if (!words || words.length === 0) {
        return '<div>No results found</div>';
    }
    
    let html = '<div class="result-list">';
    words.forEach(word => {
        const isSecondary = word.startsWith('_') && word.endsWith('_');
        const displayWord = isSecondary ? word.slice(1, -1) : word;
        html += `<span class="result-item ${isSecondary ? 'secondary' : ''}">${escapeHtml(displayWord)}</span>`;
    });
    html += '</div>';
    return html;
}

function renderAnalysis(results) {
    if (!results || results.length === 0) {
        return '<div>Nothing interesting found</div>';
    }
    
    let html = '<ul>';
    results.forEach(item => {
        html += `<li>${escapeHtml(item)}</li>`;
    });
    html += '</ul>';
    return html;
}

function renderChunks(chunks) {
    return `<div class="result-list">${chunks.map(chunk => 
        `<span class="result-item">${escapeHtml(chunk)}</span>`
    ).join('')}</div>`;
}

function renderLetterDistribution(distribution) {
    let html = '<div>';
    let max = Math.max(...Object.values(distribution));
    
    for (let letter of 'ABCDEFGHIJKLMNOPQRSTUVWXYZ') {
        const count = distribution[letter] || 0;
        const width = max > 0 ? (count / max * 200) : 0;
        html += `
            <div class="letter-bar">
                <span class="letter">${letter}:</span>
                <div class="bar" style="width: ${width}px;"></div>
                <span class="count">${count}</span>
            </div>
        `;
    }
    html += '</div>';
    return html;
}

function renderShifts(shifts) {
    let html = '<div>';
    shifts.forEach(shift => {
        const highlight = shift.score > 0.5 ? 'highlight' : '';
        html += `
            <div class="shift-item ${highlight}">
                <strong>${shift.shift}:</strong> ${escapeHtml(shift.text)} 
                <span style="color: #7f8c8d;">(${shift.score.toFixed(5)})</span>
            </div>
        `;
    });
    html += '</div>';
    return html;
}

function renderWordSearch(result) {
    let html = '<div>';
    html += '<h4>Normal:</h4>';
    html += renderWordList(result.normal);
    html += '<h4>Up/Down:</h4>';
    html += renderWordList(result.updown);
    html += '</div>';
    return html;
}

function renderColours(result) {
    let html = `<div>Total colours: ${result.totalColours}`;
    if (result.truncated) {
        html += ' (showing first 25)';
    }
    html += '</div><div>';
    
    result.colours.forEach(color => {
        html += `
            <div class="color-item">
                <div class="color-swatch" style="background-color: ${color.hex};"></div>
                <div class="color-info">
                    ${color.hex} | RGB(${color.r}, ${color.g}, ${color.b})
                    ${color.a < 255 ? `| A(${color.a})` : ''}
                    | ${color.count} pixels
                </div>
            </div>
        `;
    });
    html += '</div>';
    return html;
}

function renderImage(base64Data, alt) {
    if (base64Data === '[IMAGE_DATA_REMOVED]') {
        return '<div class="image-result"><em>Image data not available in history</em></div>';
    }
    return `<div class="image-result">
        <img src="data:image/png;base64,${base64Data}" alt="${alt}">
    </div>`;
}

function renderRGBImages(result) {
    return `
        <div class="image-grid">
            <div>
                <h4>Red Channel</h4>
                ${renderImage(result.red, 'Red channel')}
            </div>
            <div>
                <h4>Green Channel</h4>
                ${renderImage(result.green, 'Green channel')}
            </div>
            <div>
                <h4>Blue Channel</h4>
                ${renderImage(result.blue, 'Blue channel')}
            </div>
        </div>
    `;
}

function renderFSTMatches(matches) {
    if (!matches || matches.length === 0) {
        return '<div>No results found</div>';
    }
    
    let html = '<div class="result-list">';
    matches.forEach(match => {
        html += `<span class="result-item">${escapeHtml(match.term)} (${match.score})</span>`;
    });
    html += '</div>';
    return html;
}

function renderWordLink(result) {
    let html = `<div>Linking words for '${escapeHtml(result.words[0])}' &lt;&gt; '${escapeHtml(result.words[1])}':</div>`;
    html += '<div class="result-list">';
    result.links.forEach(link => {
        html += `<span class="result-item">${escapeHtml(link.term)} (${link.score})</span>`;
    });
    html += '</div>';
    return html;
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
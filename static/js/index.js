/**
 * @type HTMLLIElement
 */
let generationsInfo;
/**
 * @type HTMLLIElement
 */
let bestInfo;
/**
 * @type HTMLLIElement
 */
let allInfo;
/**
 * @type HTMLElement
 */
let stateMyX;
/**
 * @type HTMLElement
 */
let stateMyY;
/**
 * @type HTMLElement
 */
let stateGoalX;
/**
 * @type HTMLElement
 */
let stateGoalY;

/**
 * @type HTMLCanvasElement
 */
let simulatorCanvas;

/**
 * @type HTMLFormElement
 */
let simulationParams;
/**
 * @type HTMLInputElement
 */
let simulationDelay;

let targetPosition = {
    x: Math.ceil(Math.random() * 500),
    y: Math.ceil(Math.random() * 500),
};

let agentPosition = {
    x: Math.ceil(Math.random() * 500),
    y: Math.ceil(Math.random() * 500),
};

let isRunning = false;
let currentAgent = null;

const bgColor = "#30CBE8";
const targetColor = "#FF4431";
const agentColor = "#3a4a4d";
const targetSize = 6;
const agentSize = 4;

let lastStatus = null;

const patternCanvas = document.createElement('canvas');
const patternContext = patternCanvas.getContext('2d');

// Give the pattern a width and height of 50
patternCanvas.width = 4;
patternCanvas.height = 4;

// Give the pattern a background color and draw an arc
for (let x = 0; x < patternCanvas.width; x += 1) {
    for (let y = 0; y < patternCanvas.height; y += 1) {
        if ((x + y) & 1) {
            patternContext.fillStyle = agentColor;
        } else {
            patternContext.fillStyle = bgColor;
        }
        patternContext.fillRect(x, y, 1, 1);
    }
}
let agentPattern;

function printUl(target, elements) {
    target.innerHTML = `
        <ul>${Object.keys(elements).map(key => `<li><b>${key}: </b> ${elements[key]}</li>`).join('')}</ul>
    `;
}

function updateInfo() {
    generationsInfo.innerHTML = `<b>Current generation: </b>${lastStatus.generation}`;
    const best = lastStatus.best;
    const genes = lastStatus.genes;

    let minFitness = Number.MAX_SAFE_INTEGER;
    let avgFitness = 0;
    for (let i = 0; i < genes.length; i += 1) {
        const fitness = genes[i].fitness;
        minFitness = Math.min(minFitness, fitness);
        avgFitness += fitness;
    }

    avgFitness = avgFitness / genes.length;

    bestInfo.innerHTML = `
        <b>Best agent:</b>
        <ul>
            <li>Score: ${best.fitness.toFixed(3)}</li>
            <li>Survived: ${best.generations}</li>
        </ul>
    `;
    allInfo.innerHTML = `
        <b>Agents:</b>
        <ul>
            <li>Total: ${genes.length}</li>
            <li>Worst: ${minFitness.toFixed(3)}</li>
            <li>Average: ${avgFitness.toFixed(3)}</li>
        </ul>
    `
}

function drawMarker(x, y, size, color) {
    const ctx = simulatorCanvas.getContext('2d')
    const halfSize = size / 2
    ctx.fillStyle = color;
    ctx.fillRect(x - halfSize, y - halfSize, size, size);
}

function moveTarget(x, y) {
    drawMarker(targetPosition.x, targetPosition.y, targetSize, bgColor);
    drawMarker(x, y, targetSize, targetColor);
    targetPosition.x = x;
    targetPosition.y = y;
    actionsQueue = [];

    stateGoalX.innerText = x;
    stateGoalY.innerText = y;

    if (!isRunning && currentAgent) {
        console.trace("RESTARTING !!!")
        isRunning = true;
        makeStep();
    }
}

function moveAgent(x, y) {
    // drawMarker(agentPosition.x, agentPosition.y, agentSize, bgColor);
    drawMarker(x, y, agentSize, agentPattern);
    agentPosition.x = x;
    agentPosition.y = y;
    stateMyX.innerText = x;
    stateMyY.innerText = y;
}

function scale(coordinate) {
    return Math.ceil(coordinate * (simulatorCanvas.width / simulatorCanvas.offsetWidth));
}

let actionsQueue = [];

function runSimulation() {
    let params = {
        agent: currentAgent,
        steps: 500,
        agentX: agentPosition.x,
        agentY: agentPosition.y,
        targetX: targetPosition.x,
        targetY: targetPosition.y,
    };

    return fetch('/catcher', {
        method: 'POST',
        credentials: 'same-origin',
        headers: {
            'Content-Type': 'application/json'
            // 'Content-Type': 'application/x-www-form-urlencoded',
        },
        redirect: 'follow', // manual, *follow, error
        body: JSON.stringify(params)
    }).then(resp => resp.json())
        .then(simulation => {
            if (params.agent !== currentAgent
                || params.agentX !== agentPosition.x
                || params.agentY !== agentPosition.y
                || params.targetX !== targetPosition.x
                || params.targetY !== targetPosition.y
            ) {
                console.error('... old params, should be ignores ...')
                return [];
            }

            actionsQueue = simulation.actions;
            return simulation;
        })
}

function handleAction(action) {
    switch (action) {
        case 'Up':
            return moveAgent(agentPosition.x, agentPosition.y + 1);
        case 'Down':
            return moveAgent(agentPosition.x, agentPosition.y - 1);
        case 'Left':
            return moveAgent(agentPosition.x - 1, agentPosition.y);
        case 'Right':
            return moveAgent(agentPosition.x + 1, agentPosition.y);
    }
}

function makeStep() {
    if (!actionsQueue.length) {
        return runSimulation().then(makeStep);
    }

    handleAction(actionsQueue.shift());

    if (agentPosition.x === targetPosition.x && agentPosition.y === targetPosition.y) {
        isRunning = false;
        console.log("DONE!")
        return;
    }

    const strategy = document.getElementById("agent-func").value;
    const evaluate = (myX, myY) => {
        const code = `
        (function(MyX, MyY, GoalX, GoalY) {
            return ${strategy};
        })(${myX}, ${myY}, ${targetPosition.x}, ${targetPosition.y})
        `
        return eval(code).toFixed(3);
    }

    const myX = agentPosition.x;
    const myY = agentPosition.y;

    printUl(document.getElementById("test-agent-result"), {
        UP: evaluate(myX, myY + 1),
        DOWN: evaluate(myX, myY - 1),
        LEFT: evaluate(myX - 1, myY),
        RIGHT: evaluate(myX + 1, myY)
    });

    setTimeout(makeStep, +simulationDelay.value)
}

function run() {
    isRunning = true;
    agentPosition = {
        x: Math.ceil(Math.random() * 500),
        y: Math.ceil(Math.random() * 500),
    };
    loadAgent();
    const ctx = simulatorCanvas.getContext('2d')
    ctx.fillStyle = bgColor;
    ctx.fillRect(0, 0, simulatorCanvas.width, simulatorCanvas.height);
    moveTarget(targetPosition.x, targetPosition.y);
    moveAgent(agentPosition.x, agentPosition.y);
    runSimulation().then(result => {
        document.getElementById("agent-func").innerText = result.asString;
        makeStep();
    })
}

function loadAgent() {
    const mode = simulationParams.elements['mode'].value;
    let index = simulationParams.elements['index'].value;
    if (mode === 'best') {
        currentAgent = lastStatus.best.agent;
    } else if (mode === 'index') {
        if (!lastStatus.genes[index]) {
            simulationParams.elements['index'].value = 0;
            index = 0;
        }

        currentAgent = lastStatus.genes[index].agent;
    }
}

function refresh() {
    return fetch("/status").then(response => response.json())
        .then(result => {
            lastStatus = result;
            updateInfo();
            setTimeout(refresh, 3000);
        })
}

function initApp() {
    generationsInfo = document.getElementById('info-generations');
    bestInfo = document.getElementById('info-best');
    allInfo = document.getElementById('info-all');
    simulationDelay = document.getElementById('simulation-delay');

    stateMyX = document.getElementById('state-myx');
    stateMyY = document.getElementById('state-myy');
    stateGoalX = document.getElementById('state-goalx');
    stateGoalY = document.getElementById('state-goaly');

    simulatorCanvas = document.getElementById('renderer');
    agentPattern = simulatorCanvas.getContext('2d').createPattern(patternCanvas, 'repeat')
    simulationParams = document.getElementById('simulation-params');
    simulationParams.addEventListener('change', () => {
        simulationParams.elements['index'].disabled = simulationParams.elements['mode'].value === 'best'
        loadAgent();
    });
    simulationParams.elements['index'].disabled = simulationParams.elements['mode'].value === 'best'

    simulatorCanvas.addEventListener('click', (e) => {
        var rect = e.target.getBoundingClientRect();
        var x = scale(e.clientX - rect.left); //x position within the element.
        var y = scale(e.clientY - rect.top);  //y position within the element.
        moveTarget(x, y)
    });

    document.getElementById("run").addEventListener("click", run);

    refresh().then(() => {
        document.getElementById('loader').remove();
        const content = document.getElementById('app');
        content.classList.remove('hidden')
    });
}


document.addEventListener("DOMContentLoaded", initApp)
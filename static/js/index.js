/**
 * @type HTMLLIElement
 */
let generationsInfo;
/**
 * @type HTMLLIElement
 */
let statusInfo;
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

/**
 * @type HTMLInputElement
 */
let sizeParam;
/**
 * @type HTMLInputElement
 */
let elitesSizeParam;
/**
 * @type HTMLInputElement
 */
let childrenSizeParam;
/**
 * @type HTMLInputElement
 */
let mutationChanceParam;

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
const startColor = "#00ff00"
const targetSize = 6;
const agentSize = 2;

let lastStatus = null;
let editingParams = false;

function printUl(target, elements) {
    target.innerHTML = `
        <ul>${Object.keys(elements).map(key => `<li><b>${key}: </b> ${elements[key]}</li>`).join('')}</ul>
    `;
}

function updateInfo() {
    statusInfo.innerHTML = `<b>Status: </b>${lastStatus.finished ? 'Finished!' : 'Running...'}`;
    generationsInfo.innerHTML = `<b>Current generation: </b>${lastStatus.generation}`;
    const best = lastStatus.best;
    const genes = lastStatus.genes;

    if (!editingParams) {
        sizeParam.value = lastStatus.settings.size;
        elitesSizeParam.value = lastStatus.settings.elitesSize;
        childrenSizeParam.value = lastStatus.settings.childrenSize;
        mutationChanceParam.value = lastStatus.settings.mutationChance;
    }

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
        isRunning = true;
        resetSimulationScreen();
        makeStep();
    }
}

function moveAgent(x, y) {
    // drawMarker(agentPosition.x, agentPosition.y, agentSize, bgColor);
    drawMarker(x, y, agentSize, agentColor);
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
    if (agentPosition.x === targetPosition.x && agentPosition.y === targetPosition.y) {
        isRunning = false;
        console.log("DONE!")
        return;
    }

    if (!actionsQueue.length) {
        return runSimulation().then(makeStep);
    }

    handleAction(actionsQueue.shift());

    const strategy = document.getElementById("agent-func-optimised").value;
    const evaluateF = (myX, myY) => {
        const code = `
        (function(MyX, MyY, GoalX, GoalY) {
            return ${strategy};
        })(${myX}, ${myY}, ${targetPosition.x}, ${targetPosition.y})
        `
        return eval(code).toFixed(3);
    }
    const evaluateA = (myX, myY) => {
        const code = `
        (function(MyX, MyY, GoalX, GoalY) {
            const Up = 'Down';
            const Down = 'Up';
            const Left = 'Left';
            const Right = 'Right';
            return ${strategy};
        })(${myX}, ${myY}, ${targetPosition.x}, ${targetPosition.y})
        `
        return eval(code);
    }

    const myX = agentPosition.x;
    const myY = agentPosition.y;

    const resultEl = document.getElementById("test-agent-result");
    if (currentAgent.type === "float") {
        printUl(resultEl, {
            UP: evaluateF(myX, myY + 1),
            DOWN: evaluateF(myX, myY - 1),
            LEFT: evaluateF(myX - 1, myY),
            RIGHT: evaluateF(myX + 1, myY)
        });
    } else if (currentAgent.type === "action") {
        resultEl.innerText = evaluateA(myX, myY);
    }

    setTimeout(makeStep, +simulationDelay.value)
}

function resetSimulationScreen() {
    const ctx = simulatorCanvas.getContext('2d')
    ctx.fillStyle = bgColor;
    ctx.fillRect(0, 0, simulatorCanvas.width, simulatorCanvas.height);
    moveTarget(targetPosition.x, targetPosition.y);
    moveAgent(agentPosition.x, agentPosition.y);
    drawMarker(agentPosition.x, agentPosition.y, targetSize + 6, startColor);
}

function run() {
    isRunning = true;
    loadAgent();
    resetSimulationScreen();
    runSimulation().then(result => {
        document.getElementById("agent-func-optimised").innerText = result.asString;
        document.getElementById("agent-func-original").innerText = result.asStringOriginal;
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
    statusInfo = document.getElementById('info-status');
    bestInfo = document.getElementById('info-best');
    allInfo = document.getElementById('info-all');
    simulationDelay = document.getElementById('simulation-delay');

    sizeParam = document.getElementById("sizeParam")
    elitesSizeParam = document.getElementById("elitesSizeParam")
    childrenSizeParam = document.getElementById("childrenSizeParam")
    mutationChanceParam = document.getElementById("mutationChanceParam")

    stateMyX = document.getElementById('state-myx');
    stateMyY = document.getElementById('state-myy');
    stateGoalX = document.getElementById('state-goalx');
    stateGoalY = document.getElementById('state-goaly');

    simulatorCanvas = document.getElementById('renderer');
    simulationParams = document.getElementById('simulation-params');
    simulationParams.addEventListener('change', () => {
        simulationParams.elements['index'].disabled = simulationParams.elements['mode'].value === 'best'
        loadAgent();
    });
    simulationParams.elements['index'].disabled = simulationParams.elements['mode'].value === 'best'

    simulatorCanvas.addEventListener('click', (e) => {
        const rect = e.target.getBoundingClientRect();
        const x = scale(e.clientX - rect.left); //x position within the element.
        const y = scale(e.clientY - rect.top);  //y position within the element.
        if (e.altKey) {
            moveAgent(x, y)
        } else {
            moveTarget(x, y)
        }
        resetSimulationScreen();
    });

    document.getElementById("run").addEventListener("click", run);

    const changeParams = document.getElementById("change-params");
    const applyParams = document.getElementById("apply-params");
    const cancelParams = document.getElementById("cancel-params");

    changeParams.addEventListener("click", () => {
        editingParams = true;
        enable(sizeParam, elitesSizeParam, childrenSizeParam, mutationChanceParam);
        hide(changeParams);
        show(applyParams, cancelParams);
    });

    cancelParams.addEventListener('click', () => {
        editingParams = false;
        disable(sizeParam, elitesSizeParam, childrenSizeParam, mutationChanceParam);
        show(changeParams);
        hide(applyParams, cancelParams);
        updateInfo();
    });

    applyParams.addEventListener('click', () => {
        // first, some basic validation
        const size = +sizeParam.value;
        const mutationChance = +mutationChanceParam.value;
        const childrenSize = +childrenSizeParam.value;
        const elitesSize = +elitesSizeParam.value;

        if (size < 5 || size > 100) {
            return alert("Size should be in range [5, 100]");
        }
        if (childrenSize < 0 || childrenSize > 100) {
            return alert("Children should be in range [0, 100]");
        }
        if (elitesSize < 0 || elitesSize > 10) {
            return alert("Elites should be in range [0, 10]");
        }
        if (mutationChance < 0 || mutationChance > 1) {
            return alert("Mutation chance should be in range [0, 1]");
        }
        disable(sizeParam, elitesSizeParam, childrenSizeParam, mutationChanceParam);

        fetch('/evolve', {
            method: 'POST',
            credentials: 'same-origin',
            headers: {
                'Content-Type': 'application/json'
                // 'Content-Type': 'application/x-www-form-urlencoded',
            },
            redirect: 'follow', // manual, *follow, error
            body: JSON.stringify({
                size, mutationChance, childrenSize, elitesSize
            })
        }).then(resp => resp.json())
            .then(result => {
                lastStatus = result
                show(changeParams);
                hide(applyParams, cancelParams);
                updateInfo();
            })
            .catch(error => {
                alert("Request failed")
                console.error("Failed to restart server: ", error);
                enable(sizeParam, elitesSizeParam, childrenSizeParam, mutationChanceParam);
            })

    });

    refresh().then(() => {
        document.getElementById('loader').remove();
        show(document.getElementById('app'));
    });
}

function hide(el) {
    for (let i = 0; i < arguments.length; i++) {
        arguments[i].classList.add('hidden')
    }
}

function show() {
    for (let i = 0; i < arguments.length; i++) {
        arguments[i].classList.remove('hidden')
    }
}


function enable() {
    for (let i = 0; i < arguments.length; i++) {
        arguments[i].removeAttribute('disabled');
    }
}
function disable() {
    for (let i = 0; i < arguments.length; i++) {
        arguments[i].setAttribute('disabled', 'disabled');
    }
}


document.addEventListener("DOMContentLoaded", initApp)
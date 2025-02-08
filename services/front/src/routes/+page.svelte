<script>
// @ts-nocheck

  import { onMount } from "svelte";
  let enemy_color_bindings = {
    1: "#121212",
    2: "#960469",
    3: "#A7F00F"
  }
  let tower_color_bindings = {
    1: "#121212",
    2: "#960469",
    3: "#A7F00F"
  }
  
  let gameState =  {
      map: {
        name: "Desert Plains",
        path: [
          { x: 0, y: 50 },
          { x: 500, y: 50 },
          { x: 500, y: 350 },
          { x: 150, y: 350 },
          { x: 150, y: 550 },
          { x: 700, y: 550 },
          { x: 700, y: 50 }
        ],
        dimensions: { width: 800, height: 600 }
      },
      players: { attackers: [], defenders: [] },
      enemies: [],
      available_towers: [
        { id: 1, type: "Archer Tower", attack: 15, range: 1, cost: 100},
        { id: 2, type: "Archerer Tower", attack: 30, range: 2, cost: 200},
        { id: 3, type: "Archerest Tower", attack: 45, range: 3, cost: 400}
      ]
      ,available_enemies:
      [{"id":1,"type":"Red blun","position_index":0,"health":30,"speed":10,"cost":10},
      {"id":2,"type":"Green blun","position_index":0,"health":80,"speed":10,"cost":40},
      {"id":3,"type":"Blue blun","position_index":0,"health":200,"speed":10,"cost":60}]
  };
  $: available_enemies = gameState.available_enemies
  console.log("available_enemies:", available_enemies)
  let canvas;
  let selectedTower = null;
  let selectedEnemy = null;
  const { dimensions } = gameState.map;
  const gridSize = 40; // Size of each grid cell
  let placedTowers = []; // Track placed towers
  let enemies = [];
  let path = [
    { X:0, Y:50 },
    { X:50, Y:50 },
    { X:100, Y:50 },
    { X:150, Y:50 },
    { X:200, Y:50 },
    { X:250, Y:50 },
    { X:300, Y:50 },
    { X:350, Y:50 },
    { X:400, Y:50 },
    { X:450, Y:50 },
    { X:500, Y:50 },
    { X:500, Y:100 },
    { X:500, Y:150 },
    { X:500, Y:200 },
    { X:500, Y:250 },
    { X:500, Y:300 },
    { X:500, Y:350 },
    { X:450, Y:350 },
    { X:400, Y:350 },
    { X:350, Y:350 },
    { X:300, Y:350 },
    { X:250, Y:350 },
    { X:200, Y:350 },
    { X:150, Y:350 },
    { X:150, Y:400 },
    { X:150, Y:450 },
    { X:150, Y:500 },
    { X:150, Y:550 },
    { X:200, Y:550 },
    { X:250, Y:550 },
    { X:300, Y:550 },
    { X:350, Y:550 },
    { X:400, Y:550 },
    { X:450, Y:550 },
    { X:500, Y:550 },
    { X:550, Y:550 },
    { X:600, Y:550 },
    { X:650, Y:550 },
    { X:700, Y:550 },
    { X:700, Y:500 },
    { X:700, Y:450 },
    { X:700, Y:400 },
    { X:700, Y:350 },
    { X:700, Y:300 },
    { X:700, Y:250 },
    { X:700, Y:200 },
    { X:700, Y:150 },
    { X:700, Y:100 },
    { X:700, Y:50 },
        ]
  function hexToRgb(hex) {
  const bigint = parseInt(hex.replace('#', ''), 16);
  const r = (bigint >> 16) & 255;
  const g = (bigint >> 8) & 255;
  const b = bigint & 255;
  return `${r},${g},${b}`;
}
function toGrid(x, y) {
      return { x: Math.floor(x / gridSize), y: Math.floor(y / gridSize) };
    }



  function drawMap(ctx) {
    ctx.clearRect(0, 0, gameState.map.width, gameState.map.height);

    // Draw grid
    ctx.strokeStyle = "#ccc";
    for (let x = 0; x <= dimensions.width; x += gridSize) {
      ctx.beginPath();
      ctx.moveTo(x, 0);
      ctx.lineTo(x, dimensions.height);
      ctx.stroke();
    }
    for (let y = 0; y <= dimensions.height; y += gridSize) {
      ctx.beginPath();
      ctx.moveTo(0, y);
      ctx.lineTo(dimensions.width, y);
      ctx.stroke();
    }
    console.log('creating enemies', gameState.enemies)
    // Draw enemies (balloons)
    gameState.enemies.forEach(enemy => {
    console.log('got it')
    // Get the position of the enemy from its positionIndex
    const positionIndex = enemy.position_index; // Assuming this is available
    console.log("NEWPATH: ",path, "enemy", enemy)

    const enemyPathPoint = path[positionIndex];  // Get the path point corresponding to the enemy
    
    // Convert the path position to grid coordinates
    const gridPoint = toGrid(enemyPathPoint.X, enemyPathPoint.Y);
    const posX = gridPoint.x * gridSize + gridSize / 2;
    const posY = gridPoint.y * gridSize + gridSize / 2;

    // Draw the enemy as an ellipse (balloon)
    ctx.beginPath();
    ctx.ellipse(posX, posY, 15, 20, 0, 0, Math.PI * 2);  // Ellipse shape for balloon
    ctx.fillStyle = enemy_color_bindings[enemy.id];  // Use a function to get the color based on the balloon type
    ctx.fill();
    ctx.strokeStyle = '#000';  // Outline color for the balloon
    ctx.stroke();
  });
    
    
    ctx.beginPath();
path.forEach((point, index) => {
  const gridPoint = toGrid(point.X, point.Y);  // Adjust for the new path format (X, Y)

  // Adjust for grid positioning
  const posX = gridPoint.x * gridSize + gridSize / 2;
  const posY = gridPoint.y * gridSize + gridSize / 2;



  if (index === 0) {
    ctx.moveTo(posX, posY);  // Start path at the first point
  } else {
    ctx.lineTo(posX, posY);  // Draw lines to each subsequent point
  }
});

    ctx.strokeStyle = '#000';  // Path color (black)
    ctx.lineWidth = 3;         // Line width for the path
    ctx.stroke();


    placedTowers.forEach((tower) => {
      let tower_color = tower_color_bindings[tower.id]
      ctx.fillStyle = `rgba(${hexToRgb(tower_color)}, 0.1)`;
      const rangeSize = (tower.range * 2 + 1) * gridSize;
      const rangeX = (tower.x - tower.range) * gridSize;
      const rangeY = (tower.y - tower.range) * gridSize;
      ctx.fillRect(rangeX, rangeY, rangeSize, rangeSize);


      ctx.beginPath();
      ctx.arc(
        tower.x * gridSize + gridSize / 2,
        tower.y * gridSize + gridSize / 2,
        15,
        0,
        Math.PI * 2
      );
      ctx.fillStyle = tower_color;
      ctx.fill();
      ctx.strokeStyle = "black";
      ctx.stroke();
    });
  }

  function getAttackRange(x, y, range) {
  let attacking = [];
    //TODO

  return attacking;  // Return the list of indexes in range
}

let socket

  onMount(async() => {

    await initializeWebSocketConn();


    canvas.addEventListener("click", (event) => {
      if (!selectedTower) {
        
        return;
      }
      
      const rect = canvas.getBoundingClientRect();
      const x = Math.floor((event.clientX - rect.left) / gridSize);
      const y = Math.floor((event.clientY - rect.top) / gridSize);

      if (isOccupied(x, y)) {
        return;
      }
      console.log("placed tower")
      selectedTower.attacking = getAttackRange(selectedTower.x, selectedTower.y, selectedTower.range)
      //send to server TODO
      placeTower(x, y);
      selectedTower = null;
    
    });
  });



  async function initializeWebSocketConn() {
    const ctx = canvas.getContext("2d");
    canvas.width = dimensions.width;
    canvas.height = dimensions.height;
    socket = await new WebSocket(`ws://localhost:8088/game`);

    socket.onmessage = (event) => {
          const newGameState = JSON.parse(event.data);
          console.log("GameState",newGameState)
          gameState = newGameState
          drawMap(ctx);
          console.log("new Game state", newGameState.enemies)
          available_enemies = newGameState.available_enemies
          console.log("new available_enemies:", available_enemies)
        };

    socket.onclose = async (event) => {
      console.log("WebSocket connection closed:", event.reason);
    };

    socket.onerror = async (error) => {
      console.error("WebSocket error:", error);
    };

    socket.onopen = async () => {
      console.log("WebSocket connection established");

    }

  }

  function isOccupied(x, y) {
    console.log(path)
    //this bull$hit doesn't work, TODO
    let ex = x*50
    let ey = y*50
    console.log(ex,ey)
    path.forEach((point, index) => {
      if (Number(point.X) === ex && point.Y === Number(ey)) {
        console.log(ex,point.X, ey,point.Y, ex.type, point.X.type)
        return true
      }
    })
    return placedTowers.some((tower) => tower.x === x && tower.y === y);
  }

  // function placeTower(x, y) {
  //   placedTowers.push({ ...selectedTower, x, y });
  //   const ctx = canvas.getContext("2d");
  //   drawMap(ctx);
  // }

  function selectTower(tower) {
    selectedTower = tower;
  }
  async function selectEnemy(enemy) {
   await console.log('sending data')
   console.log("SELECTED ENEMY: ", enemy)
   spawnEnemy(enemy);
   } // Ensure this is set properly (could be redundant here)
   
   async function spawnEnemy(selectedEnemy) {
// Check if the WebSocket is open (readyState 1 means open)
console.log(selectedEnemy)
if (socket && socket.readyState === WebSocket.OPEN) {
  const message = {
    operation: "addEnemy",
    enemy: selectedEnemy
  };
  console.log('Sending message:', message);

  try {
    // Send the message to the server
    socket.send(JSON.stringify(message));
    console.log("Message sent:", message);
  } catch (error) {
    console.error("Error sending message through WebSocket:", error);
  }
}else{
  console.log("psycha")
  await initializeWebSocketConn()
}
return 1
  }
let test1
</script>

<main>
  <div style="display: flex; gap: 20px;">
    <div>
      <h1>{gameState.map.name}</h1>
      <canvas bind:this={canvas}></canvas>
    </div>

    <div style="width: 300px; background-color: #f4f4f4; padding: 20px; border: 1px solid #ccc;">
      <h2>Shop</h2>

      <h3>Available Towers</h3>
      <h4> Your gold ???</h4>
      <ul>
        {#each gameState.available_towers as tower}
          <li>
            <strong>{tower.type}</strong>
            <div>Attack: {tower.attack}</div>
            <div>Range: {tower.range}</div>
            <div>Cost: {tower.cost} gold</div>
            <button on:click={() => selectTower(tower)} style="background-color: {tower_color_bindings[tower.id]}; color: white;">
              Select
            </button>
          </li>
        {/each}
      </ul>
    </div>
    <div style="width: 300px; background-color: #f4f4f4; padding: 20px; border: 1px solid #ccc;">
      <h2>Shop</h2>

      <h3>Available Ballons</h3>
      <h4> Your gold ???</h4>
      {#key available_enemies}
      <ul>
        {#each available_enemies as enemy}
          <li>
            <strong>{enemy.type}</strong>
            <div>Health: {enemy.health}</div>
            <div>Speed: {enemy.speed}</div>
            <div>Cost: {enemy.cost} gold</div>
            <button on:click={() => selectEnemy(enemy)} style="background-color: {enemy_color_bindings[enemy.id]}; color: white;">
              Select
            </button>
          <div>
            <button on:click={spawnEnemy}>Send Ballon</button>
          </div>
          </li>
        {/each}
      </ul>
      {/key}
    </div>
  </div>
</main>

<style>
  canvas {
    border: 1px solid #000;
  }

  h1, h2, h3 {
    margin: 0 0 10px;
  }

  ul {
    list-style-type: none;
    padding: 0;
  }

  li {
    margin-bottom: 10px;
    padding: 10px;
    border: 1px solid #ddd;
    background-color: #fff;
  }

  button {
    margin-top: 5px;
    border: none;
    padding: 5px 10px;
    cursor: pointer;
  }

  button:hover {
    opacity: 0.8;
  }
</style>

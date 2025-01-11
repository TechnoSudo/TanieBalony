<script>
    const gameState = {
      "gameState": {
        "map": {
          "name": "Desert Plains",
          "path": [
            { "x": 0, "y": 50 },
            { "x": 500, "y": 50 },
            { "x": 500, "y": 350 },
            { "x": 150, "y": 350 },
            { "x": 150, "y": 550 },
            { "x": 700, "y": 550 },
            { "x": 700, "y": 50 },
          ],
          "dimensions": { "width": 800, "height": 600 }
        },
        "players": {
          "attackers": [
            {
              "id": 1,
              "name": "Player1",
              "gold": 500,
              "spawnedEnemies": [
                { "type": "RedBalloon", "count": 5 },
                { "type": "BlueBalloon", "count": 3 }
              ]
            },
            {
              "id": 2,
              "name": "Player2",
              "gold": 450,
              "spawnedEnemies": [
                { "type": "GreenBalloon", "count": 2 },
                { "type": "YellowBalloon", "count": 4 }
              ]
            }
          ],
          "defenders": [
            {
              "id": 3,
              "name": "Player3",
              "gold": 600,
              "placedTowers": [
                { "id": 1, "type": "Cannon", "position": { "x": 150, "y": 250 }, "level": 2 },
                { "id": 2, "type": "Sniper", "position": { "x": 300, "y": 200 }, "level": 1 }
              ]
            },
            {
              "id": 4,
              "name": "Player4",
              "gold": 550,
              "placedTowers": [
                { "id": 3, "type": "SpikeFactory", "position": { "x": 400, "y": 300 }, "level": 3 },
                { "id": 4, "type": "DartMonkey", "position": { "x": 500, "y": 400 }, "level": 1 }
              ]
            }
          ]
        },
        "enemies": [
          { "id": 1, "type": "RedBalloon", "position": { "x": 50, "y": 200 }, "health": 1, "speed": 2 },
          { "id": 2, "type": "BlueBalloon", "position": { "x": 120, "y": 200 }, "health": 2, "speed": 3 },
          { "id": 3, "type": "GreenBalloon", "position": { "x": 300, "y": 250 }, "health": 3, "speed": 4 }
        ]
      }
    };
  
    const { path, dimensions } = gameState.gameState.map;
    const towers = gameState.gameState.players.defenders.flatMap(player => player.placedTowers);
    const enemies = gameState.gameState.enemies;
  
    // Define grid size (cell width/height)
    const gridSize = 40; // You can change this value to adjust the grid cell size
    const cols = Math.floor(dimensions.width / gridSize);
    const rows = Math.floor(dimensions.height / gridSize);
  
    // Function to map coordinates to grid positions
    function toGrid(x, y) {
      return { x: Math.floor(x / gridSize), y: Math.floor(y / gridSize) };
    }
  
    // Function to draw the map
    function drawMap(ctx) {
      // Draw grid
      ctx.strokeStyle = '#ccc';
      ctx.lineWidth = 1;
  
      for (let i = 0; i <= cols; i++) {
        ctx.beginPath();
        ctx.moveTo(i * gridSize, 0);
        ctx.lineTo(i * gridSize, dimensions.height);
        ctx.stroke();
      }
  
      for (let i = 0; i <= rows; i++) {
        ctx.beginPath();
        ctx.moveTo(0, i * gridSize);
        ctx.lineTo(dimensions.width, i * gridSize);
        ctx.stroke();
      }
  
      // Draw path (scaled to grid)
      ctx.beginPath();
      path.forEach((point, index) => {
        const gridPoint = toGrid(point.x, point.y);
        if (index === 0) {
          ctx.moveTo(gridPoint.x * gridSize + gridSize / 2, gridPoint.y * gridSize + gridSize / 2);
        } else {
          ctx.lineTo(gridPoint.x * gridSize + gridSize / 2, gridPoint.y * gridSize + gridSize / 2);
        }
      });
      ctx.strokeStyle = '#000';
      ctx.lineWidth = 3;
      ctx.stroke();
  
      // Draw towers (scaled to grid)
      towers.forEach(tower => {
        const towerPosition = toGrid(tower.position.x, tower.position.y);
        ctx.beginPath();
        ctx.arc(towerPosition.x * gridSize + gridSize / 2, towerPosition.y * gridSize + gridSize / 2, 15, 0, Math.PI * 2);
        ctx.fillStyle = 'gray'; // Placeholder color for towers
        ctx.fill();
        ctx.strokeStyle = 'black';
        ctx.stroke();
      });
  
      // Draw enemies (scaled to grid)
      enemies.forEach(enemy => {
        const enemyPosition = toGrid(enemy.position.x, enemy.position.y);
        ctx.beginPath();
        ctx.arc(enemyPosition.x * gridSize + gridSize / 2, enemyPosition.y * gridSize + gridSize / 2, 10, 0, Math.PI * 2);
        ctx.fillStyle = getBalloonColor(enemy.type);
        ctx.fill();
        ctx.strokeStyle = 'black';
        ctx.stroke();
      });
    }
  
    // Function to get balloon color based on type
    function getBalloonColor(type) {
      switch (type) {
        case 'RedBalloon': return 'red';
        case 'BlueBalloon': return 'blue';
        case 'GreenBalloon': return 'green';
        case 'YellowBalloon': return 'yellow';
        default: return 'gray';
      }
    }
  
    // When the component mounts, initialize canvas
    import { onMount } from 'svelte';
    let canvas;
    onMount(() => {
      const ctx = canvas.getContext('2d');
      canvas.width = dimensions.width;
      canvas.height = dimensions.height;
      drawMap(ctx);
    });
  </script>
  
  <main>
    <h1>{gameState.gameState.map.name}</h1>
    <canvas bind:this={canvas}></canvas>
  </main>
  
  <style>
    canvas {
      border: 1px solid #000;
    }
  
    h1 {
      font-size: 24px;
      color: #333;
    }
  </style>
  
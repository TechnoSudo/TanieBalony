<script>
    /**
   * @type {WebSocket}
   */
    let socket;
    let messageToSend = '';
    let receivedMessage = '';
  
    // Connect to WebSocket when component mounts
    import { onMount } from 'svelte';
  
    onMount(() => {
      // Replace 'ws://your-websocket-url' with your actual WebSocket URL
      socket = new WebSocket('ws://127.0.0.1:8088/game');
  
      // Listen for incoming messages
      socket.onmessage = (event) => {
        receivedMessage = event.data;  // Handle the incoming message
      };
  
      socket.onopen = () => {
        console.log('WebSocket connection established');
      };
  
      socket.onerror = (error) => {
        console.error('WebSocket error: ', error);
      };
  
      socket.onclose = () => {
        console.log('WebSocket connection closed');
      };
    });
  
    // Send a message when the user clicks a button
    function sendMessage() {
      if (socket.readyState === WebSocket.OPEN) {
        socket.send(messageToSend); // Send the message typed by the user
      }
    }
  
    // Optional: Handle the text input for the message
    /**
   * @param {{ target: { value: string; }; }} event
   */
    function updateMessage(event) {
      messageToSend = event.target.value;
    }
  </script>
  
  <main>
    <h1>Svelte WebSocket Example</h1>
  
    <div>
      <input
        type="text"
        placeholder="Type a message"
        bind:value={messageToSend}
        on:input={updateMessage}
      />
      <button on:click={sendMessage}>Send Message</button>
    </div>
  
    <div>
      <h3>Received Message:</h3>
      <p>{receivedMessage}</p>
    </div>
  </main>
  
  <style>
    main {
      font-family: Arial, sans-serif;
      padding: 20px;
    }
    
    input, button {
      padding: 10px;
      margin: 10px 0;
    }
    
    button {
      cursor: pointer;
    }
  </style>
  
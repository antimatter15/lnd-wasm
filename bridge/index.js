const WebSocket = require('ws');
const net = require('net');


const wss = new WebSocket.Server({
  port: 7281,
});


wss.on('connection', async (ws) => {

  let sockets = []
  ws.on('message', (message) => {
    console.log(message)
    let payload = JSON.parse(message);

    if(payload.action === 'dial'){
      let client = new net.Socket();
      let id = sockets.length;
      sockets.push(client)

      client.connect(payload.port, payload.host, function() {
        console.log('Connected');

        ws.send(JSON.stringify({
          action: 'dial_finish',
          id: id
        }))
        
      });

      client.on('close', () => {
        console.log('socket closed')
      })

      client.on('data', e => {
        // console.log(e)
        ws.send(JSON.stringify({
          action: 'read',
          id: id,
          data: Array.from(Uint8Array.from(e))
        }))
      })

      client.on('error', err => {
        console.log(err)
        ws.send(JSON.stringify({
          action: 'dial_finish',
          error: err.toString()
        }))
      })
    }else if(payload.action === 'write'){
      let client  = sockets[payload.id]
      client.write(Uint8Array.from(payload.data), function(){
        ws.send(JSON.stringify({
          action: 'write_finish'
        }))
      })

    }
  });
});
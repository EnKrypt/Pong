## Multiplayer pong using Golang and React

### How to run

1. Open a terminal and navigate to the backend folder
2. Compile the backend with `go build .`
3. Execute the binary with `./pong` or `./pong.exe`
4. Navigate to the frontend folder
5. Run `npm install` or `yarn`
6. Run `yarn start` which should open your browser to http://localhost:3000

### Rules of Pong

1. SPACE starts a new game
2. Player 1 on the left controls with W/S while Player 2 on the right controls with UP/DOWN
3. Players can play from different browsers/computers over the internet
4. Deflect the ball with your paddle. If the ball passes behind your paddle line, the opponent gets a point
5. First player to reach 11 points wins

### Caveats

1. This will not work on mobile devices
2. When hosted, the backend could be deployed at a geographical region that is not close to you, resulting in minor input lag. This can be potentially solved with geo load balancing or GSLB
3. An unstable network such as a weak WiFi connection will prevent the game from being smooth
4. There is currently no logic to prevent a different player on a different browser to interfere with the controls of your player
5. If you choose to contribute or modify the code, there are currently no tests to ensure that it conforms to spec

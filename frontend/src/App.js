import { useEffect, useState } from 'react';
import registerKeyboardEvents from './keyboardEvents';

/* The Markup is very barebones, so we don't need to separate
 * it out into components, but I will make use of React hooks
 * to play with client side state and conveniently leverage the
 * re-rendering strategies that come with it.
 */

let socket = {};

const App = () => {
  const [loading, setLoading] = useState(true);
  const [playing, setPlaying] = useState(false);
  const [player1Score, setPlayer1Score] = useState(0);
  const [player2Score, setPlayer2Score] = useState(0);
  const [player1Position, setPlayer1Position] = useState(300);
  const [player2Position, setPlayer2Position] = useState(300);
  const [ballPosition, setBallPosition] = useState([0, 0]);
  /* X and Y axis of the ball position are tightly coupled and
   * exist in the same state variable unlike other state values
   * like player 1 position and player 2 position which are
   * independent of each other. This helps us in the future if
   * we want to compare with previous values to update and
   * re-render only some parts as an optimization.
   */

  registerKeyboardEvents(socket, playing);

  useEffect(() => {
    socket = new WebSocket(`ws://${window.location.hostname}:8080`);

    socket.addEventListener('message', (event) => {
      setLoading(false);
      const params = event.data.split(' ');
      if (params[0] === 'READY') {
        // If a game has just ended, we want to update the winner's score
        if (params.length > 1) {
          if (parseInt(params[1]) === 1) {
            setPlayer1Score(11);
          } else if (parseInt(params[1]) === 2) {
            setPlayer2Score(11);
          }
        }
        setPlaying(false);
      } else {
        /* There is a game going on, the backend is sending us the
         * game state per tick, so we update our app's state accordingly
         */
        setPlaying(true);
        setBallPosition([params[0], params[1]]);
        setPlayer1Position(params[2]);
        setPlayer2Position(params[3]);
        setPlayer1Score(params[4]);
        setPlayer2Score(params[5]);
      }
    });
  }, []);

  return (
    <div className="app">
      <div className="controls">
        <div>Player 1 moves with W/S</div>
        <div>Player 2 moves with UP/DOWN</div>
      </div>
      <div className="pong">
        {loading ? (
          <div className="loading">Loading</div>
        ) : (
          <>
            <div className="divider divider-right"></div>
            <div className="divider divider-left"></div>
            <div className="score">
              <div className={player1Score === 11 ? 'green' : ''}>
                {player1Score}
              </div>
              <div className={player2Score === 11 ? 'green' : ''}>
                {player2Score}
              </div>
            </div>
            <div
              className="paddle player1"
              style={{ bottom: `${player1Position - 60}px` }} // Adjust offset to accommodate the size of the paddles
            ></div>
            <div
              className="paddle player2"
              style={{ bottom: `${player2Position - 60}px` }}
            ></div>
            <div
              className="ball"
              style={{
                display: playing ? 'block' : 'none',
                left: `${ballPosition[0] - 10}px`, // Adjust offset to accommodate the size of the ball
                bottom: `${ballPosition[1] - 10}px`
              }}
            ></div>
            <div
              className="info"
              style={{ display: playing ? 'none' : 'block' }}
            >
              Press SPACE to start a new game
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default App;

/* Maybe there's a more idiomatic way to do this, but this is
 * the arbitrary code system I chose to denote player movement.
 * The backend is coded to understand as such:
 *
 * W 0 = Player 1 starts moving up
 * W 1 = Player 1 stops moving up
 * S 0 = Player 1 starts moving down
 * S 1 = Player 1 stops moving down
 * U 0 = Player 2 starts moving up
 * U 1 = Player 2 stops moving up
 * D 0 = Player 2 starts moving down
 * D 1 = Player 2 stops moving down
 */

const registerKeyboardEvents = (socket, playing) => {
  document.onkeydown = (event) => {
    if (!event.repeat) {
      if (event.key === ' ' && !playing) {
        socket.send('START');
      } else if (event.key === 'w' && playing) {
        socket.send('W 0');
      } else if (event.key === 's' && playing) {
        socket.send('S 0');
      } else if (event.key === 'ArrowUp' && playing) {
        socket.send('U 0');
      } else if (event.key === 'ArrowDown' && playing) {
        socket.send('D 0');
      }
    }
  };

  document.onkeyup = (event) => {
    if (event.key === 'w' && playing) {
      socket.send('W 1');
    } else if (event.key === 's' && playing) {
      socket.send('S 1');
    } else if (event.key === 'ArrowUp' && playing) {
      socket.send('U 1');
    } else if (event.key === 'ArrowDown' && playing) {
      socket.send('D 1');
    }
  };
};

export default registerKeyboardEvents;

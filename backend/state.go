package main

type PlayerState struct {
	score    int
	position int
	moving   string // Can be "no", "up" or "down"
}

type BallState struct {
	position struct {
		x int
		y int
	}
	direction struct {
		x bool
		y bool
	}
}

type GameState struct {
	inProgress bool
	player1    PlayerState
	player2    PlayerState
	ball       BallState
}

/* After every point scored, the ball must reset such that:
 * 1. It randomly picks a player to move towards
 * 2. It randomly picks between upward or downward momentum
 * 3. It randomly picks a spot along the Y axis (height) to start from
 * 4. It must start just behind the visible center line relative to the
 *    direction of the player it is moving towards
 */
func (state *GameState) GenerateNewBallState() {
	var ball BallState
	if GetRandomBool() {
		ball.direction.x = true
	}
	if GetRandomBool() {
		ball.direction.y = true
	}
	if ball.direction.x {
		ball.position.x = (arenaWidth / 2) - ((ballSize - 1) / 2)
	} else {
		ball.position.x = (arenaWidth / 2) + ((ballSize - 1) / 2)
	}
	ball.position.y = GetRandomNumberInRange(ballCollisionBounds[0][1], ballCollisionBounds[1][1])
	state.ball = ball
}

// Calculate new positions for paddles and ball and then check for collision
func (state *GameState) GameTick() {
	if state.player1.moving == "up" && (state.player1.position+(2*velocityMultiplier)) < paddleMovementBounds[1] {
		state.player1.position += 2 * velocityMultiplier
	} else if state.player1.moving == "down" && (state.player1.position-(2*velocityMultiplier)) > paddleMovementBounds[0] {
		state.player1.position -= 2 * velocityMultiplier
	}
	if state.player2.moving == "up" && (state.player2.position+(2*velocityMultiplier)) < paddleMovementBounds[1] {
		state.player2.position += 2 * velocityMultiplier
	} else if state.player2.moving == "down" && (state.player2.position-(2*velocityMultiplier)) > paddleMovementBounds[0] {
		state.player2.position -= 2 * velocityMultiplier
	}

	if state.ball.direction.y {
		state.ball.position.y += 1 * velocityMultiplier
		if state.ball.position.y > ballCollisionBounds[1][1] {
			state.ball.position.y -= 2 * velocityMultiplier
			state.ball.direction.y = false
		}
	} else {
		state.ball.position.y -= 1 * velocityMultiplier
		if state.ball.position.y < ballCollisionBounds[0][1] {
			state.ball.position.y += 2 * velocityMultiplier
			state.ball.direction.y = true
		}
	}

	if state.ball.direction.x {
		state.ball.position.x += 1 * velocityMultiplier
		if state.ball.position.x > ballCollisionBounds[1][0] {
			// Check if player 2 lost a point or managed to deflect the ball
			if state.CouldNotDeflect(state.player2.position) {
				state.player1.score++
				if state.player1.score >= scoreLimit {
					endGame("1")
				} else {
					state.GenerateNewBallState()
				}
			} else {
				state.ball.position.x -= 2 * velocityMultiplier
				state.ball.direction.x = false
			}
		}
	} else {
		state.ball.position.x -= 1 * velocityMultiplier
		if state.ball.position.x < ballCollisionBounds[0][0] {
			// Check if player 1 lost a point or managed to deflect the ball
			if state.CouldNotDeflect(state.player1.position) {
				state.player2.score++
				if state.player2.score >= scoreLimit {
					endGame("2")
				} else {
					state.GenerateNewBallState()
				}
			} else {
				state.ball.position.x += 2 * velocityMultiplier
				state.ball.direction.x = true
			}
		}
	}

	if state.inProgress {
		announceState()
	}
}

/* When the ball approaches its X axis boundary, we consider that the ball was not deflected and the player lost the point if
 * 1. The upper bound of the ball is below the lower bound of their paddle, OR
 * 2. The lower bound of the ball is above the upper bound of their paddle
 */
func (state *GameState) CouldNotDeflect(paddlePosition int) bool {
	if ((state.ball.position.y + ((ballSize - 1) / 2)) < (paddlePosition - ((paddleHeight - 1) / 2))) || ((state.ball.position.y - ((ballSize - 1) / 2)) > (paddlePosition + ((paddleHeight - 1) / 2))) {
		return true
	}
	return false
}

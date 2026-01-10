//Author: Design-BAB
//Date: 1-10-2026
//Description: My classic catching game!
//Goal: Keep improving the game until it reaches 268 lines of code

package main

import (
	"math/rand/v2"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Size = 900
)

// Actor now embeds rl.Rectangle for position and size data.
type Actor struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	Flip         bool
	Speed        float32
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Flip: false, Speed: 7.0}
}

// Objects too embeds rl.Rectangle for position and size data.
type Object struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives object all the fields of rl.Rectangle (X, Y, Width, Height)
	Weight       int
}

func newObject(texture rl.Texture2D, x, y float32, weight int) *Object {
	return &Object{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Weight: weight}
}

// function when "creating" new snacks
func place(food *Object, textures *[4]rl.Texture2D) {
	newSnackPick := rand.IntN(4)
	food.Texture = textures[newSnackPick] //It needs to be a number higher because 0 <= x > y
	food.X = float32(rand.IntN(Size - 100))
	food.X = food.X + 50
	food.Y = float32(rand.IntN(100)) - 18.0
	food.Weight = newSnackPick
}

// time-dealing functions
func timeIsUp(startTime time.Time, gameDuration time.Duration) bool {
	elapsed := time.Since(startTime)
	return elapsed >= gameDuration
}
func howMuchTimeIsLeft(startTime time.Time, gameDuration time.Duration) string {
	timeDisplay := int(gameDuration.Seconds()) - int(time.Since(startTime).Seconds())
	if timeDisplay < 0 {
		timeDisplay = 0
	}
	return strconv.Itoa(timeDisplay)
}

func playTheGame(fox *Actor, snack, snack2 *Object, snackTextures [4]rl.Texture2D, score int, gameIsOver bool, startTime time.Time, gameDuration time.Duration, screenText string) (int, bool) {
	getInput(fox)
	//this will act as gravity
	fox.Y = fox.Y + 3.0
	snack.Y = snack.Y + 6.0 + float32(snack.Weight)
	snack2.Y = snack2.Y + 6.0 + float32(snack2.Weight)

	//game Logic
	if timeIsUp(startTime, gameDuration) == true {
		screenText = "Game over"
		gameIsOver = true
	} else {
		//for first snack
		if rl.CheckCollisionRecs(fox.Rectangle, snack.Rectangle) {
			//checking to see if it is a cookie, then extra points
			if snack.Texture == snackTextures[0] {
				score = score + 3
			} else {
				score++
			}
			place(snack, &snackTextures)
		}
		if snack.Y > Size {
			place(snack, &snackTextures)
		}
		rl.DrawTexture(snack.Texture, int32(snack.X), int32(snack.Y), rl.White)
		//for second snack
		if rl.CheckCollisionRecs(fox.Rectangle, snack2.Rectangle) {
			if snack2.Texture == snackTextures[0] {
				score = score + 3
			} else {
				score++
			}
			place(snack2, &snackTextures)
		}
		if snack2.Y > Size {
			place(snack2, &snackTextures)
		}
		rl.DrawTexture(snack2.Texture, int32(snack2.X), int32(snack2.Y), rl.White)
	}

	//On screen, draw text
	rl.DrawText("Your score is "+strconv.Itoa(score), 20, 20, 18, rl.DarkGray)
	if gameIsOver == false {
		rl.DrawText(howMuchTimeIsLeft(startTime, gameDuration), 525, 20, 18, rl.DarkGray)
	}
	rl.DrawText(screenText, 20, 45, 18, rl.DarkGray)
	return score, gameIsOver
}

func getInput(fox *Actor) {
	if rl.IsKeyDown(rl.KeyRight) {
		fox.X = fox.X + fox.Speed
		fox.Flip = false
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		fox.X = fox.X - fox.Speed
		fox.Flip = true
	}
	if rl.IsKeyDown(rl.KeyUp) {
		if fox.Y > 600.0 {
			fox.Y = fox.Y - fox.Speed
		} else {
			fox.Y = 600.0
		}
	}
	//collisions with the window
	fox.X = rl.Clamp(fox.X, 0.0, Size-fox.Width)
	if rl.IsKeyDown(rl.KeyRight) {
		fox.X = fox.X + fox.Speed
		fox.Flip = false
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		fox.X = fox.X - fox.Speed
		fox.Flip = true
	}
	if rl.IsKeyDown(rl.KeyUp) {
		if fox.Y > 600.0 {
			fox.Y = fox.Y - fox.Speed
		} else {
			fox.Y = 600.0
		}
	}
	//collisions with the window
	fox.X = rl.Clamp(fox.X, 0.0, Size-fox.Width)
	fox.Y = rl.Clamp(fox.Y, 0.0, Size-fox.Height)
	//flipping logic
	src := rl.NewRectangle(0, 0, float32(fox.Texture.Width), float32(fox.Texture.Height))
	dst := rl.NewRectangle(fox.X, fox.Y, float32(fox.Texture.Width), float32(fox.Texture.Height))
	origin := rl.NewVector2(0, 0)
	if fox.Flip {
		// Flip horizontally by making source width negative
		src.Width = -src.Width
		// Shift the source rect start so it doesn't disappear
		src.X = float32(fox.Texture.Width)
	}
	rl.DrawTexturePro(fox.Texture, src, dst, origin, 0, rl.White) //DrawTexturePro(texture Texture2D, sourceRec, destRec Rectangle, origin Vector2, rotation float32, tint color.RGBA)

	fox.Y = rl.Clamp(fox.Y, 0.0, Size-fox.Height)
	//flipping logic
	src1 := rl.NewRectangle(0, 0, float32(fox.Texture.Width), float32(fox.Texture.Height))
	dst1 := rl.NewRectangle(fox.X, fox.Y, float32(fox.Texture.Width), float32(fox.Texture.Height))
	origin1 := rl.NewVector2(0, 0)
	if fox.Flip {
		// Flip horizontally by making source width negative
		src1.Width = -src1.Width
		// Shift the source rect start so it doesn't disappear
		src1.X = float32(fox.Texture.Width)
	}
	rl.DrawTexturePro(fox.Texture, src1, dst1, origin1, 0, rl.White) //DrawTexturePro(texture Texture2D, sourceRec, destRec Rectangle, origin Vector2, rotation float32, tint color.RGBA)

}

func main() {
	//setting up window
	rl.InitWindow(Size, Size, "Catch The Snack!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	var startTime = time.Now()
	var gameDuration = 60 * time.Second
	screenText := "Welcome!"

	//game-play variables
	score := 0
	gameIsOver := false

	//load Background
	backgroundTexture := rl.LoadTexture("images/background.jpg")
	defer rl.UnloadTexture(backgroundTexture)

	//load fox
	foxTexture := rl.LoadTexture("images/fox.png")
	defer rl.UnloadTexture(foxTexture)
	fox := newActor(foxTexture, 100.0, 700.0)

	//Loading snacks
	var snackTextures [4]rl.Texture2D
	snackTextures[0] = rl.LoadTexture("images/cookie.png")
	defer rl.UnloadTexture(snackTextures[0])
	snackTextures[1] = rl.LoadTexture("images/orange.png")
	defer rl.UnloadTexture(snackTextures[1])
	snackTextures[2] = rl.LoadTexture("images/apple.png")
	defer rl.UnloadTexture(snackTextures[2])
	snackTextures[3] = rl.LoadTexture("images/pineapple.png")
	defer rl.UnloadTexture(snackTextures[3])
	snack := newObject(snackTextures[1], 200.0, 5.0, 1)
	snack2 := newObject(snackTextures[3], 400.0, 5.0, 1)

	for !rl.WindowShouldClose() {
		// Drawing What should be on screen
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(backgroundTexture, 0, 0, rl.White)
		if gameIsOver == false {
			score, gameIsOver = playTheGame(fox, snack, snack2, snackTextures, score, gameIsOver, startTime, gameDuration, screenText)
		} else if rl.IsKeyDown(rl.KeyY) && gameIsOver == true {
			gameIsOver = false
			startTime = time.Now()
			score = 0
		}
		if gameIsOver == true {
			rl.DrawText("Your score is "+strconv.Itoa(score), 100, 100, 24, rl.DarkGray)
			rl.DrawText("Press Y to play again", 100, 200, 24, rl.DarkGray)
			rl.DrawTexture(fox.Texture, int32(fox.X), int32(fox.Y), rl.White)
		}
		rl.EndDrawing()
	}
}

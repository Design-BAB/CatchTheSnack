//Author: Design-BAB
//Date: 1-10-2026
//Description: My classic catching game!
//Goal: Keep improving the game until it reaches 268 lines of code

package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	_ "github.com/glebarez/go-sqlite"
)

const (
	Size                = 900
	gravity     float32 = 0.7
	jumpImpulse float32 = 15.5
)

type GameState struct {
	Score         int
	IsOver        bool
	HighScore     []*Scoreboard
	ScoreRecorded bool
}

func newGame() *GameState {
	initialScores := []*Scoreboard{}
	return &GameState{HighScore: initialScores}
}

// Actor now embeds rl.Rectangle for position and size data.
type Actor struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
	Flip         bool
	Speed        float32
	VelY         float32
	OnGround     bool
}

func newActor(texture rl.Texture2D, x, y float32) *Actor {
	return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Flip: false, Speed: 7.0}
}

// Objects too embeds rl.Rectangle for position and size data.
type Object struct {
	Texture rl.Texture2D
	//this is the collision box``
	rl.Rectangle  // This gives object all the fields of rl.Rectangle (X, Y, Width, Height)
	Weight        int
	LastCatchTime time.Time
}

func newObject(texture rl.Texture2D, x, y float32, weight int) *Object {
	startTimeNow := time.Now()
	return &Object{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: float32(texture.Width), Height: float32(texture.Height)}, Weight: weight, LastCatchTime: startTimeNow}
}

type Scoreboard struct {
	Name  string
	Score int
}

func newScoreToBoard(name string, score int) *Scoreboard {
	return &Scoreboard{Name: name, Score: score}
}

func CreateTable(db *sql.DB) (sql.Result, error) {
	sqlCommand := `CREATE TABLE IF NOT EXISTS scores (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		score INTEGER NOT NULL);`
	return db.Exec(sqlCommand)
}

func updateHighScore(db *sql.DB, yourGame GameState) ([]*Scoreboard, error) {
	var results []*Scoreboard
	//gonna add the current score into the data base
	now := time.Now()
	dateOfToday := now.Format("Monday, January 2, 2006")
	_, err := db.Exec(`INSERT INTO scores (name, score) VALUES (?, ?)`, dateOfToday, yourGame.Score)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rows, err := db.Query("SELECT name, score FROM scores ORDER BY score DESC LIMIT 3;")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var score int
		err := rows.Scan(&name, &score)
		if err != nil {
			log.Println("Error Scanning row: ", err)
			continue
		}
		results = append(results, newScoreToBoard(name, score))
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func displayHighScore(yourGame *GameState) {
	if len(yourGame.HighScore) == 0 {
		rl.DrawText("No high scores yet!", 190, 300, 20, rl.DarkGray)
		return
	}

	if yourGame.Score > yourGame.HighScore[0].Score {
		rl.DrawText("Wow! You made the high score!", 190, 250, 20, rl.DarkGray)
		rl.DrawText(strconv.Itoa(yourGame.Score), 190, 300, 20, rl.DarkGray)
	} else {
		var y int32 = 300
		for i, theScoreToDisplay := range yourGame.HighScore {
			mssg := strconv.Itoa(theScoreToDisplay.Score) + "    " + theScoreToDisplay.Name
			rl.DrawText(mssg, 190, y, 20, rl.DarkGray)
			y = int32(i+1)*20 + 300
		}
	}
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

func playTheGame(fox *Actor, snack, snack2 *Object, snackTextures [4]rl.Texture2D, startTime time.Time, gameDuration time.Duration, screenText string, crunchFx rl.Sound, yourGame *GameState) {
	getInput(fox)
	updateFoxPhysics(fox)
	snack.Y = snack.Y + 6.0 + float32(snack.Weight)
	snack2.Y = snack2.Y + 6.0 + float32(snack2.Weight)

	//game Logic
	if timeIsUp(startTime, gameDuration) == true {
		screenText = "Game over"
		yourGame.IsOver = true
	} else {
		//for first snack
		if rl.CheckCollisionRecs(fox.Rectangle, snack.Rectangle) {
			rl.PlaySound(crunchFx)
			//checking to see if it is a cookie, then extra points
			if snack.Texture == snackTextures[0] {
				yourGame.Score += 3
			} else {
				yourGame.Score++
			}
			place(snack, &snackTextures)
		}
		if snack.Y > Size {
			place(snack, &snackTextures)
		}
		rl.DrawTexture(snack.Texture, int32(snack.X), int32(snack.Y), rl.White)
		//for second snack
		if rl.CheckCollisionRecs(fox.Rectangle, snack2.Rectangle) {
			rl.PlaySound(crunchFx)
			if snack2.Texture == snackTextures[0] {
				yourGame.Score += 3
			} else {
				yourGame.Score++
			}
			place(snack2, &snackTextures)
		}
		if snack2.Y > Size {
			place(snack2, &snackTextures)
		}
		rl.DrawTexture(snack2.Texture, int32(snack2.X), int32(snack2.Y), rl.White)
	}

	//On screen, draw text
	rl.DrawText("Your score is "+strconv.Itoa(yourGame.Score), 20, 20, 18, rl.DarkGray)
	if yourGame.IsOver == false {
		rl.DrawText(howMuchTimeIsLeft(startTime, gameDuration), 525, 20, 18, rl.DarkGray)
	}
	rl.DrawText(screenText, 20, 45, 18, rl.DarkGray)
	drawFox(fox)
}

func getInput(fox *Actor) {
	//collisions with the window
	fox.X = rl.Clamp(fox.X, 0.0, Size-fox.Width)
	//Controls for the fox
	if rl.IsKeyDown(rl.KeyRight) {
		fox.X = fox.X + fox.Speed
		fox.Flip = false
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		fox.X = fox.X - fox.Speed
		fox.Flip = true
	}
	if (rl.IsKeyPressed(rl.KeyUp) || rl.IsKeyPressed(rl.KeySpace)) && fox.OnGround {
		fox.VelY = -jumpImpulse
		fox.OnGround = false
	}
	//collisions with the window
	fox.X = rl.Clamp(fox.X, 0.0, Size-fox.Width)
}

func updateFoxPhysics(fox *Actor) {
	fox.VelY += gravity
	fox.Y += fox.VelY
	groundY := float32(Size) - fox.Height
	if fox.Y >= groundY {
		fox.Y = groundY
		fox.VelY = 0
		fox.OnGround = true
	} else {
		fox.OnGround = false
	}
	fox.Y = rl.Clamp(fox.Y, 0.0, groundY)
}

func drawFox(fox *Actor) {
	src := rl.NewRectangle(0, 0, float32(fox.Texture.Width), float32(fox.Texture.Height))
	dst := rl.NewRectangle(fox.X, fox.Y, float32(fox.Texture.Width), float32(fox.Texture.Height))
	origin := rl.NewVector2(0, 0)
	if fox.Flip {
		// Flip horizontally by making source width negative
		src.Width = -src.Width
		// Shift the source rect start so it doesn't disappear
		src.X = float32(fox.Texture.Width)
	}
	rl.DrawTexturePro(fox.Texture, src, dst, origin, 0, rl.White)
}

func main() {
	//setting up window
	rl.InitWindow(Size, Size, "Catch The Snack!")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)
	db, err := sql.Open("sqlite", "./score.db?_pragma=foreign_keys(1)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	//create the table just in case
	_, err = CreateTable(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	var startTime = time.Now()
	var gameDuration = 60 * time.Second
	screenText := "Welcome!"
	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	//game-play variables
	//score := 0
	//gameIsOver := false
	yourGame := newGame()
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
	snack2 := newObject(snackTextures[3], 400.0, 5.0, 3)

	//load sound
	crunchFx := rl.LoadSound("sounds/crunch.wav")
	defer rl.UnloadSound(crunchFx)

	for !rl.WindowShouldClose() {
		// Drawing What should be on screen
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		rl.DrawTexture(backgroundTexture, 0, 0, rl.White)
		if yourGame.IsOver == false {
			playTheGame(fox, snack, snack2, snackTextures, startTime, gameDuration, screenText, crunchFx, yourGame)
		} else if rl.IsKeyDown(rl.KeyY) && yourGame.IsOver == true {
			yourGame = newGame()
			startTime = time.Now()
		}
		if yourGame.IsOver {
			if yourGame.ScoreRecorded == false {
				var err error
				yourGame.HighScore, err = updateHighScore(db, *yourGame)
				if err != nil {
					log.Println("Arrgh, the database was unreachable")
				}
				yourGame.ScoreRecorded = true
			}
			displayHighScore(yourGame)
			rl.DrawText("Your score is "+strconv.Itoa(yourGame.Score), 100, 100, 24, rl.DarkGray)
			rl.DrawText("Press Y to play again", 100, 200, 24, rl.DarkGray)
			drawFox(fox)
		}
		rl.EndDrawing()
	}
}

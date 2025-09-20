package main

import (
  rl "github.com/gen2brain/raylib-go/raylib"
  "strconv"
  "time"
  "math/rand/v2"
)


// Actor now embeds rl.Rectangle for position and size data.
type Actor struct {
    Texture rl.Texture2D
    //this is the collision box``
    rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
    Flip bool
    Speed float32
}
func newActor(texture rl.Texture2D, x, y, width, height float32) *Actor {
  return &Actor{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: width, Height: height}, Flip: false, Speed: 7.0}
}


// Actor now embeds rl.Rectangle for position and size data.
type Object struct {
    Texture rl.Texture2D
    //this is the collision box``
    rl.Rectangle // This gives Actor all the fields of rl.Rectangle (X, Y, Width, Height)
}
func newObject(texture rl.Texture2D, x, y, width, height float32) *Object {
  return &Object{Texture: texture, Rectangle: rl.Rectangle{X: x, Y: y, Width: width, Height: height}}
}
func place(c *Object, size int)  {
  c.X = float32(rand.IntN(size - 20))
  c.Y = 2.0
}


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

func main() {
  //not in the book but, I need to set up the window and fps
  var size int32 = 800
  myGreen := rl.NewColor(143, 232, 102, 255)
  rl.InitWindow(size, size, "Catch The Snack!")
  defer rl.CloseWindow()
  rl.SetTargetFPS(60)
  var startTime = time.Now()
  var gameDuration = 60 * time.Second
  screenText := "Welcome!"



  //at this point i am following the book as close as i can
  //variables
  score := 0
  gameIsOver := false

  //load fox
  foxTexture := rl.LoadTexture("images/fox.png")
  defer rl.UnloadTexture(foxTexture)
  fox := newActor(foxTexture, 100.0, float32(size) - 100, float32(foxTexture.Width), float32(foxTexture.Height))
  //laod coin
  coinTexture := rl.LoadTexture("images/coin.png")
  defer rl.UnloadTexture(coinTexture)
  coin := newObject(coinTexture, 200.0, 50.0, float32(coinTexture.Width), float32(coinTexture.Height))

    
  for !rl.WindowShouldClose() {
    // Drawing What should be on screen
    rl.BeginDrawing()
    rl.ClearBackground(myGreen)
    //this line is like fox.draw()
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
    //this will act as gravity
    fox.Y = fox.Y + 3.0
    coin.Y = coin.Y + 7.0
    //collisions with the window
    fox.X = rl.Clamp(fox.X, 0.0, float32(size) - fox.Width)
    fox.Y = rl.Clamp(fox.Y, 0.0, float32(size) - fox.Height)
   
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
    if timeIsUp(startTime, gameDuration) == true {
      screenText = "Game over"
      gameIsOver = true
    } else {
      if rl.CheckCollisionRecs(fox.Rectangle, coin.Rectangle){
        score++
        place(coin, int(size))
      }
      if coin.Y > float32(size) {
        place(coin, int(size))
      }
      rl.DrawTexture(coin.Texture, int32(coin.X), int32(coin.Y), rl.White)
  }
    //place(coin, int(size))
    //On screen, draw text
    rl.DrawText("Your score is " + strconv.Itoa(score), 20, 20, 18, rl.DarkGray)
    if gameIsOver == false {
      rl.DrawText(howMuchTimeIsLeft(startTime, gameDuration), 525, 20, 18, rl.DarkGray)
    }
    rl.DrawText(screenText, 20, 45, 18, rl.DarkGray)
    rl.EndDrawing()
  }
}

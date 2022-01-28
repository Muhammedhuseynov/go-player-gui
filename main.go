package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"log"
	"os"
	"time"
)

//init app
var _app = app.New()

//create window
var wind = _app.NewWindow("GoMp3Player")

//current playing song name
var currentSong = binding.NewString()

var timeOfMusic = binding.NewString()
var pause = false

//store songs into map{"song_name":"song_path"}
var songNameAndPath = make(map[string]string)
var songtime string

//song list buttons
var radioBtn = widget.NewRadioGroup([]string{},
	func(val string) {
		currentSong.Set(val)
	})

//add song list button in vertical scroll layout
var musicListContainer = container.NewVScroll(radioBtn)

//to make object to center
func makeObjCenter(obj fyne.CanvasObject) *fyne.Container {
	return container.New(layout.NewHBoxLayout(), layout.NewSpacer(), obj, layout.NewSpacer())
}

//function to add song
func songAdder() {
	//folder opener
	folder := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		// if clicked cancel Button (handle crush)
		if uri == nil {
			return
		} else {
			//get files in folder
			files, _ := uri.List()
			for _, file := range files {
				//filter file extensions
				ext := storage.ExtensionFileFilter{[]string{".mp3"}}
				if ext.Matches(file) == true {
					songNameAndPath[file.Name()] = file.Path()
					radioBtn.Options = append(radioBtn.Options, file.Name())
				}
			}
			//refresh list of buttons when added songs
			radioBtn.Refresh()
		}

	}, wind)
	//open folder opener window
	folder.Show()

}

//track song index
var currenSongInd int

//song player
func RunSong(songNamePath string) {
	//get file
	f, err := os.Open(songNamePath)
	//if not chosen file then show error
	if err != nil {
		dialog.NewInformation("Warning", "Please First add the song!", wind).Show()
	} else {
		streamer, format, _ := mp3.Decode(f)
		//get song time len
		songlen := format.SampleRate.D(streamer.Len()).Round(time.Second)
		songtime = songlen.String()

		//change song time(label)
		timeOfMusic.Set(songtime)

		done := make(chan bool)
		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		// play
		speaker.Play(beep.Seq(
			streamer, beep.Callback(func() {
				done <- true
			})))

		////================================ fixxxxx track time song
		// Uncomment below | if you uncomment bellow,commment the above code
		//go func() {
		//	for {
		//		select {
		//		case <-done:
		//			return
		//		case <-time.After(time.Second):
		//			speaker.Lock()
		//			//log.Println(format.SampleRate.D(streamer.Position()).Round(time.Second).String())
		//			timeOfMusic.Set(songtime + "/" + format.SampleRate.D(streamer.Position()).Round(time.Second).String())
		//			//log.Println(format.SampleRate.D(streamer.Position()).Round(time.Second))
		//			speaker.Unlock()
		//		}
		//	}
		//}()
		//================================
		// looping
		//select {}
	}
}

//need for button function
func songPlay() {

	currenSongIndPtr := &currenSongInd
	//get current playing song name
	currentSongNAme, _ := currentSong.Get()
	//find current playing song index
	for ind, elem := range radioBtn.Options {
		if elem == currentSongNAme {
			*currenSongIndPtr = ind
		}
	}
	//play the song
	RunSong(songNameAndPath[currentSongNAme])
}

func nextSong() {
	if len(radioBtn.Options) != 0 {
		//get songlists
		songLists := radioBtn.Options
		//get last index
		lastInd := len(songLists) - 1
		currenSongIndPtr := &currenSongInd
		if currenSongInd < lastInd {
			//next song
			*currenSongIndPtr += 1
			//next song name
			currentSongNAme := songLists[currenSongInd]
			RunSong(songNameAndPath[currentSongNAme])
			//change radio button select current playing song
			radioBtn.SetSelected(radioBtn.Options[currenSongInd])
		} else {
			//if current index greater than last index ,then do not change anything keep on current index
			currentSongNAme := songLists[currenSongInd]
			RunSong(songNameAndPath[currentSongNAme])
			radioBtn.SetSelected(radioBtn.Options[currenSongInd])
			*currenSongIndPtr = currenSongInd

		}
	} else {
		dialog.NewInformation("Warning", "Not selected Song", wind).Show()
	}
}

func prevSong() {
	if len(radioBtn.Options) != 0 {
		songLists := radioBtn.Options
		lastInd := len(songLists) - 1
		currenSongIndPtr := &currenSongInd
		if currenSongInd <= lastInd {
			//keep it 0 / otherwise it will be -1
			if currenSongInd == 0 {
				*currenSongIndPtr = 0
			} else {
				*currenSongIndPtr -= 1
			}

			currentSongNAme := songLists[currenSongInd]
			RunSong(songNameAndPath[currentSongNAme])
			radioBtn.SetSelected(radioBtn.Options[currenSongInd])

		} else {
			//if current index greater than last index ,then do not change anything keep on current index
			currentSongNAme := songLists[0]
			RunSong(songNameAndPath[currentSongNAme])
			radioBtn.SetSelected(radioBtn.Options[currenSongInd])
			*currenSongIndPtr = 0
		}
	} else {
		dialog.NewInformation("Warning", "Not selected Song", wind).Show()
	}
}

func main() {

	wind.Resize(fyne.NewSize(500, 449))
	wind.CenterOnScreen()

	//Container with Songs--------
	musicListContainer.SetMinSize(fyne.NewSize(300, 300))

	timeOfMusic.Set("0m00s")
	timeOfMusicCentered := makeObjCenter(widget.NewLabelWithData(timeOfMusic))
	currentSong.Set("Current Song")
	currentPlayingSongCentered := makeObjCenter(widget.NewLabelWithData(currentSong))

	//button for add song----------
	adderSong := widget.NewButton("+", songAdder)

	playSong := widget.NewButton("Play", songPlay)
	nextSong := widget.NewButton("Next", nextSong)
	prevSong := widget.NewButton("Prev", prevSong)

	pauseAndUnPause := widget.NewButton("Un/Pause", func() {
		log.Println([]fyne.CanvasObject{})
		if pause != true {
			pause = true
			// you cannot play another song when paused so
			// Need to Disable Play,next,prec  buttons if not disable , app stops
			playSong.Disable()
			nextSong.Disable()
			prevSong.Disable()
			speaker.Lock()

		} else {
			speaker.Unlock()
			pause = false
			playSong.Enable()
			nextSong.Enable()
			prevSong.Enable()
		}

		log.Println("Un/Pause")
	})

	//centered Container
	btnContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), playSong, pauseAndUnPause, prevSong, nextSong, layout.NewSpacer())
	wind.SetContent(container.NewVBox(musicListContainer, currentPlayingSongCentered, timeOfMusicCentered, adderSong, btnContainer))
	wind.ShowAndRun()
}

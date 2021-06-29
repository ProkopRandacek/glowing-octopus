package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	scriptFolder = "./master/script-output/"
)

func (b *Bot) walkTo(p Position) error {
	b.State.Walking = true
	_, err := b.conn.Execute(fmt.Sprintf(`/walkto [%.2f,%.2f]`, p.X, p.Y))
	return err
}

func (b *Bot) waitForTaskDone() { // Waits until task is done.
	for {
		fmt.Println("Waiting for task done")
		time.Sleep(2 * time.Second)
		b.refreshState()
		if !(b.State.Walking || b.State.Mining || b.State.Placing || b.State.Clearing || b.State.Building) {
			break
		}
	}
}

func (b *Bot) getResources(box Box) ([][]Position, error) {
	filename := scriptFolder + "resrc.json"
	os.Remove(filename)
	_, err := b.conn.Execute(fmt.Sprintf("/writeresrc [[%.2f,%.2f],[%.2f,%.2f]]", box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
	if err != nil {
		return nil, err
	}

	for { // wait until the file is written
		_, err := os.Stat(filename)
		if err == nil {
			break
		}

		log(fmt.Sprintf("Waiting for the %s to be generated", filename))
		time.Sleep(time.Second)
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var resrc [][]Position
	json.Unmarshal(dat, &resrc)

	b.Mapper.LoadedBoxes = append(b.Mapper.LoadedBoxes, box)

	return resrc, nil
}

func (b *Bot) drawBox(box Box, color Color) {
	b.conn.Execute(fmt.Sprintf(`/drawbox {"color":[%2.f, %2.f, %2.f, 0.2],"x1":%2.f,"y1":%2.f,"x2":%2.f,"y2":%2.f}`, color.R, color.G, color.B, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *Bot) craft(r string, c int) {
	b.conn.Execute(fmt.Sprintf(`/craft {"recipe":"%s","count":%d}`, r, c))
}

func (b *Bot) mine(p Position) {
	b.conn.Execute(fmt.Sprintf(`/mine [%.2f,%.2f]`, p.X, p.Y))
}

func (b *Bot) build(p Position, item string) {
	b.conn.Execute(fmt.Sprintf(`/build {"pos":[%2.f,%2.f],"item":"%s"}`, p.X, p.Y, item))
}

func (b *Bot) buildDir(p Position, item string, dir int) {
	b.conn.Execute(fmt.Sprintf(`/build {"pos":[%2.f,%2.f],"item":"%s","dir":%d}`, p.X, p.Y, item, dir))
}

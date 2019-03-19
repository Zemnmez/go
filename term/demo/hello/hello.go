package main

import (
	"github.com/nsf/termbox-go"
	"time"

	"zemn.me/term"
)

const lipsum = `
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin magna arcu, dignissim eu ullamcorper at, lobortis eget felis. Vivamus tristique iaculis sem. Mauris vehicula eros consectetur molestie condimentum. Ut in efficitur lacus. Morbi viverra porta libero eu lobortis. Aenean a ipsum id lorem tempus vehicula vel a nunc. Integer quis ullamcorper odio, id elementum diam. Phasellus efficitur, arcu sed gravida feugiat, sapien ipsum lacinia nulla, nec elementum elit nisl ac metus.

Quisque sollicitudin vitae velit id tincidunt. Integer finibus nibh volutpat facilisis aliquam. Duis nec ante nibh. Maecenas lobortis vehicula augue vel ultricies. Nulla gravida aliquam metus, non interdum ex volutpat a. Donec placerat, enim et rhoncus aliquet, neque tellus sollicitudin mi, quis bibendum libero odio quis quam. Praesent pharetra tortor in neque gravida auctor.

Proin semper neque et felis maximus, id egestas risus commodo. Sed eu leo sit amet magna mollis semper. Donec luctus pretium risus, a lacinia magna accumsan in. Curabitur dignissim ipsum at magna egestas, posuere viverra metus ornare. Sed id justo eget ipsum suscipit congue. Proin lorem est, euismod at ex non, gravida imperdiet magna. Maecenas suscipit nisl neque, sed fermentum metus interdum vel. Sed molestie iaculis diam in dictum. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur sed finibus libero. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Fusce in sapien accumsan justo ullamcorper mollis. Aliquam erat volutpat. Maecenas tincidunt libero vitae luctus lacinia. Pellentesque commodo ac elit eu sodales. Curabitur pellentesque bibendum augue, vitae ultrices justo auctor id.

Sed sollicitudin elit et mi faucibus, non dapibus ligula rutrum. Aenean eu nunc quis leo iaculis euismod vitae vitae nulla. Nulla justo diam, pulvinar ut suscipit et, porta sit amet eros. Phasellus imperdiet, sem sed commodo pellentesque, nulla orci feugiat justo, vel euismod turpis purus quis tellus. Proin non lorem fringilla, euismod odio aliquam, euismod nibh. Cras id tempus lectus. Vivamus metus lacus, elementum vitae justo blandit, venenatis scelerisque dui. Vivamus vitae leo tempor, vehicula risus sed, faucibus dui. Nulla sit amet enim justo.
`

func main() {
	c, done, err := term.NewCanvas()
	defer done()
	if err != nil {
		panic(err)
	}

	_, err = term.Text(lipsum).Render(c)
	if err != nil {
		return
	}

	termbox.Flush()

	<-time.After(3 * time.Second)
}

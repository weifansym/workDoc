## Programatically generate gorgeous social media images in Go.
First impressions are important and one of the first things people see when someone shares your article is the social image.

In the worst case, it is a blank image because something went wrong or the meta tags weren’t set properly. In the best case, it is a hand made
graphic that feels luxurious somehow, and shows our users that a lot of effort has gone into producing the content.

Twitter also tells us that “Tweets with photos receive an average 35% boost in Retweets” (see [What fuels a Tweet’s engagement?](https://blog.x.com/official/en_us/a/2014/what-fuels-a-tweets-engagement.html))

Doing a custom image for each post or sharable page in your app takes a lot of hard work.

But wait, I forgot. We’re programmers. We’ll automate it.

### The goal
We’re going to look at how we can programatically create social images like these:
![image](https://github.com/weifansym/workDoc/assets/6757408/e754f488-ee9d-4180-b348-d66216019099)

These are real images that were generated for [the Pace blog](https://pace.dev/blog) using code from this post.

The images are hopefully attractive, informative and unique.
#### Render images with the standard library
The standard library is very low level. We’ll explore it a little here, but later I’m going to use [Michael Fogleman’s gg package](https://pkg.go.dev/github.com/fogleman/gg?utm_source=godoc), which provides an abstraction and much simpler API.

The standard library provides lots of packages that deal with rendering images and fonts.

The image/draw package provides a simple but powerful function called Draw:
```
func Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point, op Op)
```
This function lets you copy pixels from one image to another (specifically onto a draw.Image type).

#### Rectangles
The [image.Rect](https://pkg.go.dev/image#Rect) function creates an image.Rectangle describing an area on the image.

Here is the source code for the Rectangle type:
```
// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
//
// A Rectangle is also an Image whose bounds are the rectangle itself. At
// returns color.Opaque for points in the rectangle and color.Transparent
// otherwise.
type Rectangle struct {
	Min, Max Point
}
```
It contains two Point types:
```
// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y int
}
```
With these two structures, we can describe a 2D rectange.

When working with images in Go, you will spend a lot of time working with boxes using this type, so it’s worth a quick overview of how it works.

The rectangle holds a start position (x0, y0) and an end position (x1, y1). Notice the rectangle does not hold a width and height, instead 
the second pair of coordinates desctibe the end point.

A box that is 100 x 100 pixels might be described like this:
```
image.Rect(0, 0, 100, 100)
```
or a 100 x 100 box might look like this:
```
image.Rect(100, 100, 200, 200)
```
A 20 x 20 box in the middle of that could be described like this:
```
image.Rect(40, 40, 60, 60)
```
### Drawing solid rectangles
To give you an example of how low-level rendering with the standard library is, let’s have a quick look at how we might draw a filled red rectangle onto our image.

An image.Image can be a uniform colour if we use image.Uniform type, from which we can copy to draw solid rectangles.
```
redImage := &image.Uniform{color.RGBA{0xFF, 0x00, 0x00, 0xFF}}
draw.Draw(img, image.Rect(10, 10, 30, 30), &image.Uniform{blue}, image.ZP, draw.Src)
```
* To see why it is designed like this, take a look at the [image.Image interface](https://pkg.go.dev/image#Image). It describes the colour model and the size of the image, but the only way to read data from the image is via the At(x, y int) color.Color method, which reads the colour at a single pixel. Pretty low level, right?

### Writing text
Social images often contain the title of the article or page, so this means we need to render text onto an draw.Image.

The [https://github.com/golang/freetype](https://github.com/golang/freetype) package is the font rasteriser for Go, and after a quick glance at the [github.com/golang/freetype/raster](https://pkg.go.dev/github.com/golang/freetype/raster) package, you will have a new appreciation for how fonts work.

We can create a NewContext and use DrawString to write some text.
```
var (
	img 	 = image.NewRGBA(image.Rect(0, 0, 320, 240))
	x, y 	 = 50, 50
	fontSize = 12.0
	label 	 = "Hi there"
)

ctx := freetype.NewContext()
dc.SetDst(img)
pt := freetype.Pt(x, y+int(c.PointToFixed(fontSize)>>6))
if _, err := dc.DrawString(label, pt); err != nil {
	return err
}
```
As you can see, it’s a very different proposition than creating an HTML file with <p>Hi there</p> in it.

### Encoding and saving image files
To save an image.Image as a usable file, we need to encode it to JPG or PNG (or GIF if you want to make animated images).

We do this using the image/jpg, image/png, or image/gif packages and you usually write code like this:
```
func SavePNG(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
```
### Render images using gg package
[Michael Fogleman’s gg package](https://pkg.go.dev/github.com/fogleman/gg) provides an abstraction on top the standard library that provides a very simple, programmer friendly API with [lots of examples](https://github.com/fogleman/gg/tree/master/examples) for us to poke around and take code snippets from.

The API is nicely designed, and allows us to write much more readable code. For example, to draw a red circle and save it as a PNG, we can write code like this:
```
dc := gg.NewContext(1000, 1000) 		// canvas 1000px by 1000px
dc.DrawCircle(500, 500, 400) 			// a circle in the middle
dc.SetRGB(0xFF, 0x00, 0x00) 			// choose colour
dc.Fill() 								// fill the circle
err := dc.SavePNG("out.png") 			// save it
if err != nil {
	return errors.Wrap(err, "save png")
}
```
This is much simpler than doing so using just the standard library, so we’ll use gg to render our social images.

### Choose a size for the image
When we create a new gg.Context we specify the width and height in pixels.

At the time of writing, the [best advice I could find on the optimum size for social images on louisem.com](https://louisem.com/217438/twitter-image-size) suggested a size for mobile of 1200 x 628 (although it seems to change a lot, so it’s worth checking before you decide).

1200 x 628 is astonishing considering I remember programming on an Amiga with a screen resolution of 320×200.
### Write a program to generate the images
Now we’re ready to start our program.

We’ll create a gg.Context after some Go boilerplate:
```
import (
	"fmt"
	"os"

	"github.com/fogleman/gg"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	dc := gg.NewContext(1200, 628)
	// todo, more programming
}
```
We’ll build up our image generator by adding code to the run function.

### Render an image file
All of our blog posts have a primary image. We will load it, resize it, and use it as the background for our social image.

User experience experts will no doubt give you lots of reasons why having these images the same is important, but I know that I like the fact that the user can tap something on Twitter or Facebook, and then land immediately on the right page, and havinng the same images on both helps us tell that story.

So assuming you have a nice image file (as backgroundImageFilename), we’ll load it and draw it into the gg.Context.
```
backgroundImage, err := gg.LoadImage(backgroundImageFilename)
if err != nil {
	return errors.Wrap(err, "load background image")
}
dc.DrawImage(backgroundImage, 0, 0)
```
This article assumes the images are the same size. If they are not, you will need to consider using the github.com/disintegration/imaging package to resize them to fill the entire image.
```
backgroundImage = imaging.Fill(backgroundImage, dc.Width(), dc.Height(), imaging.Center, imaging.Lanczos)
```
### Save the image file
Getting near-to-live feedback as we code is especially important for visual tasks like these. So we’ll dedicate a little effort into creating the output image so we can see our images evolve as we go.

The gg package makes this very easy for us.
```
if err := dc.SavePNG(outputFilename); err != nil {
	return errors.Wrap(err, "save png")
}
```
Running this code will produce an image like this:
![image](https://github.com/weifansym/workDoc/assets/6757408/b9b2a40e-da09-4bd6-b0aa-b085ce23ae2a)

We’re going to add some more layers to our image, so remember to keep the SavePNG code at the end so your changes aren’t ignored. (Anything you draw onto the image after savinng it won’t be seen.)

### Add a semi-transparent overlay
To ensure our text can be seen, we’re also going to add a semi-transparent black rectangle overlay too.

We’ll leave a margin to provide a full-colour border effect.
```
margin := 20.0
x := margin
y := margin
w := float64(dc.Width()) - (2.0 * margin)
h := float64(dc.Height()) - (2.0 * margin)
dc.SetColor(color.RGBA{0, 0, 0, 204})
dc.DrawRectangle(x, y, w, h)
dc.Fill()
```
After calculating the dimensions of the overlay, we set the colour to color.RGBA{0, 0, 0, 204} which is black at 80% opacity (204 is 80% of 255).

We draw the rectangle and fill it.
![image](https://github.com/weifansym/workDoc/assets/6757408/63e2783f-4088-484b-a55a-62c99ba85952)

### Add text
Since these images are going to be used in a blog context, we are going to show the post title as the primary text on the image.

Your users will want to know where the content comes from, so we’ll write our brand logo to the bottom right hand corner, and the blog domain on the left.

Let’s start with the branding logo. We’ll need to load a font, and calculate a rectangle to write the text into.
```
fontPath := filepath.Join("fonts", "OpenSans-Bold.ttf")
if err := dc.LoadFontFace(fontPath, 80); err != nil {
	return errors.Wrap(err, "load font")
}
dc.SetColor(color.White)
s := "PACE."
marginX := 50.0
marginY := -10.0
textWidth, textHeight := dc.MeasureString(s)
x = float64(dc.Width()) - textWidth - marginX
y = float64(dc.Height()) - textHeight - marginY
dc.DrawString(s, x, y)
```
You can download a font for free from[Google Fonts](https://fonts.google.com/), this articles uses the [Open Sans](https://fonts.google.com/specimen/Open+Sans) font.
Since all coordinates originate in the top left hand corner, if we want to use the other dimensions (the right or bottom edges) we’ll need to do some simple calculations.

We can measure the rectangle that our text will take up by callinng dc.MeasureString and subtracting that from the image width and height. This approach has the additional benefit of being dynamic; if you change the image size, this code will still work.

The marginX and marginnY values allow us to fine tune the position of the logo.

We use dc.DrawString again to draw the text:
![image](https://github.com/weifansym/workDoc/assets/6757408/31025b4b-8425-4870-81fd-866fecb7b02e)

Next we’ll use the same technique to draw the domain name onto the bottom left edge of the image:
```
textColor := color.White
fontPath = filepath.Join("fonts", "Open_Sans", "OpenSans-Bold.ttf")
if err := dc.LoadFontFace(fontPath, 60); err != nil {
	return errors.Wrap(err, "load Open_Sans")
}
r, g, b, _ := textColor.RGBA()
mutedColor := color.RGBA{
	R: uint8(r),
	G: uint8(g),
	B: uint8(b),
	A: uint8(200),
}
dc.SetColor(mutedColor)
marginY = 30
s = "https://pace.dev/"
_, textHeight = dc.MeasureString(s)
x = 70
y = float64(dc.Height()) - textHeight - marginY
dc.DrawString(s, x, y)
```
Here we create a semi-transparent colour for the text, based on a textColor variable that you can control.

Again we measure the height of the text for positioning, and use dc.DrawString to draw the text.
![image](https://github.com/weifansym/workDoc/assets/6757408/5d84d345-6fa9-4130-862e-da6ae589e279)

Finally, we’ll add the title. This time we will add a text shadow (a black copy of the text drawn underneath the white one) so it pops off the image.
```
title := "Programatically generate these gorgeous social media images in Go"
textShadowColor := color.Black
textColor = color.White
fontPath = filepath.Join("fonts", "Open_Sans", "OpenSans-Bold.ttf")
if err := dc.LoadFontFace(fontPath, 90); err != nil {
	return errors.Wrap(err, "load Playfair_Display")
}
textRightMargin := 60.0
textTopMargin := 90.0
x = textRightMargin
y = textTopMargin
maxWidth := float64(dc.Width()) - textRightMargin - textRightMargin
dc.SetColor(textShadowColor)
dc.DrawStringWrapped(title, x+1, y+1, 0, 0, maxWidth, 1.5, gg.AlignLeft)
dc.SetColor(textColor)
dc.DrawStringWrapped(title, x, y, 0, 0, maxWidth, 1.5, gg.AlignLeft)
```
![image](https://github.com/weifansym/workDoc/assets/6757408/6996bbae-fc0a-4956-ad54-72c203147f15)

We’ll use the same technique to add some more bits and pieces, but I suppose you get the picture by now.

### Making a GIF
An animated GIF is just a multi-frame image. We can use the code we’ve written today to render two or more different versions of the image (called frame1 annd frame2), and stitch them into a GIF like this:
```
palettedImage1 := image.NewPaletted(frame1.Bounds(), palette.Plan9)
draw.FloydSteinberg.Draw(palettedImage1, frame1.Bounds(), frame1, image.ZP)
palettedImage2 := image.NewPaletted(frame2.Bounds(), palette.Plan9)
draw.FloydSteinberg.Draw(palettedImage2, frame2.Bounds(), frame2, image.ZP)
f, err := os.Create("/path/to/social-image.gif")
if err != nil {
	return errors.Wrap(err, "create gif file")
}
gif.EncodeAll(f, &gif.GIF{
	Image: []*image.Paletted{
		palettedImage1,
		palettedImage2,
	},
	Delay: []int{50, 50},
})
```
I’ll leave it to you to experiment with this to see what cool effects you can achieve.

In our case, we use this technique to put a pipe character | at the end of the title, which gives the impression of a code editor, which we think will appeal to a technical audience. Plus, sometimes we do things just because we like it :)

### The final image
So here is our final animated image:
![image](https://github.com/weifansym/workDoc/assets/6757408/5f855656-ffe4-4ae8-be78-bb1f084b5107)

Now, on Twitter, our blog posts should look attractive, branded, and clear.

### Set the metadata in the HTML page
For this whole thing to work, we need to set this image as the og:image metadata,
```
<meta property='og:image' content='/path/to/gorgeous-social-image.gif'>
<meta name='twitter:image' content='/path/to/gorgeous-social-image.gif'>
<meta itemprop='image' content='/path/to/gorgeous-social-image.gif'>
```
There’s a range of other metadata that you or your blogging software should be adding too, but that’s out of scope for this post.

Thanks for reading, go forth and make your social sharing experiences.

Mat @matryer


转自：https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html









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

转自：https://pace.dev/blog/2020/03/02/dynamically-generate-social-images-in-golang-by-mat-ryer.html









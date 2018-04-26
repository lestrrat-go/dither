# dither

Apply dithering to an arbitrary image

This is a fork of [https://github.com/esimov/dithergo](https://github.com/esimov/dithergo)

# SYNOPSIS

```go
func Example() {
  f, err := os.Open("file.png")
  if err != nil {
    return
  } 
  defer f.Close()
  
  img, _, err := image.Decode(f)
  if err != nil {
    return
  }

  ditheredImg := dither.Monochrome(dither.Burkes.Matrix(), 1.18)

  png.Encode(os.Stdout, ditheredImg)
}
```

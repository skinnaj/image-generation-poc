# Image generation POC

This proof of concept library is used to test out how to:

1. export an d3 chart to a png image
2. place this image into a canvas which allows to draw elements around it

## Exporting the d3 image

To retrieve the d3 chart as a png image including all styles and external images, the chromedp headless browser can be used.
This library provides a screenshot functionality to get a screenshot of a certain element on the page.

## Drawing the elements around the chart

The previously generated screenshot can now be places on a fogleman/gg canvas and elements can be drawn on it.
The final image can then be exported toa file or returned to the client.

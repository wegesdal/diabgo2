import PIL
from PIL import Image

import math

palette = [
    [0,0,0],
    [20, 12, 28],
    [68, 36, 52],
    [48, 52, 109],
    [78, 74, 78],
    [133, 76, 48],
    [52, 101, 36],
    [208, 70, 72],
    [117, 113, 97],
    [89, 125, 206],
    [210, 145, 44],
    [133, 149, 161],
    [109, 170, 44],
    [210, 170, 153],
    [109, 194, 202],
    [218, 212, 94],
    [222, 248, 214]
];

PIL.Image.MAX_IMAGE_PIXELS = 125829121

path = "assets/sprites/"

name = "sphere_blue.png"
im = Image.open(path+name)
orig = Image.open(path+name)

width, height = im.size
new_image = Image.new('RGBA', (width, height))

def distance(r1, g1, b1, r2, g2, b2):
    return math.sqrt(pow(r1-r2, 2) + pow(g1-g2, 2) + pow(b1-b2, 2))

def find_closest_color(pixel, palette):
    current_best = 1000
    best_color = 0
    for j in range(len(palette)):
        d = distance(palette[j][0], palette[j][1], palette[j][2], pixel[0], pixel[1], pixel[2])
        if d < current_best:
            current_best = d
            best_color = j
    return best_color

def clamp(a, b, c):
    return int(max(b, min(c, a)))

def calc_pq(pixel, q, z):
    a = [clamp(pixel[i] + q[i] * z / 16, 0, 255) for i in range(3)]
    return ((a[0], a[1], a[2], 255))

def sub_rgb(a, b):
    return (a[0]-b[0], a[1]-b[1], a[2]-b[2])

for y in range(height):
    for x in range(width):
        old_pixel = im.getpixel((x, y))
        closest_color = find_closest_color(old_pixel, palette)
        new_pixel = palette[closest_color]
        if orig.getpixel((x,y))[3] > 0:
            new_image.putpixel((x, y), (new_pixel[0], new_pixel[1], new_pixel[2], 255))
            q = sub_rgb(old_pixel, new_pixel)
            if x < width - 1:
                if orig.getpixel((x+1, y))[3] == 0:
                    # RIGHT EDGE HIGHLIGHT
                    new_image.putpixel((x+1, y), (0,0,0,255))
                im.putpixel((x+1, y), calc_pq(im.getpixel((x+1, y)), q, 7))
            if y < height - 1:
                if orig.getpixel((x, y+1))[3] == 0:
                    # BOTTOM EDGE BLACK
                    new_image.putpixel((x, y+1), (0, 0, 0,255))

                if x > 1:
                    if orig.getpixel((x-1, y))[3] == 0:
                        # LEFT EDGE BLACK
                        new_image.putpixel((x-1, y), (0,0,0,255))
                    im.putpixel((x-1, y+1),calc_pq(im.getpixel((x-1, y+1)), q, 3))
                im.putpixel((x, y+1), calc_pq(im.getpixel((x, y+1)), q, 5))
                if (x < width - 1):
                    im.putpixel((x+1, y+1), calc_pq(im.getpixel((x+1,y+1)), q, 1))

                if y > 1:
                    if orig.getpixel((x, y-1))[3] == 0:
                        # TOP EDGE BLACK
                        new_image.putpixel((x, y-1), (0,0,0,255))

new_image.save(path+'fs_'+name)


    